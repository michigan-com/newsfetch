package fuzzy_classifier

import (
	"strings"
)

func NewWordSetFromString(spaceSeparated string) map[string]bool {
	items := strings.Fields(spaceSeparated)
	return NewWordSetFromList(items)
}

func NewWordSetFromList(items []string) map[string]bool {
	set := make(map[string]bool, len(items))

	for _, item := range items {
		set[item] = true
	}

	return set
}

func CountWordsInWordSet(set map[string]bool, words []string) int {
	c := 0
	for _, word := range words {
		if set[word] {
			c++
		}
	}
	return c
}
