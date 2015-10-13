package fetch

import (
	"testing"
	"time"

	m "github.com/michigan-com/newsfetch/model"
)

func TestArticleModelWithPhoto(t *testing.T) {
	url := "http://www.detroitnews.com/story/news/nation/2015/10/13/planned-parenthood-fetal-tissue/73861022/"

	processor := ParseArticleAtURL(url, true)
	t.Log(processor)

	if processor.Err != nil {
		t.Fatalf("Failed to process article: %s", processor.Err)
	}

	article := processor.Article

	validateArticle(article, t)

	if article.Photo == nil {
		t.Log(article.Photo)
		t.Fatal("Article photo should not be empty.")
	}

	if article.Photo.Full == (m.PhotoInfo{}) {
		t.Fatal("Article photo full should not be empty.")
	}

	if article.Photo.Full.Url == "" {
		t.Fatal("Article photo full url should not be empty.")
	}

	if article.Photo.Thumbnail == (m.PhotoInfo{}) {
		t.Fatal("Article photo thumbnail should not be empty.")
	}

	if article.Photo.Thumbnail.Url == "" {
		t.Fatal("Article photo thumbnail url should not be empty.")
	}
}

func TestArticleModelWithoutPhoto(t *testing.T) {
	url := "http://www.battlecreekenquirer.com/story/news/local/2015/10/11/young-boy-scouts-hooked-fishing/73774328/"

	processor := ParseArticleAtURL(url, true)
	t.Log(processor)

	if processor.Err != nil {
		t.Fatalf("Failed to process article: %s", processor.Err)
	}

	article := processor.Article

	validateArticle(article, t)

	if article.Photo != nil {
		t.Log(article.Photo)
		t.Fatal("This article should not contain a photo")
	}
}

func validateArticle(article *m.Article, t *testing.T) {
	if article.ArticleId == 0 {
		t.Fatal("Article Id should not be empty.")
	}

	if article.Headline == "" {
		t.Fatal("Article headline should not be empty.")
	}

	if article.BodyText == "" {
		t.Fatal("Article body should not be empty.")
	}

	if article.Url == "" {
		t.Fatal("Article url should not be empty.")
	}

	if article.Source == "" {
		t.Fatal("Article source should not be empty.")
	}

	if article.Section == "" {
		t.Fatal("Article section should not be empty.")
	}

	if article.Timestamp == (time.Time{}) {
		t.Fatal("Article timestamp should not be empty.")
	}
}
