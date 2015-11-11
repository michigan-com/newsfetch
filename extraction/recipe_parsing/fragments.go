package recipe_parsing

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/net/html"

	gq "github.com/PuerkitoBio/goquery"
	m "github.com/michigan-com/newsfetch/model"
)

var ingredientsSubtitle = regexp.MustCompile(`(?i)^Ingredients$`)
var directionsSubtitle = regexp.MustCompile(`(?i)^Directions$`)

var servesRe = regexp.MustCompile(`(?i)^(?:Makes|Serves):`)
var servesInlineRe = regexp.MustCompile(`(?i)(?:Serves)\s*(\d.*)\.`)
var prepTimeRe = regexp.MustCompile(`(?i)(?:Preparation time):?`)
var totalTimeRe = regexp.MustCompile(`(?i)(?:Total time|Start to finish):?`)

var caloriesRe = regexp.MustCompile(`(?i)(\d calories)`)

var startsWithNumberRe = regexp.MustCompile(`^\d`)
var ingredientRe = regexp.MustCompile(`(?i)(?:teaspoons?|tablespoons?|cups?|ice cubes?)`)

var createdByRe = regexp.MustCompile(`(?i)(?:created by)`)
var testedByRe = regexp.MustCompile(`(?i)(?:tested by)`)
var createdAndTestedByRe = regexp.MustCompile(`(?i)(?:from and tested by|created and tested by)`)
var testKitchenRe = regexp.MustCompile(`(?i)(?:for the )?(?:Free Press )?(?:Test Kitchen)`)
var noNutricionRe = regexp.MustCompile(`(?i)(?:Nutrition information not available\.?)`)

var recipeMakesRe = regexp.MustCompile(`(?i)^(?:This recipe makes (.*?) of (.*)\.)$`)

var copyrightRe = regexp.MustCompile(`(?i)(?:Copyright|All rights reserved)`)

var whitespaceRe = regexp.MustCompile(`\s+`)
var unwantedSentenceBreakRe = regexp.MustCompile(`[?!]`)
var potentialSentenceBreakRe = regexp.MustCompile(`\.(\s+\p{Lu}|$)`)

func extractRecipeFragments(doc *gq.Document, msg *m.Messages) []m.RecipeFragment {
	paragraphs := doc.Find("div[itemprop=articleBody] > p, div[itemprop=articleBody] li")
	if paragraphs.Length() == 0 {
		paragraphs = doc.Find("body > p")
	}

	fragments := make([]m.RecipeFragment, 0, paragraphs.Length())

	paragraphs.Each(func(_ int, paragraph *gq.Selection) {
		html, _ := paragraph.Html()
		fixupTextNodesBeforeSubSup(paragraph)

		text := strings.TrimSpace(paragraph.Text())
		text = whitespaceRe.ReplaceAllLiteralString(text, " ")

		var fragment m.RecipeFragment

		if hasSingleChildMatching(paragraph, ".-newsgate-paragraph-cci-howto-head-") {
			fragment = &m.ParagraphFragment{RawHtml: html, Text: text, TagF: m.TitleTag}
		} else if hasSingleChildMatching(paragraph, ".-newsgate-paragraph-cci-howto-components-") {
			fragment = &m.ParagraphFragment{RawHtml: html, Text: text, TagF: m.IngredientTag}
		} else if paragraph.Is("li") {
			fragment = &m.ParagraphFragment{RawHtml: html, Text: text, TagF: m.IngredientTag}
		} else if hasSingleChildMatching(paragraph, ".-newsgate-element-cci-howto--end") {
			fragment = &m.RecipeMarkerFragment{m.EndMarkerTag}

		} else if ingredientsSubtitle.MatchString(text) {
			fragment = &m.ParagraphFragment{RawHtml: html, Text: text, TagF: m.IngredientsSubtitleTag}
		} else if directionsSubtitle.MatchString(text) {
			fragment = &m.ParagraphFragment{RawHtml: html, Text: text, TagF: m.DirectionsSubtitleTag}

		} else if servesRe.MatchString(text) || prepTimeRe.MatchString(text) || totalTimeRe.MatchString(text) {
			timing := parseTimingFragment(text, msg)
			fragment = &timing

		} else if caloriesRe.MatchString(text) {
			fragment = &m.NutritionData{Text: text}

		} else if createdByRe.MatchString(text) || testedByRe.MatchString(text) || createdAndTestedByRe.MatchString(text) || testKitchenRe.MatchString(text) || noNutricionRe.MatchString(text) {
			fragment = &m.ParagraphFragment{RawHtml: html, Text: text, TagF: m.SignatureTag}

		} else if copyrightRe.MatchString(text) {
			fragment = &m.ParagraphFragment{RawHtml: html, Text: text, TagF: m.CopyrightTag}

		} else if match := recipeMakesRe.FindStringSubmatch(text); match != nil && isShortTextWithMaxWordCount(match[1], 5) && isShortTextWithMaxWordCount(match[2], 5) {
			fragment = &m.RecipeTimingFragment{TagF: m.ServingSizeAltTimingTag, ServingSize: match[1]}

		} else if startsWithNumberRe.MatchString(text) {
			if ingredientRe.MatchString(text) {
				fragment = &m.ParagraphFragment{RawHtml: html, Text: text, TagF: m.IngredientTag}
			} else {
				fragment = &m.ParagraphFragment{RawHtml: html, Text: text, TagF: m.PossibleIngredientTag}
			}

		} else if text == strings.ToUpper(text) {
			fragment = &m.ParagraphFragment{RawHtml: html, Text: text, TagF: m.PossibleIngredientSubdivisionTag}

		} else if hasSingleChildMatching(paragraph, "strong") || hasSingleChildMatching(paragraph, "b") {
			fragment = &m.ParagraphFragment{RawHtml: html, Text: text, TagF: m.PossibleTitleTag}

		} else if isShortParagraph(text) {
			fragment = &m.ParagraphFragment{RawHtml: html, Text: text, TagF: m.ShortParagraphTag}

		} else {
			fragment = &m.ParagraphFragment{RawHtml: html, Text: text, TagF: m.ParagraphTag}
			if servesInlineRe.MatchString(text) {
				// TODO
			}
		}

		fragments = append(fragments, fragment)
	})

	return fragments
}

func parseTimingFragment(text string, msg *m.Messages) m.RecipeTimingFragment {
	result := m.RecipeTimingFragment{TagF: m.TimingTag}

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
		if startsWithDigit(el) && len(result) > 0 {
			last := result[len(result)-1]
			if endsWithDigit(last) {
				result[len(result)-1] = last + "/" + el
				continue
			}
		}

		result = append(result, el)
	}
	return result
}

func startsWithDigit(s string) bool {
	r := []rune(s)
	return len(s) > 0 && unicode.IsDigit(r[0])
}
func endsWithDigit(s string) bool {
	r := []rune(s)
	return len(s) > 0 && unicode.IsDigit(r[len(r)-1])
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

func parseDuration(text string) *m.RecipeDuration {
	if text != "" {
		return &m.RecipeDuration{Text: text}
	} else {
		return nil
	}
}

func isShortParagraph(text string) bool {
	return isShortTextWithMaxWordCount(text, 15)
}

func isShortTextWithMaxWordCount(text string, maxWords int) bool {
	words := whitespaceRe.Split(text, -1)
	return (len(words) <= maxWords) && !unwantedSentenceBreakRe.MatchString(text) && !potentialSentenceBreakRe.MatchString(text)
}

func hasSingleChildMatching(s *gq.Selection, selector string) bool {
	parent := s.Nodes[0]
	childElCount := 0
	for child := parent.FirstChild; child != nil; child = child.NextSibling {
		switch child.Type {
		case html.CommentNode:
		case html.TextNode:
			if child.Data != "" {
				return false
			}
		case html.ElementNode:
			childElCount++
		default:
			return false
		}
	}

	if childElCount != 1 {
		return false
	}

	children := s.Children()
	return children.Length() == 1 && children.Is(selector)
}

func fixupTextNodesBeforeSubSup(s *gq.Selection) {
	for _, node := range s.Nodes {
		fixupTextNodesBeforeSubSupInNode(node)
	}
}

func fixupTextNodesBeforeSubSupInNode(parent *html.Node) {
	for child := parent.FirstChild; child != nil; child = child.NextSibling {
		if (child.Type == html.ElementNode) && (child.Data == "sup") {
			if prev := child.PrevSibling; prev != nil && prev.Type == html.TextNode {
				if endsWithDigit(prev.Data) {
					prev.Data = prev.Data + " "
				}
			}

		}
		if child.Type == html.ElementNode {
			fixupTextNodesBeforeSubSupInNode(child)
		}
	}
}
