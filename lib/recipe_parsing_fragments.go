package lib

import (
	"regexp"
	"strings"

	gq "github.com/PuerkitoBio/goquery"
)

var ingredientsSubtitle = regexp.MustCompile(`(?i)^Ingredients$`)
var directionsSubtitle = regexp.MustCompile(`(?i)^Directions$`)

var servesRe = regexp.MustCompile(`(?i)^(Makes|Serves):`)
var servesInlineRe = regexp.MustCompile(`(?i)(Serves)\s*(\d.*)\.`)
var prepTimeRe = regexp.MustCompile(`(?i)(Preparation time):?`)
var totalTimeRe = regexp.MustCompile(`(?i)(Total time|Start to finish):?`)

var caloriesRe = regexp.MustCompile(`(?i)(\d calories)`)

var startsWithNumberRe = regexp.MustCompile(`^\d`)
var ingredientRe = regexp.MustCompile(`(?i)(teaspoons?|tablespoons?|cups?|ice cubes?)`)

var createdByRe = regexp.MustCompile(`(?i)(created by)`)
var testedByRe = regexp.MustCompile(`(?i)(tested by)`)
var createdAndTestedByRe = regexp.MustCompile(`(?i)(from and tested by|created and tested by)`)
var testKitchenRe = regexp.MustCompile(`(?i)(for the )?(Free Press )?(Test Kitchen)`)
var noNutricionRe = regexp.MustCompile(`(?i)(Nutrition information not available\.?)`)

var copyrightRe = regexp.MustCompile(`(?i)(Copyright|All rights reserved)`)

var whitespaceRe = regexp.MustCompile(`\s+`)
var unwantedSentenceBreakRe = regexp.MustCompile(`[?!]`)
var potentialSentenceBreakRe = regexp.MustCompile(`\.(\s+\p{Lu}|$)`)

func extractRecipeFragments(doc *gq.Document) []RecipeFragment {
	paragraphs := doc.Find("div[itemprop=articleBody] > p, div[itemprop=articleBody] li")

	fragments := make([]RecipeFragment, 0, paragraphs.Length())

	paragraphs.Each(func(_ int, paragraph *gq.Selection) {
		html, _ := paragraph.Html()
		text := strings.TrimSpace(paragraph.Text())

		var fragment RecipeFragment

		if hasSingleChildMatching(paragraph, ".-newsgate-paragraph-cci-howto-head-") {
			fragment = &ParagraphFragment{RawHtml: html, Text: text, tag: TitleTag}
		} else if hasSingleChildMatching(paragraph, ".-newsgate-paragraph-cci-howto-components-") {
			fragment = &RecipeIngredient{Text: text}
		} else if paragraph.Is("li") {
			fragment = &RecipeIngredient{Text: text}
		} else if hasSingleChildMatching(paragraph, ".-newsgate-element-cci-howto--end") {
			fragment = &RecipeMarkerFragment{EndMarkerTag}

		} else if ingredientsSubtitle.MatchString(text) {
			fragment = &ParagraphFragment{RawHtml: html, Text: text, tag: IngredientsSubtitleTag}
		} else if directionsSubtitle.MatchString(text) {
			fragment = &ParagraphFragment{RawHtml: html, Text: text, tag: DirectionsSubtitleTag}

		} else if servesRe.MatchString(text) || prepTimeRe.MatchString(text) || totalTimeRe.MatchString(text) {
			timing := parseTimingFragment(text)
			fragment = &timing

		} else if caloriesRe.MatchString(text) {
			fragment = &NutritionData{Text: text}

		} else if createdByRe.MatchString(text) || testedByRe.MatchString(text) || createdAndTestedByRe.MatchString(text) || testKitchenRe.MatchString(text) || noNutricionRe.MatchString(text) {
			fragment = &ParagraphFragment{RawHtml: html, Text: text, tag: SignatureTag}

		} else if copyrightRe.MatchString(text) {
			fragment = &ParagraphFragment{RawHtml: html, Text: text, tag: CopyrightTag}

		} else if startsWithNumberRe.MatchString(text) {
			if ingredientRe.MatchString(text) {
				fragment = &RecipeIngredient{Text: text}
			} else {
				fragment = &ParagraphFragment{RawHtml: html, Text: text, tag: PossibleIngredientTag}
			}

		} else if text == strings.ToUpper(text) {
			fragment = &ParagraphFragment{RawHtml: html, Text: text, tag: PossibleIngredientSubdivisionTag}

		} else if hasSingleChildMatching(paragraph, "strong") {
			fragment = &ParagraphFragment{RawHtml: html, Text: text, tag: PossibleTitleTag}

		} else if isShortParagraph(text) {
			fragment = &ParagraphFragment{RawHtml: html, Text: text, tag: ShortParagraphTag}

		} else {
			fragment = &ParagraphFragment{RawHtml: html, Text: text, tag: ParagraphTag}
			if servesInlineRe.MatchString(text) {
				// TODO
			}
		}

		fragments = append(fragments, fragment)
	})

	return fragments
}

func parseTimingFragment(text string) RecipeTimingFragment {
	result := RecipeTimingFragment{}

	for _, component := range strings.Split(text, "/") {
		component = strings.TrimSpace(component)
		if value, ok := extractComponent(component, servesRe); ok {
			result.ServingSize = value
		} else if value, ok := extractComponent(component, totalTimeRe); ok {
			result.TotalTime = parseDuration(value)
		} else if value, ok := extractComponent(component, prepTimeRe); ok {
			result.PreparationTime = parseDuration(value)
		} else {
			println("Unknown duration component: >>>", component, "<<<")
			panic("Unknown duration component")
		}
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

func parseDuration(text string) *RecipeDuration {
	if text != "" {
		return &RecipeDuration{Text: text}
	} else {
		return nil
	}
}

func isShortParagraph(text string) bool {
	words := whitespaceRe.Split(text, -1)
	return (len(words) <= 15) && !unwantedSentenceBreakRe.MatchString(text) && !potentialSentenceBreakRe.MatchString(text)
}

func hasSingleChildMatching(s *gq.Selection, selector string) bool {
	children := s.Children()
	return children.Length() == 1 && children.Is(selector)
}
