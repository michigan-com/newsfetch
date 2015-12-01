package recipematcher

import (
	"regexp"
)

// TODO: replace most of these with Fuzzy Classifier

var ingredientsSubtitleRe = regexp.MustCompile(`(?i)^Ingredients$`)
var directionsSubtitleRe = regexp.MustCompile(`(?i)^Directions$`)

var servesRe = regexp.MustCompile(`(?i)^(?:Makes|Serves):`)
var servesInlineRe = regexp.MustCompile(`(?i)(?:Serves)\s*(\d.*)\.`)
var prepTimeRe = regexp.MustCompile(`(?i)(?:Preparation time):?`)
var totalTimeRe = regexp.MustCompile(`(?i)(?:Total time|Start to finish):?`)

var caloriesRe = regexp.MustCompile(`(?i)(\d calories)`)

var createdByRe = regexp.MustCompile(`(?i)(?:created by)`)
var testedByRe = regexp.MustCompile(`(?i)(?:tested by)`)
var createdAndTestedByRe = regexp.MustCompile(`(?i)(?:from and tested by|created and tested by)`)
var testKitchenRe = regexp.MustCompile(`(?i)(?:for the )?(?:Free Press )?(?:Test Kitchen)`)
var noNutricionRe = regexp.MustCompile(`(?i)(?:Nutrition information not available\.?)`)

var recipeMakesRe = regexp.MustCompile(`(?i)^(?:This recipe makes (.*?) of (.*)\.)$`)

var copyrightRe = regexp.MustCompile(`(?i)(?:Copyright|All rights reserved)`)
