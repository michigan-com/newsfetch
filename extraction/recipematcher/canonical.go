package recipematcher

import (
	"github.com/michigan-com/newsfetch/util/fuzzy"
)

func CanonicalString(paragraph string) string {
	return fuzzy.CanonicalString(paragraph)
}
