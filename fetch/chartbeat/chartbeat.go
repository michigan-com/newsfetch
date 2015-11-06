package fetch

import (
	"net/url"
	"net/http"
	"fmt"
	"errors"

	"gopkg.in/mgo.v2"

	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model/chartbeat"
)

const chartbeatApiUrlFormat = "http://api.chartbeat.com/%s/?apikey=%s&host=%s&limit=100"

var chartbeatDebugger = lib.NewCondLogger("newsfetch:fetch:chartbeat")
var chartbeatError = lib.NewCondLogger("newsfetch:fetch:chartbeat:error")

type ChartbeatFetch interface {
	Fetch([]string) m.Snapshot
}

// Types
// Beat type to be called from command package
type Beat interface {
	Run(*mgo.Session, string)
}

// ChartbeatApi to be used in the chartbeat package
type ChartbeatApi struct {
	Url          ChartbeatUrl
	MapiEndpoint string         // https://api.michigan.com/<MapiEndpoint> API endpointto update the sockets
	Fetch        ChartbeatFetch // Interface in fetch/ that will fetch the chartbeat info
}

type ChartbeatUrl struct {
	ChartbeatEndpoint string // https://api.chartbeat.com API endpoint
	ChartbeatParams   string // Urls params for chartbeat url, "" to specify none
}

func (u ChartbeatUrl) Urls(apiKey string) []string {
	urls, err := FormatChartbeatUrls(u.ChartbeatEndpoint, lib.Sites, apiKey)
	if err != nil {
		chartbeatDebugger.Println(err)
	}
	return AddUrlParams(urls, u.ChartbeatParams)
}

func (c ChartbeatApi) Run(session *mgo.Session, apiKey string) {
	urls := c.Url.Urls(apiKey)
	snapshot := c.Fetch.Fetch(urls)

	if session != nil {
		snapshot.Save(session)
	}

	// TODO hit mapi
	_, err := http.Get(fmt.Sprintf("https://api.michigan.com/%s/", c.MapiEndpoint))
	if err != nil {
		chartbeatDebugger.Println("Failed to update mapi")
	}
}

/** The beats */
var TrafficSeriesApi = ChartbeatApi{
	ChartbeatUrl{"historical/traffic/series", ""},
	"historical-traffic",
	TrafficSeries{},
}
var QuickStatsApi = ChartbeatApi{
	ChartbeatUrl{"live/quickstats/v4", "all_platforms=1&loyalty=1"},
	"quickstats",
	Quickstats{},
}

var RecentApi = ChartbeatApi{
	ChartbeatUrl{"live/recent/v3", ""},
	"recent",
	Recent{},
}

var ReferrersApi = ChartbeatApi{
	ChartbeatUrl{"live/referrers/v3", ""},
	"referrers",
	Referrers{},
}

var TopGeoApi = ChartbeatApi{
	ChartbeatUrl{"live/top_geo/v1", ""},
	"topgeo",
	TopGeo{},
}

// TODO add back visits calculations
var TopPagesApi = ChartbeatApi{
	ChartbeatUrl{"live/toppages/v3", "all_platforms=1&loyalty=1"},
	"popular",
	TopPages{},
}

/** End beats */

/*
  Format chartbeat URLs based on a chartbeat API endpoint

  Format: http://api.chartbeat.com/<endPoint>/?apikey=<key>&host=<site[i]>&

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
