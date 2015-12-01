package stringutil

import (
	"fmt"
)

func QuoteStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for _, el := range values {
		result = append(result, fmt.Sprintf("%#v", el))
	}
	return result
}
