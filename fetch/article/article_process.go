package fetch

import (
	"errors"
	"strings"

	"github.com/michigan-com/newsfetch/extraction"
	m "github.com/michigan-com/newsfetch/model"
)

func ParseArticleAtURL(articleUrl string, runExtraction bool) (*m.Article, string, *m.ExtractedBody, error) {
	article := &m.Article{}

	articleIn := NewArticleIn(articleUrl)
	err := articleIn.GetData()

	if err != nil {
		return nil, "", nil, err
	}

	if !articleIn.IsValid() {
		artDebugger.Println("Article is not valid: ", article)
		return nil, "", nil, errors.New("Article is not valid: " + articleUrl)
	}

	err = articleIn.Process(article)
	if err != nil {
		artDebugger.Println("Article could not be processed: %s", articleIn)
	}

	var html string
	var bodyExtract *m.ExtractedBody
	if runExtraction {
		html = articleIn.BodyHTML()
		bodyExtract = extraction.ExtractDataFromHTMLString(html, articleUrl, false)

		if bodyExtract.Text != "" {
			artDebugger.Printf(
				"Extracted extracted contains %d characters, %d paragraphs.",
				len(strings.Split(bodyExtract.Text, "")),
				len(strings.Split(bodyExtract.Text, "\n\n")),
			)
			article.BodyText = bodyExtract.Text
		}
	}

	return article, html, bodyExtract, nil
}
