package words2num

import (
	"strings"
	"testing"
)

func TestTransform(t *testing.T) {
	cases := []struct{ in, out string }{
		// The examples.
		{"Eight hundred and fifty five", "855"},
		{"one", "1"},
		{"Five hundred", "500"},
		{"One million three hundred thousand fifty five", "1,300,055"},

		// Plain values.
		{"zero", "0"},
		{"nine", "9"},
		{"ten", "10"},
		{"nineteen", "19"},
		{"twenty", "20"},
		{"twenty one", "21"},
		{"ninety nine", "99"},

		// Hundreds.
		{"one hundred", "100"},
		{"one hundred one", "101"},
		{"one hundred twenty three", "123"},
		{"nine hundred ninety nine", "999"},
		{"five hundred", "500"},
		{"twenty hundred", "2,000"},

		// And is only filler inside a number.
		{"one hundred and five", "105"},
		{"one thousand and fifty five", "1,055"},

		// Scales.
		{"one thousand", "1,000"},
		{"one thousand one", "1,001"},
		{"nine hundred ninety nine thousand", "999,000"},
		{"one million", "1,000,000"},
		{"one million one thousand one", "1,001,001"},
		{"twelve million three hundred forty five thousand six hundred seventy eight", "12,345,678"},

		// Billions and trillions.
		{"one billion", "1,000,000,000"},
		{"one billion one million", "1,001,000,000"},
		{"five hundred billion", "500,000,000,000"},
		{"one trillion", "1,000,000,000,000"},
		{"one trillion one", "1,000,000,000,001"},
		{"one trillion two hundred billion thirty four million five hundred sixty seven thousand eight hundred ninety", "1,200,034,567,890"},
		{"nine hundred ninety nine trillion nine hundred ninety nine billion nine hundred ninety nine million nine hundred ninety nine thousand nine hundred ninety nine", "999,999,999,999,999"},

		// Decimals: point splits a number we already started reading.
		{"forty two point one", "42.1"},
		{"zero point five", "0.5"},
		{"two point zero five", "2.05"},
		{"one point two three", "1.23"},
		{"ten point one two three", "10.123"},
		{"one hundred point five", "100.5"},
		{"one thousand point nine", "1,000.9"},
		{"it costs thirty point nine nine dollars", "it costs 30.99 dollars"},

		// Point on its own is only a word.
		{"the point is clear", "the point is clear"},
		{"point", "point"},
		{"what is the point of two things", "what is the point of 2 things"},
		{"point five", "point 5"},
		{"five point", "five point"},
		{"one hundred and point five", "one hundred and point 5"},
		{"three point five point five", "3.5 point 5"},
		{"two point twenty", "two point 20"},

		// In text.
		{"I have twenty three apples", "I have 23 apples"},
		{"It is one hundred and five degrees, give or take two", "It is 105 degrees, give or take 2"},
		{"ONE HUNDRED TWENTY THREE", "123"},
		{"twenty-one", "21"},
		{"line one: ten, line two: twenty", "line 1: 10, line 2: 20"},

		// A number never crosses the end of a sentence.
		{"The number is one hundred. Fifty five of them.", "The number is 100. 55 of them."},
		{"one hundred! fifty five", "100! 55"},
		{"one hundred\nfifty five", "100\n55"},
		{"five point. one", "five point. 1"},
		{"twenty-one", "21"},

		// Not numbers.
		{"", ""},
		{"no numbers here", "no numbers here"},
		{"their extraordinarily long words stay put", "their extraordinarily long words stay put"},
		{"hundred of them", "hundred of them"},
		{"thousand of them", "thousand of them"},
		{"and", "and"},
		{"one hundred and", "one hundred and"},
		{"and one", "and 1"},
	}

	w := Words2Num{}
	for _, c := range cases {
		if got := w.Transform(c.in); got != c.out {
			t.Errorf("Transform(%q) = %q, want %q", c.in, got, c.out)
		}
	}
}

// Runs are greedy: the longest well formed number wins, and whatever is left
// over is parsed again on its own.
func TestTransformRuns(t *testing.T) {
	cases := []struct{ in, out string }{
		{"one two", "1 2"},
		{"five five", "5 5"},
		{"twenty twenty", "20 20"},
		{"one hundred hundred", "100 hundred"},
		{"one million million", "1,000,000 million"},
		{"million trillion", "million trillion"},
		{"trillion trillion", "trillion trillion"},
		{"five hundred two two", "502 2"},
		{"two point five hundred", "2.5 hundred"},
	}

	w := Words2Num{}
	for _, c := range cases {
		if got := w.Transform(c.in); got != c.out {
			t.Errorf("Transform(%q) = %q, want %q", c.in, got, c.out)
		}
	}
}

func TestTransformNoCommas(t *testing.T) {
	cases := []struct{ in, out string }{
		{"One million three hundred thousand fifty five", "1300055"},
		{"one thousand", "1000"},
		{"nine hundred ninety nine trillion", "999000000000000"},
		{"one thousand point nine", "1000.9"},
		{"eight hundred and fifty five", "855"},
		{"one hundred. Fifty five", "100. 55"},
	}

	w := Words2Num{NoCommas: true}
	for _, c := range cases {
		if got := w.Transform(c.in); got != c.out {
			t.Errorf("Transform(%q) = %q, want %q", c.in, got, c.out)
		}
	}
}

// Commas in the input are only separators, they never end a number.
func TestTransformInputCommas(t *testing.T) {
	cases := []struct{ in, out string }{
		{"one, thousand", "1,000"},
		{"one million, three hundred thousand, fifty five", "1,300,055"},
		{"one million three hundred thousand fifty five", "1,300,055"},
		{"two, two", "2, 2"},
		{"one, two, three", "1, 2, 3"},
	}

	w := Words2Num{}
	for _, c := range cases {
		if got := w.Transform(c.in); got != c.out {
			t.Errorf("Transform(%q) = %q, want %q", c.in, got, c.out)
		}
	}
}

// Text with no numbers in it must cost nothing at all.
func TestTransformNoAllocations(t *testing.T) {
	inputs := []string{
		"",
		"no numbers here",
		"the point is clear",
		"hundred of them, and you",
		strings.Repeat("words that are nothing like digits at all, never here. ", 40),
	}

	w := Words2Num{}
	for _, in := range inputs {
		allocs := testing.AllocsPerRun(100, func() {
			if got := w.Transform(in); got != in {
				t.Errorf("Transform(%q) = %q", in, got)
			}
		})
		if allocs != 0 {
			t.Errorf("Transform(%q) allocated %g times per run, want 0", in, allocs)
		}
	}
}

func BenchmarkTransformNoNumbers(b *testing.B) {
	in := strings.Repeat("words that are nothing like digits at all, never here. ", 40)
	w := Words2Num{}
	b.ReportAllocs()
	for b.Loop() {
		sink = w.Transform(in)
	}
}

func BenchmarkTransform(b *testing.B) {
	in := strings.Repeat("twenty three, one hundred and four, one million three, ", 50)
	w := Words2Num{}
	b.ReportAllocs()
	for b.Loop() {
		sink = w.Transform(in)
	}
}

// FuzzTransform keeps Transform honest: text with no numbers is never touched,
// the result is stable, and no word shows up in the output that was not in the
// input.
func FuzzTransform(f *testing.F) {
	f.Add("Eight hundred and fifty five")
	f.Add("Forty two point one")
	f.Add("One million three hundred thousand fifty five")
	f.Add("one hundred. fifty five")
	f.Add("the point is clear")
	f.Add("one million, three hundred thousand, and fifty five")
	w := Words2Num{}
	f.Fuzz(func(t *testing.T, s string) {
		got := w.Transform(s)
		if !hasNumberWord(s) && got != s {
			t.Fatalf("changed text with no numbers: %q -> %q", s, got)
		}
		if again := w.Transform(got); again != got {
			t.Fatalf("not idempotent: %q -> %q -> %q", s, got, again)
		}
		seen := map[string]bool{}
		for _, word := range strings.Fields(s) {
			seen[strings.ToLower(word)] = true
		}
		for _, word := range strings.Fields(got) {
			word = strings.ToLower(word)
			if isNumberWord(word) && !seen[word] {
				t.Fatalf("made up the word %q: %q -> %q", word, s, got)
			}
		}
	})
}

// isNumberWord reports whether s is a word that carries a number.
func isNumberWord(s string) bool {
	w, ok := lookup(s)
	return ok && w.kind != kindAnd && w.kind != kindPoint
}

var sink string
