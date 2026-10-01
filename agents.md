# words2num

Go library that rewrites numbers written as words into digits inside text.
Numbers come from speech to text most of the time.

## Layout

- One package, one source file (`words2num.go`) and one test file
  (`words2num_test.go`), both named after the package. Keep it that way.
- No dependencies, standard library only. Tests use plain `testing` and plain
  comparisons, not a testify/ensure style helper.
- Format with gofmt, check with `go vet ./...`, run `go test ./...`.

## Parsing

Vocabulary covers `zero`..`nineteen`, `twenty`..`ninety`, `hundred`, the scales
`thousand`, `million`, `billion`, `trillion`, and the filler `and`. Plurals like
"thousands" are not numbers, and neither is a scale word on its own.

`point` splits a number into a whole and a fractional part, and only starts a
number that is already being read. Only the digit words `zero`..`nine` follow it,
one digit each, so `two point five hundred` is "2.5 hundred" and a run that ends
on `point` is not a number at all.

`Words2Num` is the config and `Replace` is the only method. It is named after
`strings.Replacer.Replace` so both can be used through the same
`Replace(string) string` interface. `NoCommas` turns off grouping of three
digits; grouping only ever applies to the whole part, so `one thousand point
nine` is "1,000.9".

`Replace` walks the text word by word and feeds each recognized number word to
a small state machine (`run`) that accumulates `total` and `cur`:

- Only a unit or a tens word may start a number, so "hundred of them" and "and
  you" stay as they are.
- A unit word only starts a number when the word before it allows one: after a
  determiner, "one" is the pronoun and not a count, so "this one", "the one"
  and "no one" stay as they are. `determiners` is deliberately short and
  `isPronounOne` is the only place that reads it. "twenty one" is unaffected,
  because there the run is already going.
- "one" is the pronoun in one other shape, when the word after it is "of":
  "one of my friends" is not "1 of my friends". The run machine only ever sees
  behind itself, so `isPartitiveOne` looks ahead instead, and both
  `hasNumberWord` and `Replace` call it. A count in front of "of" is
  unaffected ("two of my friends" is "2 of my friends"), and so is a scale
  word ("one hundred of them" is "100 of them") because there the run is no
  longer a bare "one".
- Runs are greedy and never ambiguous: the longest well formed number wins and
  leftovers are parsed again on their own ("five hundred two two" is "502 2").
- A run ends at a hard separator (`.!?;:` and newlines) but not at a soft one,
  so "twenty-one" and "one million, three hundred thousand" are still one
  number while "one hundred. Fifty five" is two.
- A word that does not fit closes the run and is then given a chance to start a
  run of its own. Nothing is ever dropped: a run that turns out not to be a
  number is left in the text untouched.

## Costs

Text with no numbers must cost zero allocations. These are the load bearing
choices that keep it that way:

- `hasNumberWord` pre-scans and returns the input unchanged. Every number starts
  with a unit or tens word that is not one of the pronouns, so this prescan is
  exact, and text like "this one" or "one of them" still costs nothing.
- `lookup` lower cases into a fixed stack array and only then does the map
  lookup, which Go does not allocate. Words longer than `maxWordLen` are
  rejected early to keep that array on the stack, so bump `maxWordLen` whenever a
  longer word is added to the vocabulary.
- `Replace` builds the result in a single byte slice and everything formats
  into it (`run.format` and `appendInt` take a destination, `strconv.AppendInt`
  fills a stack array). Text with numbers costs two allocations however many
  numbers it holds; `TestReplaceNoAllocations` and the benchmarks in
  `words2num_test.go` keep an eye on this.
- `FuzzReplace` checks the two properties that are easy to break while
  changing the parser: text with no number words comes back untouched, and
  transforming twice changes nothing after the first pass. Run it with
  `go test -fuzz FuzzReplace` before shipping a parser change.
