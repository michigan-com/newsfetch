package body_parsing

import (
	"strings"

	gq "github.com/PuerkitoBio/goquery"
	"github.com/michigan-com/newsfetch/extraction/classify"
	"github.com/michigan-com/newsfetch/extraction/dateline"
	"github.com/michigan-com/newsfetch/extraction/recipe_parsing"
	m "github.com/michigan-com/newsfetch/model"
)

func withoutEmptyStrings(strings []string) []string {
	result := make([]string, 0, len(strings))
	for _, el := range strings {
		if el != "" {
			result = append(result, el)
		}
	}
	return result
}

func ExtractBodyFromDocument(doc *gq.Document, fromJSON bool, includeTitle bool) *m.ExtractedBody {
	msg := new(m.Messages)

	var paragraphs *gq.Selection
	if fromJSON {
		paragraphs = doc.Find("p")
	} else {
		paragraphs = doc.Find("div[itemprop=articleBody] > p")
	}

	// remove contact info at the end of the article (might not be needed any more when parsing
	// HTML from JSON?)
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

		if worthy, _ := classify.IsWorthyParagraph(text); !worthy {
			return ""
		}

		//marker := ""

		for _, selector := range [...]string{"span.-newsgate-paragraph-cci-subhead-lead-", "span.-newsgate-paragraph-cci-subhead-"} {
			if el := paragraph.Find(selector); el.Length() > 0 {
				//marker = "### "
				return ""
				break
			}
		}

		return text
	})

	if len(paragraphStrings) > 0 {
		paragraphStrings[0] = dateline.RmDateline(paragraphStrings[0])
	}

	content := make([]string, 0, len(paragraphStrings)+1)
	if includeTitle {
		title := ExtractTitleFromDocument(doc)
		content = append(content, title)
	}

	content = append(content, withoutEmptyStrings(paragraphStrings)...)

	body := strings.Join(content, "\n")
	recipeData, recipeMsg := recipe_parsing.ExtractRecipes(doc)
	msg.AddMessages("recipes", recipeMsg)
	extracted := m.ExtractedBody{body, recipeData, msg}
	return &extracted
}

func ExtractTitleFromDocument(doc *gq.Document) string {
	title := doc.Find("h1[itemprop=headline]")
	return strings.TrimSpace(title.Text())
}
