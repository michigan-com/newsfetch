package stringutil

import (
	"strings"
	"unicode"
)

func IsEntirelyUppercase(s string) bool {
	return s == strings.ToUpper(s) && strings.IndexFunc(s, unicode.IsLetter) != -1
}

// TODO: a better way!
func IsSingleSentence(text string) bool {
	return !unwantedSentenceBreakRe.MatchString(text) && !potentialSentenceBreakRe.MatchString(text)
}

func IsSingleSentenceWithMaxWordCount(text string, maxWords int) bool {
	words := strings.Fields(text)
	return (len(words) <= maxWords) && IsSingleSentence(text)
}

func StartsWithDigit(s string) bool {
	r := []rune(s)
	return len(s) > 0 && unicode.IsDigit(r[0])
}
func EndsWithDigit(s string) bool {
	r := []rune(s)
	return len(s) > 0 && unicode.IsDigit(r[len(r)-1])
}
