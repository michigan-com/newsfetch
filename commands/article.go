package commands

import (
	"fmt"
	"strings"
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

		articleIn := lib.NewArticleIn(articleUrl)
		err := articleIn.Process(article)
		if err != nil {
			lib.Debugger.Println("Article could not be processed: %s", articleIn)
		}

		if body {
			ch := make(chan *lib.ExtractedBody)
			go lib.ExtractBodyFromURL(ch, articleUrl, false)
			bodyExtract := <-ch

			if bodyExtract.Text != "" {
				lib.Debugger.Printf(
					"Extracted extracted contains %d characters, %d paragraphs.",
					len(strings.Split(bodyExtract.Text, "")),
					len(strings.Split(bodyExtract.Text, "\n\n")),
				)
				article.BodyText = bodyExtract.Text
			}
		}

		if globalConfig.MongoUrl != "" {
			lib.Debugger.Println("Attempting to save article ...")
			session := lib.DBConnect(globalConfig.MongoUrl)
			defer lib.DBClose(session)
			article.Save(session)

			if article.BodyText != "" {
				lib.Debugger.Println("Attempting to process summary ...")
				go lib.ProcessSummaries(nil)
			}
		}

		if output {
			fmt.Println(article)
		}

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}
