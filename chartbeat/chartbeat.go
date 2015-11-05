package chartbeat

import (
	"errors"
	"fmt"

	"gopkg.in/mgo.v2"

	f "github.com/michigan-com/newsfetch/fetch/chartbeat"
	"github.com/michigan-com/newsfetch/lib"
)

const chartbeatApiUrlFormat = "http://api.chartbeat.com/%s/?apikey=%s&host=%s&limit=100"

var debugger = lib.NewCondLogger("newsfetch:chartbeat")

// Types
// Beat type to be called from command package
type Beat interface {
	Run(*mgo.Session, string)
}

// ChartbeatApi to be used in the chartbeat package
type ChartbeatApi struct {
	Endpoint string
	Fetch    f.ChartbeatFetch
}

func (c ChartbeatApi) Urls(apiKey string) []string {
	urls, err := FormatChartbeatUrls(c.Endpoint, lib.Sites, apiKey)
	if err != nil {
		debugger.Println(err)
	}
	return urls
}

func (c ChartbeatApi) Run(session *mgo.Session, apiKey string) {
	urls := c.Urls(apiKey)
	snapshot := c.Fetch.Fetch(urls)

	if session != nil {
		snapshot.Save(session)
	}

	// TODO hit mapi
}

/** The beats */
var Historical = ChartbeatApi{
	"historical/traffic/series",
	f.Historical{},
}

var QuickStats = ChartbeatApi{
	"live/quickstats/v4",
	f.Quickstats{},
}

var Recent = ChartbeatApi{
	"live/recent/v3",
	f.Recent{},
}

var Referrers = ChartbeatApi{
	"live/referrers/v3",
	f.Referrers{},
}

var TopGeo = ChartbeatApi{
	"live/top_geo/v1",
	f.TopGeo{},
}

var TopPages = ChartbeatApi{
	"live/toppages/v3",
	f.TopPages{},
}

/** End beats */

/*
  Format chartbeat URLs based on a chartbeat API endpoint

  Format: http://api.chartbeat.com/<endPoint>/?apikey=<key>host=<site[i]>

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
