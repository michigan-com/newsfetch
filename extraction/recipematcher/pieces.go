package recipematcher

import (
	"regexp"
	"strings"

	"github.com/michigan-com/newsfetch/model/recipetypes"
	"github.com/michigan-com/newsfetch/util/messages"
	"github.com/michigan-com/newsfetch/util/stringutil"
)

func parseTimingFragment(text string, msg *messages.Messages) recipetypes.RecipeTimingFragment {
	result := recipetypes.RecipeTimingFragment{TagF: recipetypes.TimingTag}

	for _, component := range fractionAwareSplitBySlashes(text) {
		component = strings.TrimSpace(component)
		if value, ok := extractComponent(component, servesRe); ok {
			result.ServingSize = value
		} else if value, ok := extractComponent(component, totalTimeRe); ok {
			result.TotalTime = parseDuration(value)
		} else if value, ok := extractComponent(component, prepTimeRe); ok {
			result.PreparationTime = parseDuration(value)
		} else {
			msg.AddWarningf("Unknown duration component: %#v", component)
		}
	}

	return result
}

func fractionAwareSplitBySlashes(text string) []string {
	components := strings.Split(text, "/")

	// rejoin fractions
	result := make([]string, 0, len(components))
	for _, el := range components {
		if stringutil.StartsWithDigit(el) && len(result) > 0 {
			last := result[len(result)-1]
			if stringutil.EndsWithDigit(last) {
				result[len(result)-1] = last + "/" + el
				continue
			}
		}

		result = append(result, el)
	}
	return result
}

func extractComponent(component string, re *regexp.Regexp) (string, bool) {
	// TODO: make sure there's no extra text before the match of the regexp
	if re.MatchString(component) {
		value := strings.TrimSpace(re.ReplaceAllLiteralString(component, ""))
		return value, true
	} else {
		return "", false
	}
}

func parseDuration(text string) *recipetypes.RecipeDuration {
	if text != "" {
		return &recipetypes.RecipeDuration{Text: text}
	} else {
		return nil
	}
}
