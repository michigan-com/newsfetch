package commands

import (
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
		urls := lib.FormatChartbeatUrls("live/toppages/v3", lib.Sites)
		snapshot := lib.FetchTopPages(urls)

		if mongoUri != "" {
			logger.Info("Saving toppages snapshot")
			err := lib.SaveTopPagesSnapshot(mongoUri, snapshot)

			if err != nil {
				panic(err)
			}

		} else {
			logger.Warning("Variable 'mongoUri' not specified, no data will be saved")
		}
	},
}
