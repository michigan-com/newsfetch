package fetch

import (
	"errors"
	"fmt"
	"github.com/michigan-com/newsfetch/lib"
)

var chartbeatDebugger = lib.NewCondLogger("fetch:chartbeat")
var chartbeatError = lib.NewCondLogger("fetch:chartbeat:error")

const chartbeatApiUrlFormat = "http://api.chartbeat.com/%s/?apikey=%s&host=%s&limit=100"

/*
	Format chartbeat URLs based on a chartbeat API endpoint

	Format: http://api.chartbeat.com/<endPoint>/?apikey=<key>&host=<site[i]>

	Example endPoint (NOTE no starting or ending slashes): live/toppages/v3
*/
func FormatChartbeatUrls(endPoint string, sites []string, apiKey string) ([]string, error) {
	urls := make([]string, 0, len(sites))

	if apiKey == "" {
		return urls, errors.New(fmt.Sprintf("No API key specified. Use the -k flag to specify (Run ./newsfetch chartbeat --help for more info)"))
	}

	for i := 0; i < len(sites); i++ {
		site := sites[i]

		url := fmt.Sprintf(chartbeatApiUrlFormat, endPoint, apiKey, site)

		urls = append(urls, url)
	}

	return urls, nil
}

/*
	Add query string to the end of the each url in an array of urls.
	Expects that some url params are already added

	AddUrlParam(["http://google.com?test=123", "http://yahoo.com?test=abc"], "test2=added")

	Result:

		["http://google.com?test=123&test2=added", "http://yahoo.com?test=abc&test2=added"]
*/
func AddUrlParams(urls []string, queryString string) []string {
	for i := 0; i < len(urls); i++ {
		urls[i] = fmt.Sprintf("%s&%s", urls[i], queryString)
	}
	return urls
}
