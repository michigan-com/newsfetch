package fetch

import (
	"errors"
	"fmt"
	"strings"

	"github.com/michigan-com/newsfetch/extraction"
	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model"
)

var artDebugger = lib.NewCondLogger("newsfetch:fetch:article")

// Processor object that contains the Article to be saved
// as well as the body text
type ArticleProcess struct {
	*m.Article
	*m.ExtractedBody
	Html string
	Err  error
}

func (p *ArticleProcess) String() string {
	return fmt.Sprintf("<ArticleProcess %s\n %s\n Error: %v>\n", p.Article, p.ExtractedBody, p.Err)
}

// Primary entry point to process an article's json
// based on the article Url
func ParseArticleAtURL(articleUrl string, runExtraction bool) *ArticleProcess {
	processor := &ArticleProcess{}
	article := &m.Article{}

	articleIn := NewArticleIn(articleUrl)
	if articleIn == nil {
		processor.Err = fmt.Errorf("Article Url was blacklisted")
		return processor
	}

	err := articleIn.GetData()
	if err != nil {
		processor.Err = err
		return processor
	}

	if !articleIn.IsValid() {
		artDebugger.Println("Article is not valid: ", article)
		processor.Err = errors.New("Article is not valid: " + articleUrl)
		return processor
	}

	err = articleIn.Process(article)
	if err != nil {
		artDebugger.Println("Article could not be processed: %s", articleIn)
	}

	var html string
	var bodyExtract *m.ExtractedBody
	if runExtraction {
		bodyExtract = extraction.ExtractDataFromDocument(articleIn.Doc, articleIn.Url, false, false)

		if bodyExtract.Text != "" {
			artDebugger.Printf(
				"Extracted extracted contains %d characters, %d paragraphs. for article %s",
				len(strings.Split(bodyExtract.Text, "")),
				len(strings.Split(bodyExtract.Text, "\n\n")),
				articleIn.Url,
			)
			article.BodyText = bodyExtract.Text
		}
	}

	processor.Article = article
	processor.ExtractedBody = bodyExtract
	processor.Html = html

	return processor
}
