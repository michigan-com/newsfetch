package lib

import (
	"strings"

	gq "github.com/PuerkitoBio/goquery"
)

var log = GetLogger()

func withoutEmptyStrings(strings []string) []string {
	result := make([]string, 0, len(strings))
	for _, el := range strings {
		if el != "" {
			result = append(result, el)
		}
	}
	return result
}

func extractBodyFromDocument(doc *gq.Document) (string, error) {
	paragraphs := doc.Find("div[itemprop=articleBody] p")

	// remove contact info at the end of the article
	paragraphs.Find("span.-newsgate-paragraph-cci-endnote-contact-").Remove()

	paragraphStrings := paragraphs.Map(func(i int, paragraph *gq.Selection) string {
		return strings.TrimSpace(paragraph.Text())
	})

	paragraphStrings = withoutEmptyStrings(paragraphStrings)

	body := strings.Join(paragraphStrings, "\n\n")

	return body, nil
}

func ExtractBodyFromURL(url string) (string, error) {
	log.Info("Fetching %s...\n", url)
	doc, err := gq.NewDocument(url)
	if err != nil {
		return "", err
	}

	return extractBodyFromDocument(doc)
}
