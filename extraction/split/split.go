package split

import (
	"strings"
	"unicode"
)

func SplitWords(text string) []string {
	runs := strings.Fields(text)

	words := make([]string, 0, len(runs))
	for _, run := range runs {
		trimmed := strings.TrimFunc(run, isTrimmable)
		if len(trimmed) > 0 {
			words = append(words, trimmed)
		}
	}
	return words
}

func IsRegularWord(word string) bool {
	return strings.IndexFunc(word, isIrregularRune) < 0
}

func isTrimmable(r rune) bool {
	return !unicode.IsLetter(r) && !isSpecialCharacter(r)
}

func isIrregularRune(r rune) bool {
	return !unicode.IsLetter(r)
}

// Twitter usernames, hashtags, URLs
const special = "#@/"

func isSpecialCharacter(r rune) bool {
	return unicode.IsNumber(r) || strings.ContainsRune(special, r)
}
