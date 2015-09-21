package lib

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestFormatChartbeatUrls(t *testing.T) {
	t.Log("Testing the formatting of Chartbeat URLs")

	apiKey := "asdf"

	// Test the toppages api
	endPoint := "live/toppages/v3"
	formattedUrls, err := FormatChartbeatUrls(endPoint, Sites, apiKey)

	if err != nil {
		t.Fatalf("%v", err)
	}

	// Check to make sure we have the right numnber of urls
	if len(formattedUrls) != len(Sites) {
		t.Fatalf(fmt.Sprintf("Expected %d urls, got %d", len(Sites), len(formattedUrls)))
	}

	// Test to make sure the URLs formatted correctly
	for i := 0; i < len(formattedUrls); i++ {
		url := formattedUrls[i]
		site := Sites[i]
		if !strings.Contains(url, endPoint) {
			t.Fatalf(fmt.Sprintf("Url %s does not contain endPoint %s", url, endPoint))
		} else if !strings.Contains(url, apiKey) {
			t.Fatalf(fmt.Sprintf("Url %s does not contain the apiKey %s", url, apiKey))
		} else if !strings.Contains(url, site) {
			t.Fatalf(fmt.Sprintf("Url %s should have site %s as a parameter", url, site))
		}
	}

	// Test with no sites
	endPoint = "blah"
	formattedUrls, err = FormatChartbeatUrls(endPoint, []string{}, apiKey)
	if len(formattedUrls) != 0 {
		t.Fatalf(fmt.Sprintf("%d urls created, should have been 0", len(formattedUrls)))
	}

	// Test and make sure that no api key returns an error
	_, err = FormatChartbeatUrls(endPoint, Sites, "")
	if err == nil {
		t.Fatalf("Should have thrown an error when no API key was set")
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
		t.Fatalf(fmt.Sprintf("Url '%s' should have thrown an error", badUrl))
	}
}

func TestSaveSnapshot(t *testing.T) {
	mongoUri := os.Getenv("MONGOURI")
	if mongoUri == "" {
		t.Fatalf("%v", "No mongo URI specified, failing test")
	}

	// Make an article snapshot and save it
	numArticles := 20
	toppages := make([]*TopArticle, 0, numArticles)
	for i := 0; i < numArticles; i++ {
		article := &TopArticle{}
		article.ArticleId = i
		article.Headline = fmt.Sprintf("Article %d", i)
		article.Visits = 100

		toppages = append(toppages, article)
	}

	// Add the collection 4 times
	SaveTopPagesSnapshot(mongoUri, toppages)
	SaveTopPagesSnapshot(mongoUri, toppages)
	SaveTopPagesSnapshot(mongoUri, toppages)
	SaveTopPagesSnapshot(mongoUri, toppages)

	// Now verify
	session := DBConnect(mongoUri)
	defer DBClose(session)

	col := session.DB("").C("Toppages")
	numSnapshots, err := col.Count()

	if err != nil {
		t.Fatalf("%v", err)
	}
	if numSnapshots != 1 {
		t.Fatalf("Should only be one collection")
	}

	snapshot := &Snapshot{}
	err = col.Find(bson.M{}).One(&snapshot)

	if len(snapshot.Articles) != numArticles {
		t.Fatalf("Should be %d values in the snapshot, but there are %d", numArticles, len(snapshot.Articles))
	}
}

func TestSortTopArticles(t *testing.T) {
	// This is a test URL provided by chartbeat
	testUrl := "http://api.chartbeat.com/live/toppages/v3/?apikey=317a25eccba186e0f6b558f45214c0e7&host=gizmodo.com"
	topPages, _ := GetTopPages(testUrl)

	// Compile the top articles together
	topArticles := make([]*TopArticle, 0, len(topPages.Pages))
	for i := 0; i < len(topPages.Pages); i++ {
		page := topPages.Pages[i]
		article := &TopArticle{}
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
