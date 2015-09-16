package lib

import (
	"fmt"
	"strings"
	"testing"
)

func TestFormatChartbeatUrls(t *testing.T) {
	t.Log("Testing the formatting of Chartbeat URLs")

	apiKey := "asdf"

	// Test the toppages api
	endPoint := "live/toppages/v3"
	formattedUrls, err := FormatChartbeatUrls(endPoint, Sites, apiKey)

	if err != nil {
		panic(err)
	}

	// Check to make sure we have the right numnber of urls
	if len(formattedUrls) != len(Sites) {
		panic(fmt.Sprintf("Expected %d urls, got %d", len(Sites), len(formattedUrls)))
	}

	// Test to make sure the URLs formatted correctly
	for i := 0; i < len(formattedUrls); i++ {
		url := formattedUrls[i]
		site := Sites[i]
		if !strings.Contains(url, endPoint) {
			panic(fmt.Sprintf("Url %s does not contain endPoint %s", url, endPoint))
		} else if !strings.Contains(url, apiKey) {
			panic(fmt.Sprintf("Url %s does not contain the apiKey %s", url, apiKey))
		} else if !strings.Contains(url, site) {
			panic(fmt.Sprintf("Url %s should have site %s as a parameter", url, site))
		}
	}

	// Test with no sites
	endPoint = "blah"
	formattedUrls, err = FormatChartbeatUrls(endPoint, []string{}, apiKey)
	if len(formattedUrls) != 0 {
		panic(fmt.Sprintf("%d urls created, should have been 0", len(formattedUrls)))
	}

	// Test and make sure that no api key returns an error
	_, err = FormatChartbeatUrls(endPoint, Sites, "")
	if err == nil {
		panic("Should have thrown an error when no API key was set")
	}

}

func TestGetTopPages(t *testing.T) {
	// This is a test URL provided by chartbeat
	testUrl := "http://api.chartbeat.com/live/toppages/v3/?apikey=317a25eccba186e0f6b558f45214c0e7&host=gizmodo.com"
	topPages, err := GetTopPages(testUrl)

	if err != nil {
		panic(err)
	}
	if len(topPages.Pages) == 0 {
		panic("No pages got returned in the test API call")
	}

	// This is a url that should return an error
	badUrl := "this is not a real url"
	_, err = GetTopPages(badUrl)

	if err == nil {
		panic(fmt.Sprintf("Url '%s' should have thrown an error", badUrl))
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
			panic(fmt.Sprintf("sorted[%d] == %d, sorted[%d] == %d. Should be sorted in descending order", i-1, lastVal, i, thisVal))
		}
		lastVal = thisVal
	}
}
