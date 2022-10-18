package matchers

var (
	NumberMatcher         = NewRegexMather(`[\-]{0,1}[0-9]+[\.][0-9]+|[\-]{0,1}[0-9]+`)
	PositiveNumberMatcher = NewRegexMather(`[0-9]+`)
	FirstWord             = NewRegexMather(`[^\s]+`)
)

const (
	AgreeString = "Точно!"
)

var (
	Agree = NewSelectorMatcher(
		[]string{
			"Точно",
			"Согласен",
			"Да",
			"Ага",
		},
		AgreeString,
	)
)
