package recipematcher

import (
	"fmt"
	"strings"

	"github.com/michigan-com/newsfetch/model/bodytypes"
	"github.com/michigan-com/newsfetch/model/recipetypes"
	"github.com/michigan-com/newsfetch/util/fuzzy"
	"github.com/michigan-com/newsfetch/util/messages"
	"github.com/michigan-com/newsfetch/util/stringutil"
)

func Match(paragraph bodytypes.Paragraph, msg *messages.Messages) (recipetypes.Role, bodytypes.Match, []recipetypes.Fragment) {
	text := paragraph.Text
	r := classifier.Process(text)

	if ingredientsSubtitleRe.MatchString(text) {
		return recipetypes.IngredientsHeading, bodytypes.Match{bodytypes.Likely, "matched ingredients subtitle regexp"}, nil
	}
	if directionsSubtitleRe.MatchString(text) {
		return recipetypes.DirectionsHeading, bodytypes.Match{bodytypes.Likely, "matched directions subtitle regexp"}, nil
	}

	if paragraph.HasTag(bodytypes.NewsgateHead) {
		return recipetypes.Title, bodytypes.Match{bodytypes.Likely, "found NewsgateHead"}, nil
	}

	if paragraph.HasTag(bodytypes.NewsgateEnd) {
		return recipetypes.EndMarker, bodytypes.Match{bodytypes.Perfect, "found NewsgateEnd"}, nil
	}

	ishm := matchIngredientSubhead(text)
	if ishm.Confidence >= bodytypes.Likely {
		return recipetypes.IngredientSubsectionHeading, ishm, nil
	}

	im := matchIngredient(text, r)
	dm := matchDirection(text, r)

	if im.Confidence >= bodytypes.Likely && dm.Confidence >= bodytypes.Likely {
		return recipetypes.Conflict, bodytypes.Match{im.Confidence, fmt.Sprintf("ingredient (%v) AND direction (%v)", im.Rationale, dm.Rationale)}, nil

	} else if im.Confidence >= bodytypes.Possible && im.Confidence > dm.Confidence {
		return recipetypes.Ingredient, im, nil

	} else if dm.Confidence >= bodytypes.Possible && dm.Confidence > im.Confidence {
		return recipetypes.Direction, dm, nil

		// can never be triggered currently
		// } else if im.Confidence == dm.Confidence && im.Confidence >= bodytypes.Possible {
		// 	return recipetypes.Unknown, bodytypes.Match{im.Confidence, fmt.Sprintf("soft conflict: ingredient (%v) AND direction (%v)", im.Rationale, dm.Rationale)}, nil

	} else if im.Confidence == bodytypes.Negative {
		return recipetypes.Direction, bodytypes.Match{bodytypes.Likely, fmt.Sprintf("matched blacklist for ingredient: %v", im.Rationale)}, nil
	}

	if paragraph.HasTag(bodytypes.NewsgateComponent) {
		return recipetypes.Ingredient, bodytypes.Match{bodytypes.Likely, "found NewsgateComponent"}, nil
	}
	if paragraph.HasTag(bodytypes.ListItem) {
		return recipetypes.Ingredient, bodytypes.Match{bodytypes.Likely, "found ListItem"}, nil
	}

	if servesRe.MatchString(text) || prepTimeRe.MatchString(text) || totalTimeRe.MatchString(text) {
		timing := parseTimingFragment(text, msg)
		return recipetypes.Special, bodytypes.Match{bodytypes.Likely, "matched serving/timing regexp"}, []recipetypes.Fragment{timing}
	}

	if caloriesRe.MatchString(text) {
		fragment := &recipetypes.NutritionData{Text: text}
		return recipetypes.Special, bodytypes.Match{bodytypes.Likely, "matched calories regexp"}, []recipetypes.Fragment{fragment}
	}

	if createdByRe.MatchString(text) || testedByRe.MatchString(text) || createdAndTestedByRe.MatchString(text) || testKitchenRe.MatchString(text) || noNutricionRe.MatchString(text) {
		return recipetypes.ShitTail, bodytypes.Match{bodytypes.Likely, "matched regexp"}, nil
	}

	if copyrightRe.MatchString(text) {
		return recipetypes.ShitTail, bodytypes.Match{bodytypes.Likely, "matched copyright regexp"}, nil
	}

	if match := recipeMakesRe.FindStringSubmatch(text); match != nil && stringutil.IsSingleSentenceWithMaxWordCount(match[1], 5) && stringutil.IsSingleSentenceWithMaxWordCount(match[2], 5) {
		fragment := &recipetypes.RecipeTimingFragment{TagF: recipetypes.ServingSizeAltTimingTag, ServingSize: match[1]}
		return recipetypes.OutOfBandSpecial, bodytypes.Match{bodytypes.Likely, "matched alt makes regexp"}, []recipetypes.Fragment{fragment}
	}

	if paragraph.HasTag(bodytypes.Bold) {
		return recipetypes.Title, bodytypes.Match{bodytypes.Possible, "bold"}, nil
	}

	// TODO: remove this particular case?
	if stringutil.IsSingleSentenceWithMaxWordCount(text, 15) {
		return recipetypes.Ingredient, bodytypes.Match{bodytypes.DistantlyPossible, "short paragraph"}, nil
	}

	return recipetypes.Direction, bodytypes.Match{bodytypes.DistantlyPossible, "nothing matched"}, nil
}

func matchIngredientSubhead(s string) bodytypes.Match {
	s = strings.Replace(s, "(optional)", "", 1)

	if stringutil.IsEntirelyUppercase(s) {
		return bodytypes.Match{bodytypes.Likely, "entirely uppercase"}
	}

	return bodytypes.Match{bodytypes.None, "nothing matched"}
}

func matchIngredient(s string, r *fuzzy.Result) bodytypes.Match {
	if r.HasTag("@not_ingredient") {
		return bodytypes.Match{bodytypes.Negative, "blacklisted"}
	}

	startsWithQuantity := r.HasTagAt("@quantity", 0)

	ingredientRange, ingredientFound := r.GetTagMatch("@ingredient")

	if ingredientFound && r.RangeCoversEntireInput(ingredientRange) {
		return bodytypes.Match{bodytypes.Perfect, "perfect match"}
	}

	if !ingredientFound && startsWithQuantity {
		return bodytypes.Match{bodytypes.Likely, "starts with a quantity"}
	}

	if ingredientFound {
		ingredientString := r.GetRangeString(ingredientRange, fuzzy.Raw)
		return bodytypes.Match{bodytypes.Possible, fmt.Sprintf("matched a substring: %#v", ingredientString)}
	}

	return bodytypes.Match{bodytypes.None, "nothing matched"}
}

func matchDirection(s string, r *fuzzy.Result) bodytypes.Match {
	if r.HasTag("@direction") {
		all := r.GetAllTagMatchStrings("@direction", fuzzy.Trimmed)
		return bodytypes.Match{bodytypes.Likely, fmt.Sprintf("matched %s", strings.Join(stringutil.QuoteStrings(all), ", "))}
	}
	return bodytypes.Match{bodytypes.None, "nothing matched"}
}
