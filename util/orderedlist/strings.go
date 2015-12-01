package orderedlist

import (
	"sort"
)

func ContainsString(ordered []string, element string) bool {
	i := sort.SearchStrings(ordered, element)
	if i == len(ordered) {
		return false
	} else {
		return ordered[i] == element
	}
}

func InsertString(ordered []string, element string, unique bool) []string {
	n := len(ordered)
	i := sort.SearchStrings(ordered, element)
	if i == n {
		return append(ordered, element)
	} else {
		if unique && (ordered[i] == element) {
			return ordered
		}

		ordered = append(ordered, element) // extend the array; the actual value does not matter
		copy(ordered[i+1:n+1], ordered[i:n])
		ordered[i] = element
		return ordered
	}
}
