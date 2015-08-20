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

		if verbose {
			Verbose("")
		}

		if len(args) > 0 && args[0] != "" {
			articleUrl = args[0]
		}

		var body string
		ch := make(chan string)
		go lib.ExtractBodyFromURL(ch, articleUrl, includeTitle)
		body = <-ch

		if output {
			bodyFmt := strings.Join(strings.Split(body, "\n"), "\n\n")
			fmt.Println(bodyFmt)
		}

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}
