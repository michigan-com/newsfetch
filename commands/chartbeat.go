package commands

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/michigan-com/newsfetch/lib"
	"github.com/spf13/cobra"
)

// Beats
type Beat interface {
	Run(string)
}

type TopPages struct{}
type QuickStats struct{}

var debugger = lib.NewCondLogger("chartbeat")

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
	if mongoUri == "" {
		mongoUri = os.Getenv("MONGO_URI")
	}
	if apiKey == "" {
		apiKey = os.Getenv("CHARTBEAT_API_KEY")
	}

	for {
		startTime := time.Now()

		// Run the actual meat of the program
		var beatWait sync.WaitGroup

		for _, beat := range beats {
			beatWait.Add(1)

			go func(beat Beat) {
				beat.Run(mongoUri)
				beatWait.Done()
			}(beat)
		}

		beatWait.Wait()

		if timeit {
			getElapsedTime(&startTime)
		}

		if loop != -1 {
			debugger.Printf("Looping! Sleeping for %d seconds...", loop)
			time.Sleep(time.Duration(loop) * time.Second)
			debugger.Printf("...and now I'm awake!")
		} else {
			break
		}
	}
}

func (t *TopPages) Run(mongoUri string) {
	debugger.Println("Fetching toppages")
	urls, err := lib.FormatChartbeatUrls("live/toppages/v3", lib.Sites, apiKey)
	if err != nil {
		debugger.Printf("ERROR: %v", err)
		return
	}

	snapshot := lib.FetchTopPages(urls)

	if mongoUri != "" {
		debugger.Println("Saving toppages snapshot")
		err := lib.SaveTopPagesSnapshot(snapshot, mongoUri)
		if err != nil {
			debugger.Printf("ERROR: %v", err)
			return
		}

		lib.CalculateTimeInterval(snapshot, mongoUri)

		// Update mapi to let it know that a new snapshot has been saved
		_, err = http.Get("https://api.michigan.com/popular/")
		if err != nil {
			debugger.Printf("%v", err)
		} else {
			now := time.Now()
			debugger.Printf("Updated snapshot at Mapi at %v", now)
		}
	} else {
		debugger.Printf("Variable 'mongoUri' not specified, no data will be saved")
	}
}

func (q *QuickStats) Run(mongoUri string) {
	debugger.Printf("Quickstats")
}
