package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var cmdArticle = &cobra.Command{
	Use:   "article",
	Short: "Command to process a single article", Run: func(cmd *cobra.Command, args []string) {
		if timeit {
			startTime = time.Now()
		}

		if len(args) > 0 && args[0] != "" {
			articleUrl = args[0]
		}

		if output {
			w.Init(os.Stdout, 0, 8, 0, '\t', 0)
			fmt.Fprintln(w, "Source\tSection\tHeadline\tURL\tTimestamp")
		}

		ProcessArticle(articleUrl)

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}
