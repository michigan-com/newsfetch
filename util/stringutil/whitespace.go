package stringutil

import (
	"strings"
)

func NormalizeSpace(s string) string {
	return whitespaceRe.ReplaceAllLiteralString(strings.TrimSpace(s), " ")
}
