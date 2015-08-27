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

	if feed.Body.Content == nil {
		t.Error(`"content" property not found in response JSON`)
	}

	foundFullPhotoDim := false
	foundThumbPhotoDim := false
	for _, articleJson := range feed.Body.Content {
		url := articleJson.Url
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

		if article.Photo == nil {
			t.Error("Photo should not be empty")
		}

		if article.Photo.Full.Width != 0 || article.Photo.Full.Height != 0 {
			foundFullPhotoDim = true
		}

		if article.Photo.Thumbnail.Width != 0 || article.Photo.Thumbnail.Height != 0 {
			foundThumbPhotoDim = true
		}
	}

	if !foundThumbPhotoDim {
		t.Error("Could not find a single thumbnail image width or height dimension")
	}

	if !foundFullPhotoDim {
		t.Error("Could not find a single full image width or height dimension")
	}
}

func TestRemoveArticles(t *testing.T) {}
func TestArticlesMongo(t *testing.T) {
	logger.Info("Ensure no duplicate articles are in the database")

	uri := "mongodb://localhost:27017/mapi_test"

	RemoveArticles(uri)

	urls := FormatFeedUrls([]string{"freep.com"}, []string{"news"})
	getBody := false
	for i := 0; i < 2; i++ {
		/*if i == 0 {
			getBody = true
		}*/
		articles := FetchAndParseArticles(urls, getBody)

		if i == 0 {
			for _, art := range articles {
				art.Headline = ""
			}
		}

		SaveArticles(uri, articles)
	}

	session := DBConnect(uri)
	defer DBClose(session)
	articleCol := session.DB("").C("Article")

	var arts []Article
	err := articleCol.Find(bson.M{}).All(&arts)
	if err != nil {
		t.Error("No articles found in database, should not happen after adding them")
	}

	artMap := map[string]int{}
	for _, art := range arts {
		if artMap[art.Url] != 0 {
			t.Fatalf("Article %s already exists in database, should never happen", art.Url)
		}
		artMap[art.Url] = 1

		logger.Info("Determine if headline gets updated")
		if art.Headline == "" {
			t.Fatalf("The article's headline did not get properly updated")
		}

		/*if art.BodyText == "" {
			t.Errorf("Body text was overwritten on article update")
		}*/
	}
}

func TestFetchAndParseArticles(t *testing.T) {
	logger.Info("Compile articles and check for duplicate URLs")
	urls := FormatFeedUrls(Sites, Sections)
	articles := FetchAndParseArticles(urls, false)

	artMap := map[string]int{}
	for _, art := range articles {
		if artMap[art.Url] != 0 {
			t.Fatal("%s duplicate found, this should never happen")
		}
		artMap[art.Url] = 1
	}
}
