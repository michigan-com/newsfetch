package lib

import (
	"strings"
	"testing"
	//"time"
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
		t.Fatalf("Expected %d urls, got %d", len(Sites), len(formattedUrls))
	}

	// Test to make sure the URLs formatted correctly
	for i := 0; i < len(formattedUrls); i++ {
		url := formattedUrls[i]
		site := Sites[i]
		if !strings.Contains(url, endPoint) {
			t.Fatalf("Url %s does not contain endPoint %s", url, endPoint)
		} else if !strings.Contains(url, apiKey) {
			t.Fatalf("Url %s does not contain the apiKey %s", url, apiKey)
		} else if !strings.Contains(url, site) {
			t.Fatalf("Url %s should have site %s as a parameter", url, site)
		}
	}

	// Test with no sites
	endPoint = "blah"
	formattedUrls, err = FormatChartbeatUrls(endPoint, []string{}, apiKey)
	if len(formattedUrls) != 0 {
		t.Fatalf("%d urls created, should have been 0", len(formattedUrls))
	}

	// Test and make sure that no api key returns an error
	_, err = FormatChartbeatUrls(endPoint, Sites, "")
	if err == nil {
		t.Fatalf("Should have thrown an error when no API key was set")
	}
}

func TestCalculateTimeInterval(t *testing.T) {
	articles := []*Article{&Article{}}
	articles[0].ArticleId = 1
	//articles[0].Visits = []TimeInterval{
	//TimeInterval{
	//10,
	//time.Now(),
	//},
	//}

	articleVisits := map[int]int{
		1: 100,
	}

	calculateTimeInterval(articles, articleVisits)

	if articles[0].Visits[0].Max != 100 {
		t.Fatalf("Should be 100, actual %d", articles[0].Visits[0].Max)
	}
	logger.Info("%v", articles[0].Visits)
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
}

func TestSaveSnapshot(t *testing.T) {
	//TODO figure out a test DB situation and how to read that mongoUri into the
	// test cases
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
			t.Fatalf("sorted[%d] == %d, sorted[%d] == %d. Should be sorted in descending order", i-1, lastVal, i, thisVal)
		}
		lastVal = thisVal
	}
}
