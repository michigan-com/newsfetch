package stringutil

import (
	"sort"
)

func SortedStrings(list []string) []string {
	result := make([]string, 0, len(list))
	copy(result, list)
	sort.Strings(result)
	return result
}

func StringSetKeys(set map[string]bool) []string {
	result := make([]string, 0, len(set))
	for el, _ := range set {
		result = append(result, el)
	}
	return result
}

func SortedStringSetKeys(set map[string]bool) []string {
	result := StringSetKeys(set)
	sort.Strings(result)
	return result
}

func SetFromStrings(list []string) map[string]bool {
	set := make(map[string]bool, len(list))

	for _, item := range list {
		set[item] = true
	}

	return set
}
