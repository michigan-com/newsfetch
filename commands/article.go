package commands

import (
	"fmt"
	"time"

	"github.com/michigan-com/newsfetch/lib"
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

		article := &lib.Article{}
		lib.TheOneRing(article,
			lib.NewArticleIn(articleUrl),
			lib.NewBodyParser(),
		)

		if output {
			fmt.Println(article)
		}

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}
