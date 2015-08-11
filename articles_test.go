package main

import (
	"fmt"
	"testing"

	"github.com/michigan-com/newsfetch/lib"
)

func TestFormatFeedUrls(t *testing.T) {
	logger.Info("Compile feed urls from sites and sections")

	expected := []string{
		"http://freep.com/feeds/live/sports/json",
		"http://freep.com/feeds/live/news/json",
		"http://freep.com/feeds/live/life/json",
	}

	actual := FormatFeedUrls([]string{"freep.com"}, []string{"sports", "news", "life"})

	if len(expected) != len(actual) {
		t.Errorf("%d != %d", expected, actual)
	}

	for _, eurl := range expected {
		found := false
		for _, aurl := range actual {
			if eurl == aurl {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%s not found in url list", eurl)
		}
	}
}

func TestGetFeedContent(t *testing.T) {
	logger.Info("Download feed and store its content in a Feed struct")

	url := "http://detroitnews.com/feeds/live/news/json"
	feed, err := GetFeedContent(url)
	if err != nil {
		t.Error(err)
	}

	expected := "detroitnews"
	if feed.Site != expected {
		t.Errorf("%s != %s", feed.Site, expected)
	}

	logger.Info("Download feed from the wrong url")

	url = "http://fuckyoubitch.com"
	_, err = GetFeedContent(url)
	if err == nil {
		t.Errorf("%s should fail because it's not a Gannett site", url)
	}
}

func TestParseArticle(t *testing.T) {
	logger.Info("Parse news article")

	url := "http://detroitnews.com/feeds/live/news/json"
	feed, err := GetFeedContent(url)
	if err != nil {
		t.Error(err)
	}

	data, ok := feed.Body.CheckGet("content")
	if !ok {
		t.Error(`"content" property not found in response JSON`)
	}

	contentArr, err := data.Array()
	for i := 0; i < len(contentArr); i++ {
		articleUrl := fmt.Sprintf("http://%s.com%s", feed.Site, data.GetIndex(i).Get("url").MustString())
		article, err := ParseArticle(articleUrl, data.GetIndex(i), false)
		if err != nil {
			continue
		}

		if article.ArticleId == 0 {
			t.Error("ArticleId should not be empty")
		}

		if article.Headline == "" {
			t.Error("Headline should not be empty")
		}

		if article.Url == "" {
			t.Error("Url should not be empty")
		}

		if *article.Photo == (Photo{}) {
			t.Error("Photo should not be empty")
		}
	}

}

func TestRemoveArticles(*testing.T) {}
func TestSaveArticles(*testing.T)   {}

func TestFetchAndParseArticles(t *testing.T) {
	logger.Info("Compile articles and check for duplicate URLs")
	urls := FormatFeedUrls(lib.Sites, lib.Sections)
	articles := FetchAndParseArticles(urls, false)

	artMap := map[string]int{}
	for _, art := range articles {
		if artMap[art.Url] != 0 {
			t.Error("%s duplicate found, this should never happen")
		}
		artMap[art.Url] = 1
	}
}
