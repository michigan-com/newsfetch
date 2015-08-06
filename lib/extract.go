package lib

import (
	"strings"

	gq "github.com/PuerkitoBio/goquery"
)

var logger = GetLogger()

func withoutEmptyStrings(strings []string) []string {
	result := make([]string, 0, len(strings))
	for _, el := range strings {
		if el != "" {
			result = append(result, el)
		}
	}
	return result
}

func extractBodyFromDocument(doc *gq.Document, includeTitle bool) (string, error) {
	paragraphs := doc.Find("div[itemprop=articleBody] p")

	// remove contact info at the end of the article
	paragraphs.Find("span.-newsgate-paragraph-cci-endnote-contact-").Remove()

	paragraphStrings := paragraphs.Map(func(i int, paragraph *gq.Selection) string {
		return strings.TrimSpace(paragraph.Text())
	})

	content := make([]string, 0, len(paragraphStrings)+1)
	if includeTitle {
		title := ExtractTitleFromDocument(doc)
		content = append(content, title)
	}

	content = append(content, withoutEmptyStrings(paragraphStrings)...)

	body := strings.Join(content, "\n")

	return body, nil
}

func ExtractTitleFromDocument(doc *gq.Document) string {
	title := doc.Find("h1[itemprop=headline]")
	return strings.TrimSpace(title.Text())
}

func ExtractBodyFromURL(url string, includeTitle bool) (string, error) {
	logger.Info("Fetching %s ...\n", url)
	doc, err := gq.NewDocument(url)
	if err != nil {
		return "", err
	}

	return extractBodyFromDocument(doc, includeTitle)
}
