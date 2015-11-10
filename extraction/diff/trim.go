package diff

import (
	"strings"
)

func TrimLines(input []string) []string {
	result := make([]string, 0, len(input))
	for _, el := range input {
		result = append(result, strings.TrimSpace(el))
	}
	return result
}

func TrimLinesInString(input string) string {
	return strings.Join(TrimLines(strings.Split(strings.TrimSpace(input), "\n")), "\n")
}
