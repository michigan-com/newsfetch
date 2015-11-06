package fetch

import (
	"fmt"
	"os"
	"testing"

	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model"
	mc "github.com/michigan-com/newsfetch/model/chartbeat"
	"gopkg.in/mgo.v2/bson"
)

func TestSaveTimeInterval(t *testing.T) {
	t.Skip("No Mongo tests allowed")
	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri == "" {
		t.Fatalf("No MONGO_URI env variable set")
	}

	session := lib.DBConnect(mongoUri)
	defer lib.DBClose(session)

	// Save a bunch of articles
	numArticles := 20
	articles := make([]*m.Article, 0, numArticles)
	for i := 0; i < numArticles; i++ {
		article := &m.Article{}
		article.ArticleId = i + 1
		articles = append(articles, article)
	}
	/*err := a.SaveArticles(mongoUri, articles)
	if err != nil {
		t.Fatalf("%v", err)
	}*/

	// Calculate a bunch of the time intervals
	topPages := make([]*mc.TopArticle, 0, numArticles)
	visits := map[int]int{}
	for i := 0; i < numArticles; i++ {
		article := &mc.TopArticle{}
		articleId := i + 1
		numVisits := lib.RandomInt(500)

		article.ArticleId = articleId
		article.Visits = numVisits

		visits[articleId] = numVisits
		topPages = append(topPages, article)
	}
	CalculateTimeInterval(topPages, session)

	// Now check the articles saved and make sure they updated the visits
	savedArticles := make([]*m.Article, 0, numArticles)
	articleCol := session.DB("").C("Article")
	articleCol.Find(bson.M{
		"article_id": bson.M{
			"$gte": 1,
			"$lte": 20,
		},
	}).All(&savedArticles)

	if len(savedArticles) != numArticles {
		t.Fatalf("Failed to get the right number of articles from the DB")
	}

	// Verify the visits match up
	for _, article := range savedArticles {
		articleId := article.ArticleId
		if len(article.Visits) != 1 {
			t.Fatalf("Should be exactly one visit in the array, there are %d", len(article.Visits))
		}

		numVisits, ok := visits[articleId]
		if !ok {
			t.Fatalf("Failed to find visits in map for articleId=%d", articleId)
		}

		if numVisits != article.Visits[0].Max {
			t.Fatalf("Article %d: Expected: %d visits. Actual: %d visits", articleId, numVisits, article.Visits[0].Max)
		}
	}
}

func TestCalculateTimeInterval(t *testing.T) {
	articles := []*m.Article{&m.Article{}, &m.Article{}, &m.Article{}, &m.Article{}}
	articles[0].ArticleId = 1
	articles[1].ArticleId = 2
	articles[2].ArticleId = 3
	articles[3].ArticleId = 4

	articleVisits := map[int]int{
		1: 100,
		2: 600,
		3: 12,
		4: 566,
	}

	calculateTimeInterval(articles, articleVisits)

	for _, article := range articles {
		id := article.ArticleId
		visits := article.Visits[0].Max
		expectedVisits := articleVisits[id]

		if visits != expectedVisits {
			t.Fatalf("Expected %d visits, got %d", expectedVisits, visits)
		}
	}
}

func TestGetTopPages(t *testing.T) {
	// This is a test URL provided by chartbeat
	testUrl := "http://api.chartbeat.com/live/toppages/v3/?apikey=317a25eccba186e0f6b558f45214c0e7&host=gizmodo.com"
	topPages, err := GetTopPages(testUrl)

	if err != nil {
		t.Fatalf("%v", err)
	}
	if len(topPages.Pages) == 0 {
		t.Fatalf("No pages got returned in the test API call")
	}

	// This is a url that should return an error
	badUrl := "this is not a real url"
	_, err = GetTopPages(badUrl)
	if err == nil {
		t.Fatalf("Url '%s' should have thrown an error", badUrl)
	}

	// This is a valid url that should return an error
	badUrl = "http://google.com"
	_, err = GetTopPages(badUrl)
	if err == nil {
		t.Fatalf("Url '%s' should have thrown an error", badUrl)
	}
}

func TestSortTopArticles(t *testing.T) {
	// This is a test URL provided by chartbeat
	testUrl := "http://api.chartbeat.com/live/toppages/v3/?apikey=317a25eccba186e0f6b558f45214c0e7&host=gizmodo.com"
	topPages, _ := GetTopPages(testUrl)

	// Compile the top articles together
	topArticles := make([]*mc.TopArticle, 0, len(topPages.Pages))
	for i := 0; i < len(topPages.Pages); i++ {
		page := topPages.Pages[i]
		article := &mc.TopArticle{}
		article.Visits = page.Stats.Visits
		article.Url = page.Path
		article.Headline = page.Title

		topArticles = append(topArticles, article)
	}

	// Now sort them and check
	sorted := SortTopArticles(topArticles)
	var lastVal int = -1
	for i := 0; i < len(sorted); i++ {
		thisVal := sorted[i].Visits
		if lastVal == -1 {
			lastVal = sorted[i].Visits
			continue
		}

		// Fail if we have a value that is larger than the preceding value
		if thisVal > lastVal {
			t.Fatalf(fmt.Sprintf("sorted[%d] == %d, sorted[%d] == %d. Should be sorted in descending order", i-1, lastVal, i, thisVal))
		}
		lastVal = thisVal
	}
}
