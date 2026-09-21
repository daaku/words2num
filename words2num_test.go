package words2num

import "testing"

func TestTransform(t *testing.T) {
	cases := []struct{ in, out string }{
		// The examples.
		{"Eight hundred and fifty five", "855"},
		{"one", "1"},
		{"Five hundred", "500"},

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
		{"twenty hundred", "2000"},

		// And is only filler inside a number.
		{"one hundred and five", "105"},
		{"one thousand and fifty five", "1055"},

		// Scales.
		{"one thousand", "1000"},
		{"one thousand one", "1001"},
		{"nine hundred ninety nine thousand", "999000"},
		{"one million", "1000000"},
		{"one million one thousand one", "1001001"},
		{"twelve million three hundred forty five thousand six hundred seventy eight", "12345678"},

		// In text.
		{"I have twenty three apples", "I have 23 apples"},
		{"It is one hundred and five degrees, give or take two", "It is 105 degrees, give or take 2"},
		{"ONE HUNDRED TWENTY THREE", "123"},
		{"twenty-one", "21"},
		{"line one: ten, line two: twenty", "line 1: 10, line 2: 20"},

		// Not numbers.
		{"", ""},
		{"no numbers here", "no numbers here"},
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
		{"one million million", "1000000 million"},
		{"five hundred two two", "502 2"},
	}

	w := Words2Num{}
	for _, c := range cases {
		if got := w.Transform(c.in); got != c.out {
			t.Errorf("Transform(%q) = %q, want %q", c.in, got, c.out)
		}
	}
}
