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

`Transform` walks the text word by word and feeds each recognized number word to
a small state machine (`run`) that accumulates `total` and `cur`:

- Only a unit or a tens word may start a number, so "hundred of them" and "and
  you" stay as they are.
- Runs are greedy and never ambiguous: the longest well formed number wins and
  leftovers are parsed again on their own ("five hundred two two" is "502 2").
- A run ends at a hard separator (`.!?;:` and newlines) but not at a soft one,
  so "twenty-one" and "one million, three hundred thousand" are still one
  number while "one hundred. Fifty five" is two.
- A word that does not fit closes the run and is then given a chance to start a
  run of its own. Nothing is ever dropped: a run that turns out not to be a
  number is left in the text untouched.

## Costs

Text with no numbers must cost zero allocations. Two things hold that up, and
both are load bearing:

- `hasNumberWord` pre-scans and returns the input unchanged. Every number starts
  with a unit or tens word, so this prescan is exact.
- `lookup` lower cases into a fixed stack array and only then does the map
  lookup, which Go does not allocate. Words longer than `maxWordLen` are
  rejected early to keep that array on the stack, so bump `maxWordLen` whenever a
  longer word is added to the vocabulary.
