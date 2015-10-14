package commands

import (
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

		artDebugger.Println(articleUrl)
		/*if output {
			w.Init(os.Stdout, 0, 8, 0, '\t', 0)
			fmt.Fprintln(w, "Source\tSection\tHeadline\tURL\tTimestamp")
		}*/

		processArticle(articleUrl, nil)

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}
