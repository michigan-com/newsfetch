package lib

import (
	"fmt"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func TestFormatFeedUrls(t *testing.T) {
	t.Log("Compile feed urls from sites and sections")

	expected := []string{
		"http://freep.com/feeds/live/sports/json",
		"http://freep.com/feeds/live/news/json",
		"http://freep.com/feeds/live/life/json",
	}

	actual := FormatFeedUrls([]string{"freep.com"}, []string{"sports", "news", "life"})

	if len(expected) != len(actual) {
		t.Fatalf("%d != %d", expected, actual)
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
			t.Fatalf("%s not found in url list", eurl)
		}
	}
}

func TestGetFeedContent(t *testing.T) {
	t.Log("Download feed and store its content in a Feed struct")

	url := "http://detroitnews.com/feeds/live/news/json"
	feed, err := GetFeedContent(url)
	if err != nil {
		t.Fatal(err)
	}

	expected := "detroitnews"
	if feed.Site != expected {
		t.Fatalf("%s != %s", feed.Site, expected)
	}

	t.Log("Download feed from the wrong url")

	url = "http://fuckyoubitch.com"
	_, err = GetFeedContent(url)
	if err == nil {
		t.Fatalf("%s should fail because it's not a Gannett site", url)
	}
}

func TestParseArticle(t *testing.T) {
	t.Log("Parse news article")

	url := "http://detroitnews.com/feeds/live/news/json"
	feed, err := GetFeedContent(url)
	if err != nil {
		t.Fatal(err)
	}

	if feed.Body.Content == nil {
		t.Fatal(`"content" property not found in response JSON`)
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
			t.Fatal("ArticleId should not be empty")
		}

		if article.Headline == "" {
			t.Fatal("Headline should not be empty")
		}

		if article.Url == "" {
			t.Fatal("Url should not be empty")
		}

		if article.Photo == nil {
			t.Fatal("Photo should not be empty")
		}

		if article.Photo.Full.Width != 0 || article.Photo.Full.Height != 0 {
			foundFullPhotoDim = true
		}

		if article.Photo.Thumbnail.Width != 0 || article.Photo.Thumbnail.Height != 0 {
			foundThumbPhotoDim = true
		}
	}

	if !foundThumbPhotoDim {
		t.Fatal("Could not find a single thumbnail image width or height dimension")
	}

	if !foundFullPhotoDim {
		t.Fatal("Could not find a single full image width or height dimension")
	}
}

func TestArticlesMongo(t *testing.T) {
	t.Log("Ensure no duplicate articles are in the database")

	uri := "mongodb://localhost:27017/mapi_test"

	RemoveArticles(uri)

	urls := FormatFeedUrls([]string{"freep.com"}, []string{"news"})
	getBody := false
	created_at := time.Now()
	for i := 0; i < 2; i++ {
		/*if i == 0 {
			getBody = true
		}*/
		articles := FetchAndParseArticles(urls, getBody)
		if i == 0 {
			for _, art := range articles {
				art.Headline = ""
				art.Created_at = created_at
				t.Log(art.Created_at)
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
		t.Fatal("No articles found in database, should not happen after adding them")
	}

	artMap := map[string]int{}
	for _, art := range arts {
		if artMap[art.Url] != 0 {
			t.Fatalf("Article %s already exists in database, should never happen", art.Url)
		}
		artMap[art.Url] = 1

		t.Log("Determine if headline gets updated")
		if art.Headline == "" {
			t.Fatalf("The article's headline did not get properly updated")
		}

		t.Log("Created_at should never be updated")
		if art.Created_at.Sub(created_at) > time.Second {
			t.Logf("expected: %s, actual: %s", created_at.String(), art.Created_at.String())
			t.Log(art.Headline)
			t.Fatalf("Created_at was updated, which should not happen")
		}

		/*if art.BodyText == "" {
			t.Fatalf("Body text was overwritten on article update")
		}*/
	}
}

func TestFetchAndParseArticles(t *testing.T) {
	t.Log("Compile articles and check for duplicate URLs")
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
