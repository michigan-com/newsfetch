package commands

import (
	"net/http"
	"sync"
	"time"

	f "github.com/michigan-com/newsfetch/fetch/chartbeat"
	"github.com/michigan-com/newsfetch/lib"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2"
)

// Beats
type Beat interface {
	Run(*mgo.Session)
}

type TopPages struct{}
type QuickStats struct{}
type TopGeo struct{}

var chartbeatDebugger = lib.NewCondLogger("chartbeat")

var cmdChartbeat = &cobra.Command{
	Use:   "chartbeat",
	Short: "Hit the Chartbeat API",
}

var cmdAllBeats = &cobra.Command{
	Use:   "all",
	Short: "Fetch all Chartbeat Beats (toppages, quickstats)",
	Run: func(cmd *cobra.Command, argv []string) {
		RunChartbeatCommands([]Beat{&TopPages{}, &QuickStats{}, &TopGeo{}})
	},
}

var cmdTopPages = &cobra.Command{
	Use:   "toppages",
	Short: "Fetch toppages snapshot for Chartbeat",
	Run: func(cmd *cobra.Command, argv []string) {
		RunChartbeatCommands([]Beat{&TopPages{}})
	},
}

var cmdQuickStats = &cobra.Command{
	Use:   "quickstats",
	Short: "Fetch quickstats snapshot for Chartbeat",
	Run: func(cmd *cobra.Command, argv []string) {
		RunChartbeatCommands([]Beat{&QuickStats{}})
	},
}

var cmdTopGeo = &cobra.Command{
	Use:   "topgeo",
	Short: "Fetch topgeo snapshot for Chartbeat",
	Run: func(cmd *cobra.Command, argv []string) {
		RunChartbeatCommands([]Beat{&TopGeo{}})
	},
}

func RunChartbeatCommands(beats []Beat) {
	// Set up environment
	var session *mgo.Session
	if globalConfig.MongoUrl != "" {
		session = lib.DBConnect(globalConfig.MongoUrl)
		defer lib.DBClose(session)
	}

	for {
		startTime := time.Now()

		// Run the actual meat of the program
		var beatWait sync.WaitGroup

		for _, beat := range beats {
			beatWait.Add(1)
			copy := session.Copy()
			defer copy.Close()

			go func(beat Beat) {
				beat.Run(copy)
				beatWait.Done()
			}(beat)
		}

		beatWait.Wait()

		getElapsedTime(&startTime)

		if loop != -1 {
			chartbeatDebugger.Printf("Looping! Sleeping for %d seconds...", loop)
			time.Sleep(time.Duration(loop) * time.Second)
			chartbeatDebugger.Printf("...and now I'm awake!")
			session.Refresh()
		} else {
			break
		}
	}
}

func (t *TopPages) Run(session *mgo.Session) {
	chartbeatDebugger.Println("Fetching toppages")
	urls, err := f.FormatChartbeatUrls("live/toppages/v3", lib.Sites, globalConfig.ChartbeatApiKey)
	urls = f.AddUrlParams(urls, "loyalty=1")

	if err != nil {
		chartbeatDebugger.Printf("ERROR: %v", err)
		return
	}

	snapshot := f.FetchTopPages(urls)

	if session != nil {
		chartbeatDebugger.Println("Saving toppages snapshot")
		err := f.SaveTopPagesSnapshot(snapshot, session)
		if err != nil {
			chartbeatDebugger.Printf("ERROR: %v", err)
			return
		}

		f.CalculateTimeInterval(snapshot, session)

		// Update mapi to let it know that a new snapshot has been saved
		resp, err := http.Get("https://api.michigan.com/popular/")
		if err != nil {
			chartbeatDebugger.Printf("%v", err)
		} else {
			now := time.Now()
			chartbeatDebugger.Printf("Updated toppages snapshot at Mapi at %v", now)
		}

		defer resp.Body.Close()
	} else {
		chartbeatDebugger.Printf("Variable 'mongoUri' not specified, no data will be saved")
	}
}

func (q *QuickStats) Run(session *mgo.Session) {
	chartbeatDebugger.Printf("Quickstats")

	urls, err := f.FormatChartbeatUrls("live/quickstats/v4", lib.Sites, globalConfig.ChartbeatApiKey)
	urls = f.AddUrlParams(urls, "all_platforms=1&loyalty=1")
	if err != nil {
		chartbeatDebugger.Println("ERROR: %v", err)
		return
	}

	quickStats := f.FetchQuickStats(urls)

	if session != nil {
		chartbeatDebugger.Printf("Saving quickstats...")

		f.SaveQuickStats(quickStats, session)

		// Update mapi
		resp, err := http.Get("https://api.michigan.com/quickstats/")
		if err != nil {
			chartbeatDebugger.Printf("%v", err)
		} else {
			defer resp.Body.Close()
			chartbeatDebugger.Printf("Updated quickstats snapshot at Mapi at %v", time.Now())
		}
	} else {
		chartbeatDebugger.Printf("Variable 'mongoUri' not specified, no data will be saved")
		chartbeatDebugger.Printf("%v", quickStats)
	}
}

func (t *TopGeo) Run(session *mgo.Session) {
	chartbeatDebugger.Printf("Topgeo")

	urls, err := f.FormatChartbeatUrls("live/top_geo/v1", lib.Sites, globalConfig.ChartbeatApiKey)
	if err != nil {
		chartbeatDebugger.Println("ERROR: %v", err)
		return
	}

	topGeo := f.FetchTopGeo(urls)

	if session != nil {
		chartbeatDebugger.Printf("Saving topgeo...")

		f.SaveTopGeo(topGeo, session)

		// Update mapi
		resp, err := http.Get("https://api.michigan.com/topgeo/")
		if err != nil {
			chartbeatDebugger.Printf("%v", err)
		} else {
			defer resp.Body.Close()
			chartbeatDebugger.Printf("Updated topgeo snapshot at Mapi at %v", time.Now())
		}
	} else {
		chartbeatDebugger.Printf("Variable 'mongoUri' not specified, no data will be saved")
		chartbeatDebugger.Printf("%v", topGeo)
	}
}
