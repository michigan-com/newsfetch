package main

import (
	"fmt"
	"github.com/bmizerany/assert"
	"github.com/michigan-com/newsFetch/lib"
	"regexp"
	"testing"
)

// Tests the parsing of a url for the article ID
func TestGetArticleId(t *testing.T) {

	type idTest struct {
		url       string
		returnVal int
	}

	validTestStrings := []idTest{
		idTest{
			url:       "http://www.freep.com/story/sports/mlb/tigers/2015/07/30/detroit-tigers-daniel-norris/30883693/",
			returnVal: 30883693,
		},
		idTest{
			url:       "http://www.freep.com/story/news/local/michigan/oakland/2015/07/29/hoffa-disappearance-anniversary-teamsters/30862419/",
			returnVal: 30862419,
		},
		idTest{
			url:       "asdfasdfasdf/asdfasd/9/",
			returnVal: 9,
		},
		idTest{
			url:       "asdfasdf",
			returnVal: -1,
		},
	}

	for _, test := range validTestStrings {
		returnVal := getArticleId(test.url)

		assert.Equal(t, returnVal, test.returnVal, "Return values don't match")
	}
}

func TestFormatUrls(t *testing.T) {

	// Maps strings to an array of strings
	urlMap := make(map[string][]string)
	sectionMap := make(map[string]bool)

	formattedUrls := formatUrls()
	assert.Equal(t, len(formattedUrls), len(lib.Sites)*len(lib.Sections), "Incorrect number of sites")

	// Create a map of sections for constant lookup time
	for _, section := range lib.Sections {
		sectionMap[section] = true
	}

	// Iterate over all the returned urls and save the site/sections
	for _, val := range formattedUrls {
		hostRegex := regexp.MustCompile("http://(.+[.].{2,3})/feeds/live/(.+)/json")
		match := hostRegex.FindStringSubmatch(val)
		assert.Equal(t, len(match), 3, fmt.Sprintf("Url %s not successfully parsed", val))

		site := match[1]
		section := match[2]

		_, ok := urlMap[site]
		if !ok {
			urlMap[site] = make([]string, 0, len(lib.Sections))
		}
		urlMap[site] = append(urlMap[site], section)
	}

	// Now iterate over all the expected sites/sections and make sure they exist
	for _, site := range lib.Sites {
		sections, ok := urlMap[site]
		assert.Equal(t, ok, true, fmt.Sprintf("Can't find site %s in urlMap", site))

		// Now look at sections
		assert.Equal(t, len(sections), len(lib.Sections))
		for _, section := range sections {
			if section == "life-home" { // Special case
				section = "life"
			}
			_, ok := sectionMap[section]
			assert.Equal(t, ok, true, fmt.Sprintf("Section %s not recognized", section))
		}
	}
}

func TestGetFeedUrl(t *testing.T) {
	url := "http://google.com"
	articles, err := getFeedUrl(url)
	assert.NotEqual(t, err, nil, "Should have an error")
	assert.Equal(t, len(articles), 0, "No articles shold have been returned")
}
