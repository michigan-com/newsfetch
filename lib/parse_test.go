package lib

import (
	"fmt"
	"testing"
)

type UrlIds struct {
	Url string
	Id  int
}

type UrlHosts struct {
	Url  string
	Host string
}

func TestGetArticleId(t *testing.T) {
	t.Log("Testing the parsing of article IDs from urls")

	testCases := []UrlIds{{
		"http://www.freep.com/story/money/2015/09/16/uaw-fca-showdown/32488047/",
		32488047,
	}, {
		"http://www.freep.com/story/sports/nfl/lions/2015/09/15/melvin-ingram-matthew-stafford-detroit-lions/72323384/",
		72323384,
	}, {
		"http://www.freep.com/story/sports/college/university-michigan/wolverines/2015/09/15/michigan-wolverines-offensive-line/72292760",
		72292760,
	}, {
		"http://www.detroitnews.com/story/entertainment/arts/2015/09/16/dia-names-new-director/32497959/",
		32497959,
	}, {
		"this should return -1",
		-1,
	}}

	for i := 0; i < len(testCases); i++ {
		testCase := testCases[i]

		id := GetArticleId(testCase.Url)
		if id != testCase.Id {
			panic(fmt.Sprintf("Url %s should have generated ID %d, instead it generated %d", testCase.Url, testCase.Id, id))
		}
	}

}

func TestGetHost(t *testing.T) {
	t.Log("testing the parsing of a host in a url string")

	// These are valid test cases that should return a host
	testCases := []UrlHosts{{
		"http://google.com",
		"google",
	}, {
		"http://freep.com/sports/lions/",
		"freep",
	}, {
		"http://www.detroitnews.com/story/entertainment/arts/2015/09/16/dia-names-new-director/32497959/",
		"detroitnews",
	}, {
		"http://this.that.domain.subdomain.wtf.freep.com/this/is/some/path",
		"freep",
	}}

	for i := 0; i < len(testCases); i++ {
		testCase := testCases[i]
		url := testCase.Url
		expected := testCase.Host

		result, err := GetHost(url)
		if err != nil {
			panic(err)
		} else if result != expected {
			panic(fmt.Sprintf("Url %s should have host %s. Instead got %s", url, expected, result))
		}
	}

	// These are invalid test cases that should return an error
	errorCases := []string{
		"invalid url",
		"freep.this.com/asdfasdf",
		"1231231231",
		"----------123what?",
		"http:::::asdfasdf",
	}

	for i := 0; i < len(errorCases); i++ {
		testCase := errorCases[i]
		result, err := GetHost(testCase)

		if err == nil {
			panic(fmt.Sprintf("Test case '%s' should have returned an error", testCase))
		} else if result != "" {
			panic(fmt.Sprintf("Test case '%s' should have return '' for a host", testCase))
		}
	}
}
