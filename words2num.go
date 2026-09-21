// Package words2num converts numbers written as words into digits. It is meant
// for text, usually speech to text output, where numbers show up as words: "I
// have twenty three apples" becomes "I have 23 apples".
package words2num

import (
	"math"
	"strconv"
	"strings"
)

// The kinds of words we recognize. Words are looked up whole, never split, so
// the kind carries the grammatical role of the word inside a number.
const (
	kindUnit    = iota // zero .. nineteen
	kindTens           // twenty .. ninety
	kindHundred        // hundred
	kindScale          // thousand, million, billion, trillion
	kindAnd            // and
	kindPoint          // point, splits a number into its whole and fractional part
)

// maxWordLen is the length of the longest word we recognize.
const maxWordLen = 9

type word struct {
	kind  int
	value int64
}

var words = map[string]word{
	"zero":      {kindUnit, 0},
	"one":       {kindUnit, 1},
	"two":       {kindUnit, 2},
	"three":     {kindUnit, 3},
	"four":      {kindUnit, 4},
	"five":      {kindUnit, 5},
	"six":       {kindUnit, 6},
	"seven":     {kindUnit, 7},
	"eight":     {kindUnit, 8},
	"nine":      {kindUnit, 9},
	"ten":       {kindUnit, 10},
	"eleven":    {kindUnit, 11},
	"twelve":    {kindUnit, 12},
	"thirteen":  {kindUnit, 13},
	"fourteen":  {kindUnit, 14},
	"fifteen":   {kindUnit, 15},
	"sixteen":   {kindUnit, 16},
	"seventeen": {kindUnit, 17},
	"eighteen":  {kindUnit, 18},
	"nineteen":  {kindUnit, 19},

	"twenty":  {kindTens, 20},
	"thirty":  {kindTens, 30},
	"forty":   {kindTens, 40},
	"fifty":   {kindTens, 50},
	"sixty":   {kindTens, 60},
	"seventy": {kindTens, 70},
	"eighty":  {kindTens, 80},
	"ninety":  {kindTens, 90},

	"hundred": {kindHundred, 100},

	"thousand": {kindScale, 1000},
	"million":  {kindScale, 1000 * 1000},
	"billion":  {kindScale, 1000 * 1000 * 1000},
	"trillion": {kindScale, 1000 * 1000 * 1000 * 1000},

	"and": {kindAnd, 0},

	"point": {kindPoint, 0},
}

// The states of a run of number words. Only units and tens may start a run,
// which is what keeps "hundred of them" and "and you" untouched.
const (
	stateEmpty   = iota // nothing consumed
	stateUnit           // last word was zero .. nine
	stateTens           // last word was ten .. ninety
	stateHundred        // last word was hundred
	stateScale          // last word was thousand, million, billion, trillion
	stateAnd            // last word was and
	statePoint          // last word was point
	stateDigit          // last word was a digit of the fractional part
)

// run parses one run of consecutive number words into a value. It is fed every
// number word it will accept and closed when a word does not fit.
type run struct {
	state int8
	total int64
	cur   int64
	frac  []byte // digits after point, if any
	start int    // byte offset of the first word
	end   int    // byte offset after the last word consumed
}

// add feeds the word w found at [start,end). It reports false, leaving the run
// unchanged, when the word does not continue this number.
func (r *run) add(w word, start, end int) bool {
	total, cur, state := r.total, r.cur, r.state
	if state == statePoint || state == stateDigit {
		// Past the point a number only takes digits, one word each: two point
		// one three. Nothing else fits, not even hundred or another point.
		if w.kind != kindUnit || w.value > 9 {
			return false
		}
		r.frac = append(r.frac, byte('0'+w.value))
		r.state = stateDigit
		r.end = end
		return true
	}
	switch w.kind {
	case kindUnit, kindTens:
		switch {
		case state == stateEmpty:
			cur = w.value
			r.start = start
		case state == stateTens && w.value < 10:
			cur += w.value // twenty one
		case state == stateHundred || state == stateScale || state == stateAnd:
			cur += w.value // one hundred five, one thousand five
		default:
			return false // two two, fifty fifty
		}
		state = stateUnit
		if w.value >= 10 {
			state = stateTens
		}
	case kindHundred:
		// Five hundred, twenty hundred, but not hundred hundred.
		if cur <= 0 || cur >= 100 {
			return false
		}
		cur *= 100
		state = stateHundred
	case kindScale:
		if cur <= 0 {
			return false
		}
		if state != stateUnit && state != stateTens && state != stateHundred {
			return false
		}
		// Scales up to trillion fit in an int64 with room to spare. This guards
		// the vocabulary we may add later.
		if cur > math.MaxInt64/w.value || cur > (math.MaxInt64-total)/w.value {
			return false
		}
		total += cur * w.value
		cur = 0
		state = stateScale
	case kindAnd:
		// And only carries a number along: one hundred and five.
		if state != stateHundred && state != stateScale {
			return false
		}
		state = stateAnd
	case kindPoint:
		// Point only splits a number we already started reading, so "the point
		// of it" stays words.
		switch state {
		case stateUnit, stateTens, stateHundred, stateScale:
			state = statePoint
		default:
			return false
		}
	}
	r.total, r.cur, r.state = total, cur, state
	r.end = end
	return true
}

// value returns the number the run parsed to.
func (r *run) value() int64 { return r.total + r.cur }

// format returns the digits of the number, with its fractional part if it has
// one.
func (r *run) format() string {
	b := strconv.AppendInt(nil, r.value(), 10)
	if len(r.frac) > 0 {
		b = append(append(b, '.'), r.frac...)
	}
	return string(b)
}

// done reports whether the run is empty.
func (r *run) done() bool { return r.state == stateEmpty }

// valid reports whether the run parsed to a number worth replacing. A run that
// ends in "and", or in a "point" with no digits after it, was never a number.
func (r *run) valid() bool {
	return r.state != stateEmpty && r.state != stateAnd && r.state != statePoint
}

// lookup finds the word in a lower cased copy of s. The copy is made in a
// stack array so looking up words costs nothing.
func lookup(s string) (word, bool) {
	if len(s) == 0 || len(s) > maxWordLen {
		return word{}, false
	}
	var buf [maxWordLen]byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
		}
		buf[i] = c
	}
	w, ok := words[string(buf[:len(s)])]
	return w, ok
}

// wordAt returns the run of letters at i, and the index after it.
func wordAt(s string, i int) (string, int) {
	if i >= len(s) || !isLetter(s[i]) {
		return "", i
	}
	j := i + 1
	for j < len(s) && isLetter(s[j]) {
		j++
	}
	return s[i:j], j
}

func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

// hardBreak reports whether c ends a sentence, which is where a number ends
// too: "one hundred. Fifty five" is 100. 55, never 150. Commas and hyphens are
// soft, so "twenty-one" and "one million, three hundred thousand" stay one
// number.
func hardBreak(c byte) bool {
	switch c {
	case '.', '!', '?', ';', ':', '\n', '\r':
		return true
	}
	return false
}

// hasNumberWord reports whether s holds a word that could start a number.
// Every number starts with a unit or a tens word, so this is exact.
func hasNumberWord(s string) bool {
	for i := 0; i < len(s); {
		tok, end := wordAt(s, i)
		if end == i {
			i++
			continue
		}
		if w, ok := lookup(tok); ok && (w.kind == kindUnit || w.kind == kindTens) {
			return true
		}
		i = end
	}
	return false
}

// Words2Num converts numbers written as words in text into digits. The zero
// value is ready to use.
type Words2Num struct{}

// Transform replaces every number written in words in s with digits. Words
// that do not form a number are left exactly as they were.
func (w Words2Num) Transform(s string) string {
	if !hasNumberWord(s) {
		return s
	}
	var (
		out  strings.Builder
		last int
		r    run
	)
	flush := func() {
		if !r.valid() {
			return
		}
		out.WriteString(s[last:r.start])
		out.WriteString(r.format())
		last = r.end
	}
	for i := 0; i < len(s); {
		tok, end := wordAt(s, i)
		if end == i {
			if hardBreak(s[i]) && !r.done() {
				flush()
				r = run{}
			}
			i++
			continue
		}
		wd, known := lookup(tok)
		if known && r.add(wd, i, end) {
			i = end
			continue
		}
		if r.done() {
			i = end // just a word
			continue
		}
		// The run ends before this word, which may start one of its own.
		flush()
		r = run{}
		if known {
			r.add(wd, i, end)
		}
		i = end
	}
	flush()
	out.WriteString(s[last:])
	return out.String()
}
