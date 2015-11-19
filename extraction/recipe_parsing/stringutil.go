package recipe_parsing

import (
	"strings"
	"unicode"
)

func IsEntirelyUppercase(s string) bool {
	return s == strings.ToUpper(s) && strings.IndexFunc(s, unicode.IsLetter) != -1
}
