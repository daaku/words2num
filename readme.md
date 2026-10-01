# words2num

words2num converts numbers in words to digits.

```go
w := words2num.Words2Num{}
w.Replace("Buy twenty three apples")                       // Buy 23 apples
w.Replace("Forty two point one")                           // 42.1
w.Replace("One million three hundred thousand fifty five") // 1,300,055

w = words2num.Words2Num{NoCommas: true}
w.Replace("One million three hundred thousand fifty five") // 1300055
```

Words that do not make a number are left alone, so "the point is clear" and
"hundred of them" come back as they were. A "one" that stands on its own is
left as the word too, because it is the pronoun about as often as it is a
count: "one apple", "this one" and "one of my friends" are untouched, while
"one hundred" is "100", "one twenty three" is "1 23" and "twenty one" is 21.
