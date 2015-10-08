package commands

import (
	"net/http"
	"sync"
	"time"

	f "github.com/michigan-com/newsfetch/fetch/chartbeat"
	"github.com/michigan-com/newsfetch/lib"
	"github.com/spf13/cobra"
)

// Beats
type Beat interface {
	Run(string)
}

type TopPages struct{}
type QuickStats struct{}

var chartbeatDebugger = lib.NewCondLogger("chartbeat")

var cmdChartbeat = &cobra.Command{
	Use:   "chartbeat",
	Short: "Hit the Chartbeat API",
}

var cmdAllBeats = &cobra.Command{
	Use:   "all",
	Short: "Fetch all Chartbeat Beats (toppages, quickstats)",
	Run: func(cmd *cobra.Command, argv []string) {
		RunChartbeatCommands([]Beat{&TopPages{}, &QuickStats{}})
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

func RunChartbeatCommands(beats []Beat) {
	// Set up environment

	for {
		startTime := time.Now()

		// Run the actual meat of the program
		var beatWait sync.WaitGroup

		for _, beat := range beats {
			beatWait.Add(1)

			go func(beat Beat) {
				beat.Run(globalConfig.MongoUrl)
				beatWait.Done()
			}(beat)
		}

		beatWait.Wait()

		if timeit {
			getElapsedTime(&startTime)
		}

		if loop != -1 {
			chartbeatDebugger.Printf("Looping! Sleeping for %d seconds...", loop)
			time.Sleep(time.Duration(loop) * time.Second)
			chartbeatDebugger.Printf("...and now I'm awake!")
		} else {
			break
		}
	}
}

func (t *TopPages) Run(mongoUri string) {
	chartbeatDebugger.Println("Fetching toppages")
	urls, err := f.FormatChartbeatUrls("live/toppages/v3", lib.Sites, globalConfig.ChartbeatApiKey)

	if err != nil {
		chartbeatDebugger.Printf("ERROR: %v", err)
		return
	}

	snapshot := f.FetchTopPages(urls)

	if globalConfig.MongoUrl != "" {
		chartbeatDebugger.Println("Saving toppages snapshot")
		err := f.SaveTopPagesSnapshot(snapshot, globalConfig.MongoUrl)
		if err != nil {
			chartbeatDebugger.Printf("ERROR: %v", err)
			return
		}

		f.CalculateTimeInterval(snapshot, globalConfig.MongoUrl)

		// Update mapi to let it know that a new snapshot has been saved
		_, err = http.Get("https://api.michigan.com/popular/")
		if err != nil {
			chartbeatDebugger.Printf("%v", err)
		} else {
			now := time.Now()
			chartbeatDebugger.Printf("Updated toppages snapshot at Mapi at %v", now)
		}
	} else {
		chartbeatDebugger.Printf("Variable 'mongoUri' not specified, no data will be saved")
	}
}

func (q *QuickStats) Run(mongoUri string) {
	chartbeatDebugger.Printf("Quickstats")

	urls, err := f.FormatChartbeatUrls("live/quickstats/v4", lib.Sites, globalConfig.ChartbeatApiKey)
	if err != nil {
		chartbeatDebugger.Println("ERROR: %v", err)
		return
	}

	quickStats := f.FetchQuickStats(urls)

	if globalConfig.MongoUrl != "" {
		chartbeatDebugger.Printf("Saving quickstats...")

		f.SaveQuickStats(quickStats, globalConfig.MongoUrl)

		// Update mapi
		_, err = http.Get("https://api.michigan.com/quickstats/")
		if err != nil {
			chartbeatDebugger.Printf("%v", err)
		} else {
			chartbeatDebugger.Printf("Updated quickstats snapshot at Mapi at %v", time.Now())
		}
	} else {
		chartbeatDebugger.Printf("Variable 'mongoUri' not specified, no data will be saved")
		chartbeatDebugger.Printf("%v", quickStats)
	}
}
