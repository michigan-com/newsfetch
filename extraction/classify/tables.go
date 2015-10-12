package classify

import (
	"strings"
)

func makeTable(spaceSeparated string) map[string]bool {
	items := strings.Fields(spaceSeparated)

	table := make(map[string]bool, len(items))

	for _, item := range items {
		table[item] = true
	}

	return table
}

func countInTable(table map[string]bool, words []string) int {
	c := 0
	for _, word := range words {
		if table[word] {
			c++
		}
	}
	return c
}
