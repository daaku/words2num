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
"hundred of them" come back as they were.
