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

		// Billions and trillions.
		{"one billion", "1000000000"},
		{"one billion one million", "1001000000"},
		{"five hundred billion", "500000000000"},
		{"one trillion", "1000000000000"},
		{"one trillion one", "1000000000001"},
		{"one trillion two hundred billion thirty four million five hundred sixty seven thousand eight hundred ninety", "1200034567890"},
		{"nine hundred ninety nine trillion nine hundred ninety nine billion nine hundred ninety nine million nine hundred ninety nine thousand nine hundred ninety nine", "999999999999999"},

		// Decimals: point splits a number we already started reading.
		{"forty two point one", "42.1"},
		{"zero point five", "0.5"},
		{"two point zero five", "2.05"},
		{"one point two three", "1.23"},
		{"ten point one two three", "10.123"},
		{"one hundred point five", "100.5"},
		{"one thousand point nine", "1000.9"},
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
