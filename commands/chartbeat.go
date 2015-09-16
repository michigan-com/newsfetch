package commands

import (
	"net/http"
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
	Run: func(cmd *cobra.Command, args []string) {
		if timeit {
			startTime = time.Now()
		}

		if verbose {
			Verbose("")
		}

		logger.Info("Fetching toppages")
		urls, err := lib.FormatChartbeatUrls("live/toppages/v3", lib.Sites, apiKey)
		if err != nil {
			panic(err)
		}

		snapshot := lib.FetchTopPages(urls)

		if mongoUri != "" {
			logger.Info("Saving toppages snapshot")
			err := lib.SaveTopPagesSnapshot(mongoUri, snapshot)

			if err != nil {
				panic(err)
			}

			// Update mapi to let it know that a new snapshot has been saved
			_, err = http.Get("https://api.michigan.com/popular/")
			if err != nil {
				logger.Error("%v", err)
			} else {
				logger.Info("Updated snapshot at Mapi at %v", time.Now())
			}
		} else {
			logger.Warning("Variable 'mongoUri' not specified, no data will be saved")
		}
	},
}
