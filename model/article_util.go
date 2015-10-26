package model

import (
	"strings"
)

func FilterArticleURLsBySection(urls []string, section string) []string {
	result := make([]string, 0, len(urls))
	for _, el := range urls {
		if strings.Contains(el, "/story/"+section+"/") {
			result = append(result, el)
		}
	}
	return result
}
