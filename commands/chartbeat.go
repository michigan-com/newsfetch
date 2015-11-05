package commands

import (
	"sync"
	"time"

	c "github.com/michigan-com/newsfetch/chartbeat"
	"github.com/michigan-com/newsfetch/lib"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2"
)

var chartbeatDebugger = lib.NewCondLogger("newsfetch:commands:chartbeat")

var cmdChartbeat = &cobra.Command{
	Use:   "chartbeat",
	Short: "Hit the Chartbeat API",
}

var cmdAllBeats = &cobra.Command{
	Use:   "all",
	Short: "Fetch all Chartbeat Beats (toppages, quickstats)",
	Run: func(cmd *cobra.Command, argv []string) {
		RunChartbeatCommands([]c.Beat{
			c.TopPages,
			c.QuickStats,
			c.TopGeo,
			c.Referrers,
			c.Recent,
			c.Historical,
		})
	},
}

var cmdTopPages = &cobra.Command{
	Use:   "toppages",
	Short: "Fetch toppages snapshot for Chartbeat",
	Run: func(cmd *cobra.Command, argv []string) {
		RunChartbeatCommands([]c.Beat{&c.TopPages})
	},
}

var cmdQuickStats = &cobra.Command{
	Use:   "quickstats",
	Short: "Fetch quickstats snapshot for Chartbeat",
	Run: func(cmd *cobra.Command, argv []string) {
		RunChartbeatCommands([]c.Beat{&c.QuickStats})
	},
}

var cmdTopGeo = &cobra.Command{
	Use:   "topgeo",
	Short: "Fetch topgeo snapshot for Chartbeat",
	Run: func(cmd *cobra.Command, argv []string) {
		RunChartbeatCommands([]c.Beat{&c.TopGeo})
	},
}

var cmdReferrers = &cobra.Command{
	Use:   "referrers",
	Short: "Fetch referrers snapshot for Chartbeat",
	Run: func(cmd *cobra.Command, arg []string) {
		RunChartbeatCommands([]c.Beat{&c.Referrers})
	},
}

var cmdRecent = &cobra.Command{
	Use:   "recent",
	Short: "Fetch recent snapshot for Chartbeat",
	Run: func(cmd *cobra.Command, arg []string) {
		RunChartbeatCommands([]c.Beat{&c.Recent})
	},
}

var cmdHistorical = &cobra.Command{
	Use:   "historical-traffic",
	Short: "Fetch recent snapshot for Chartbeat",
	Run: func(cmd *cobra.Command, arg []string) {
		RunChartbeatCommands([]c.Beat{&c.Historical})
	},
}

func RunChartbeatCommands(beats []c.Beat) {
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

			go func(beat c.Beat) {
				var _copy *mgo.Session
				if session != nil {
					_copy = session.Copy()
					defer _copy.Close()
				}
				beat.Run(_copy, globalConfig.ChartbeatApiKey)
				beatWait.Done()
			}(beat)
		}

		beatWait.Wait()

		getElapsedTime(&startTime)

		if loop != -1 {
			chartbeatDebugger.Printf("Looping! Sleeping for %d seconds...", loop)
			time.Sleep(time.Duration(loop) * time.Second)
			chartbeatDebugger.Printf("...and now I'm awake!")
			if session != nil {
				session.Refresh()
			}
		} else {
			break
		}
	}
}
