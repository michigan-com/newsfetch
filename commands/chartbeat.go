package commands

import (
	"net/http"
	"os"
	"time"

	"github.com/michigan-com/newsfetch/lib"
	"github.com/spf13/cobra"
)

var cmdChartbeat = &cobra.Command{
	Use:   "chartbeat",
	Short: "Hit the Chartbeat API",
}

var cmdTopPages = &cobra.Command{
	Use:   "toppages",
	Short: "Fetch toppages snapshot for Chartbeat",
	Run:   ChartbeatToppagesCommand,
}

func ChartbeatToppagesCommand(cmd *cobra.Command, args []string) {
	// Set up environment
	if mongoUri == "" {
		mongoUri = os.Getenv("MONGO_URI")
	}
	if apiKey == "" {
		apiKey = os.Getenv("CHARTBEAT_API_KEY")
	}

	startTime := time.Now()

	if verbose {
		Verbose("")
	}

	// Run the actual meat of the program
	ChartbeatToppages()

	if timeit {
		getElapsedTime(&startTime)
	}

	if loop != -1 {
		logger.Info("Looping! Sleeping for %d seconds...", loop)
		time.Sleep(time.Duration(loop) * time.Second)
		logger.Info("...and now I'm awake!")
		ChartbeatToppagesCommand(cmd, args)
	}
}

func ChartbeatToppages() {
	logger.Info("Fetching toppages")
	urls, err := lib.FormatChartbeatUrls("live/toppages/v3", lib.Sites, apiKey)
	if err != nil {
		logger.Error("%v", err)
		return
	}

	snapshot := lib.FetchTopPages(urls)

	if mongoUri != "" {
		logger.Info("Saving toppages snapshot")
		err := lib.SaveTopPagesSnapshot(mongoUri, snapshot)

		if err != nil {
			logger.Error("%v", err)
			return
		}

		// Update mapi to let it know that a new snapshot has been saved
		_, err = http.Get("https://api.michigan.com/popular/")
		if err != nil {
			logger.Error("%v", err)
		} else {
			now := time.Now()
			logger.Info("Updated snapshot at Mapi at %v", now)
		}
	} else {
		logger.Warning("Variable 'mongoUri' not specified, no data will be saved")
	}
}
