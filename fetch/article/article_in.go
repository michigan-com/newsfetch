package fetch

import (
	"fmt"
	"regexp"
	"time"

	gq "github.com/PuerkitoBio/goquery"

	e "github.com/michigan-com/newsfetch/extraction/body_parsing"
	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model"
)

type ArticleIn struct {
	Site      string
	Url       string
	ArticleId int
	Doc       *gq.Document
}

func NewArticleIn(url string) *ArticleIn {
	article := &ArticleIn{Url: url}
	if article.isBlacklisted() {
		return nil
	}

	return article
}

func (a *ArticleIn) String() string {
	return fmt.Sprintf("<ArticleIn Site: %s, Id: %d, Url: %s>", a.Site, a.ArticleId, a.Url)
}

func (a *ArticleIn) GetData() error {

	artDebugger.Println("Fetching: ", a.Url)

	doc, err := gq.NewDocument(a.Url)
	if err != nil {
		return err
	}

	a.Site, _ = lib.GetHost(a.Url)
	a.Doc = doc
	a.ArticleId = lib.GetArticleId(a.Url)

	return nil
}

func (a *ArticleIn) IsValid() bool {
	if a.Doc == nil {
		artDebugger.Println("Article struct missing ...")
		return false
	}

	if a.ArticleId == 0 {
		artDebugger.Println("Article ID missing ...")
		return false
	}

	return true
}

func (a *ArticleIn) isBlacklisted() bool {
	return lib.IsBlacklisted(a.Url)
}

func (a *ArticleIn) Process(article *m.Article) error {

	extractedSection := e.ExtractSectionInfo(a.Doc)

	sections := make([]string, len(extractedSection.Sections))
	for i, section := range extractedSection.Sections {
		sections[i] = section
	}

	article.Source = a.Site
	article.ArticleId = a.ArticleId
	article.Headline = e.ExtractTitleFromDocument(a.Doc)
	article.Subheadline = e.ExtractSubheadlineFromDocument(a.Doc)
	article.Section = extractedSection.Section
	article.Subsection = extractedSection.Subsection
	article.Sections = sections
	article.Created_at = time.Now()
	article.Updated_at = time.Now()
	article.Timestamp = e.ExtractTimestamp(a.Doc)
	article.Url = a.Url
	article.Photo = e.ExtractPhotoInfo(a.Doc)

	return nil
}

func (a *ArticleIn) GetSiteFromHost(host string) (string, error) {
	replace := regexp.MustCompile("^w{3}[.](.+)[.].+$")
	match := replace.FindStringSubmatch(host)
	if len(match) < 2 {
		return "", fmt.Errorf("Could not parse %s for host", host)
	}

	return match[1], nil
}
