package chartbeat

import (
	"fmt"
	"strings"
	"testing"

	"github.com/michigan-com/newsfetch/lib"
)

func TestFormatChartbeatUrls(t *testing.T) {
	t.Log("Testing the formatting of Chartbeat URLs")

	apiKey := "asdf"

	// Test the toppages api
	endPoint := "live/toppages/v3"
	formattedUrls, err := FormatChartbeatUrls(endPoint, lib.Sites, apiKey)

	if err != nil {
		t.Fatalf("%v", err)
	}

	// Check to make sure we have the right numnber of urls
	if len(formattedUrls) != len(lib.Sites) {
		t.Fatalf("Expected %d urls, got %d", len(lib.Sites), len(formattedUrls))
	}

	// Test to make sure the URLs formatted correctly
	for i := 0; i < len(formattedUrls); i++ {
		url := formattedUrls[i]
		site := lib.Sites[i]
		if !strings.Contains(url, endPoint) {
			t.Fatalf(fmt.Sprintf("Url %s does not contain endPoint %s", url, endPoint))
		} else if !strings.Contains(url, apiKey) {
			t.Fatalf(fmt.Sprintf("Url %s does not contain the apiKey %s", url, apiKey))
		} else if !strings.Contains(url, site) {
			t.Fatalf(fmt.Sprintf("Url %s should have site %s as a parameter", url, site))
		}
	}

	// Add some url params
	urlString := "this=1234&that=abcd&other=what"
	formattedUrls = AddUrlParams(formattedUrls, urlString)
	for _, url := range formattedUrls {
		if !strings.HasSuffix(url, urlString) {
			t.Fatalf("Url %s should end with %s", url, urlString)
		}
		t.Log("Url %s checks out", url)
	}

	// Test with no sites
	endPoint = "blah"
	formattedUrls, err = FormatChartbeatUrls(endPoint, []string{}, apiKey)
	if len(formattedUrls) != 0 {
		t.Fatalf(fmt.Sprintf("%d urls created, should have been 0", len(formattedUrls)))
	}

	// Test and make sure that no api key returns an error
	_, err = FormatChartbeatUrls(endPoint, lib.Sites, "")
	if err == nil {
		t.Fatalf("Should have thrown an error when no API key was set")
	}
}
