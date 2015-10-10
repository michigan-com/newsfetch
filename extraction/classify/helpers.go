package classify

import (
	"strings"
	"unicode"

	"github.com/michigan-com/newsfetch/extraction/split"
)

func stringsToLower(input []string) []string {
	result := make([]string, 0, len(input))
	for _, el := range input {
		el = strings.Map(unicode.ToLower, el)
		result = append(result, el)
	}
	return result
}

func countIrregulars(words []string) int {
	c := 0
	for _, word := range words {
		if !split.IsRegularWord(word) {
			c++
		}
	}
	return c
}

func percentageOf(amount, total int) int {
	if total == 0 {
		return 0
	} else {
		return amount * 100 / total
	}
}
