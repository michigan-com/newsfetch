package recipematcher

import (
	"github.com/michigan-com/newsfetch/util/fuzzy"
)

var classifier = fuzzy.NewFromString(ingredientRules, directionRules)
