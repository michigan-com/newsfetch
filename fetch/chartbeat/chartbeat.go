package fetch

import (
	"net/url"

	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model/chartbeat"
)

var chartbeatDebugger = lib.NewCondLogger("newsfetch:fetch:chartbeat")
var chartbeatError = lib.NewCondLogger("newsfetch:fetch:chartbeat:error")

type ChartbeatFetch interface {
	Fetch([]string) m.Snapshot
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
