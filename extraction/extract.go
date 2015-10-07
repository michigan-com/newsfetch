package extraction

import (
	"regexp"
	"strings"

	gq "github.com/PuerkitoBio/goquery"
	m "github.com/michigan-com/newsfetch/model"
)

var TWITTER_RE = regexp.MustCompile("^twitter.com/[a-zA-Z0-9_]*$")

func withoutEmptyStrings(strings []string) []string {
	result := make([]string, 0, len(strings))
	for _, el := range strings {
		if el != "" {
			result = append(result, el)
		}
	}
	return result
}

func extractBodyFromDocument(doc *gq.Document, fromJSON bool, includeTitle bool) *m.ExtractedBody {
	msg := new(m.Messages)

	var paragraphs *gq.Selection
	if fromJSON {
		paragraphs = doc.Find("p")
	} else {
		paragraphs = doc.Find("div[itemprop=articleBody] > p")
	}

	// remove contact info at the end of the article
	paragraphs.Find("span.-newsgate-paragraph-cci-endnote-contact-").Remove()
	paragraphs.Find("span.-newsgate-paragraph-cci-endnote-contrib-").Remove()

	ignoreRemaining := false
	paragraphStrings := paragraphs.Map(func(i int, paragraph *gq.Selection) string {
		if ignoreRemaining {
			return ""
		}
		for _, selector := range [...]string{"span.-newsgate-character-cci-tagline-name-", "span.-newsgate-paragraph-cci-infobox-head-"} {
			if el := paragraph.Find(selector); el.Length() > 0 {
				ignoreRemaining = true
				return ""
			}
		}

		text := strings.TrimSpace(paragraph.Text())

		if TWITTER_RE.MatchString(text) {
			return ""
		}

		marker := ""

		for _, selector := range [...]string{"span.-newsgate-paragraph-cci-subhead-lead-", "span.-newsgate-paragraph-cci-subhead-"} {
			if el := paragraph.Find(selector); el.Length() > 0 {
				marker = "### "
				break
			}
		}

		return marker + text
	})

	content := make([]string, 0, len(paragraphStrings)+1)
	if includeTitle {
		title := ExtractTitleFromDocument(doc)
		content = append(content, title)
	}

	content = append(content, withoutEmptyStrings(paragraphStrings)...)

	body := strings.Join(content, "\n")
	recipeData, recipeMsg := ExtractRecipes(doc)
	msg.AddMessages("recipes", recipeMsg)
	extracted := m.ExtractedBody{body, recipeData, msg}
	return &extracted
}

func ExtractTitleFromDocument(doc *gq.Document) string {
	title := doc.Find("h1[itemprop=headline]")
	return strings.TrimSpace(title.Text())
}
