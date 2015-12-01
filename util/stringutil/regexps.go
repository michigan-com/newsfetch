package stringutil

import (
	"regexp"
)

var whitespaceRe = regexp.MustCompile(`\s+`)
var unwantedSentenceBreakRe = regexp.MustCompile(`[?!]`)
var potentialSentenceBreakRe = regexp.MustCompile(`\.(\s+\p{Lu}|$)`)
