package lib

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"
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

	data, ok := feed.Body["content"].([]interface{})
	if !ok {
		t.Error(`"content" property not found in response JSON`)
	}

	for _, jso := range data {
		articleJson := jso.(map[string]interface{})
		url := articleJson["url"].(string)
		articleUrl := fmt.Sprintf("http://%s.com%s", feed.Site, url)
		article, err := ParseArticle(articleUrl, articleJson, false)

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

func TestRemoveArticles(t *testing.T) {}
func TestNoDuplcateArticlesMongo(t *testing.T) {
	logger.Info("Ensure no duplicate articles are in the database")

	uri := "mongodb://localhost:27017/mapi_test"
	session := DBConnect(uri)
	defer DBClose(session)

	RemoveArticles(uri)

	articleCol := session.DB("").C("Article")

	urls := FormatFeedUrls([]string{"freep.com"}, []string{"news"})
	for i := 0; i < 2; i++ {
		articles := FetchAndParseArticles(urls, false)
		SaveArticles(uri, articles)
	}

	var arts []Article
	err := articleCol.Find(bson.M{}).All(&arts)
	if err != nil {
		t.Error("No articles found in database, should not happen after adding them")
	}

	artMap := map[string]int{}
	for _, art := range arts {
		if artMap[art.Url] != 0 {
			t.Errorf("Article %s already exists in database, should never happen", art.Url)
		}
		artMap[art.Url] = 1
	}
}

func TestFetchAndParseArticles(t *testing.T) {
	logger.Info("Compile articles and check for duplicate URLs")
	urls := FormatFeedUrls(Sites, Sections)
	articles := FetchAndParseArticles(urls, false)

	artMap := map[string]int{}
	for _, art := range articles {
		if artMap[art.Url] != 0 {
			t.Error("%s duplicate found, this should never happen")
		}
		artMap[art.Url] = 1
	}
}
