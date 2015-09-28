package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/michigan-com/newsfetch/lib"
	"github.com/spf13/cobra"
)

var cmdBody = &cobra.Command{
	Use:   "body",
	Short: "Get article body content from Gannett URL",
	Run: func(cmd *cobra.Command, args []string) {
		if timeit {
			startTime = time.Now()
		}

		if len(args) > 0 && args[0] != "" {
			articleUrl = args[0]
		}

		var extracted *lib.ExtractedBody
		ch := make(chan *lib.ExtractedBody)
		go lib.ExtractBodyFromURL(ch, articleUrl, includeTitle)
		extracted = <-ch

		if output {
			bodyFmt := strings.Join(strings.Split(extracted.Text, "\n"), "\n\n")
			fmt.Println(bodyFmt)
		}

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}
