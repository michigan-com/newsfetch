package recipe_parsing

import (
	"fmt"
	"strings"
	_ "unicode"

	fuz "github.com/michigan-com/newsfetch/extraction/fuzzy_classifier"
)

type ConfidenceLevel int

const (
	Negative ConfidenceLevel = iota
	None
	Possible
	Likely
	Perfect
)

type Match struct {
	Confidence ConfidenceLevel
	Rationale  string
}

type Matcher struct {
	ingredientClassifier *fuz.Classifier
	directionClassifier  *fuz.Classifier
}

func NewMatcher() *Matcher {
	return &Matcher{
		ingredientClassifier: NewIngredientClassifier(),
		directionClassifier:  NewDirectionClassifier(),
	}
}

func (matcher *Matcher) MatchIngredientSubhead(s string) Match {
	s = strings.Replace(s, "(optional)", "", 1)

	if IsEntirelyUppercase(s) {
		return Match{Confidence: Likely, Rationale: "entirely uppercase"}
	}

	return Match{Confidence: None}
}

func (matcher *Matcher) MatchIngredient(s string) Match {
	r := matcher.ingredientClassifier.Process(s)

	if r.HasTag("@not_ingredient") {
		return Match{Confidence: Negative, Rationale: "blacklisted"}
	}

	startsWithQuantity := r.HasTagAt("@quantity", 0)

	ingredientRange, ingredientFound := r.GetTagMatch("@ingredient")

	if ingredientFound && r.RangeCoversEntireInput(ingredientRange) {
		return Match{Confidence: Perfect, Rationale: "perfect match"}
	}

	if !ingredientFound && startsWithQuantity {
		return Match{Confidence: Likely, Rationale: "starts with a quantity"}
	}

	if ingredientFound {
		ingredientString := r.GetRangeString(ingredientRange, fuz.Raw)
		return Match{Confidence: Possible, Rationale: fmt.Sprintf("matched a substring: %#v", ingredientString)}
	}

	return Match{Confidence: None}
}

func (matcher *Matcher) MatchDirection(s string) Match {
	r := matcher.directionClassifier.Process(s)

	if r.HasTag("@direction") {
		all := r.GetAllTagMatchStrings("@direction", fuz.Trimmed)
		return Match{Confidence: Likely, Rationale: fmt.Sprintf("matched %s", strings.Join(mapToQuoted(all), ", "))}
	}
	return Match{Confidence: None}
}

func mapToQuoted(values []string) []string {
	result := make([]string, 0, len(values))
	for _, el := range values {
		result = append(result, fmt.Sprintf("%#v", el))
	}
	return result
}
