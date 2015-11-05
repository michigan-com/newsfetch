package fetch

import (
	"fmt"
	"net/url"

	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model/chartbeat"
)

var chartbeatDebugger = lib.NewCondLogger("newsfetch:fetch:chartbeat")
var chartbeatError = lib.NewCondLogger("newsfetch:fetch:chartbeat:error")

type ChartbeatFetch interface {
	Fetch([]string) m.Snapshot
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

// Chartbeat queries have a GET parameter "host", which represents the host
// we're getting data on. Pull the host from the url and return it.
// Return host (e.g. freep.com)
// Return "" if we don't find one
func GetHostFromParams(inputUrl string) (string, error) {
	var host string
	var err error

	parsed, err := url.Parse(inputUrl)
	if err != nil {
		return host, err
	}

	hosts := parsed.Query()["host"]
	if len(hosts) > 0 {
		host = hosts[0]
	}

	return host, err
}
