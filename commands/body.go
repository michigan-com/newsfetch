package commands

import (
	"strings"
	"time"

	"github.com/michigan-com/newsfetch/extraction"
	"github.com/michigan-com/newsfetch/lib"
	"github.com/spf13/cobra"
)

var cmdBody = &cobra.Command{
	Use:   "body",
	Short: "Get article body content from Gannett URL",
	Run: func(cmd *cobra.Command, args []string) {
		startTime = time.Now()

		if len(args) > 0 && args[0] != "" {
			articleUrl = args[0]
		}

		extracted := extraction.ExtractDataFromHTMLAtURL(articleUrl, includeTitle)

		bodyFmt := strings.Join(strings.Split(extracted.Text, "\n"), "\n\n")
		lib.Logger.Println(bodyFmt)

		getElapsedTime(&startTime)
	},
}
