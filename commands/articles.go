package commands

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/michigan-com/newsfetch/lib"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
)

func printArticleBrief(articles []*lib.Article) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "Source\tSection\tHeadline\tURL\tTimestamp")
	for _, article := range articles {
		fmt.Fprintf(
			w, "%s\t%s\t%s\t%s\t%s\n", article.Source, article.Section,
			article.Headline, article.Url, article.Timestamp,
		)
	}
	w.Flush()
}

var cmdArticles = &cobra.Command{
	Use:   "articles",
	Short: "Command for Gannett news articles",
}

var cmdGetArticles = &cobra.Command{
	Use:   "get",
	Short: "Fetches, parses, and saves news articles",
	Run: func(cmd *cobra.Command, args []string) {
		if timeit {
			startTime = time.Now()
		}

		if verbose {
			Verbose()
		}

		var sites []string
		var sections []string

		if siteStr == "all" {
			sites = lib.Sites
		} else {
			sites = strings.Split(siteStr, ",")
		}

		if sectionStr == "all" {
			sections = lib.Sections
		} else {
			sections = strings.Split(sectionStr, ",")
		}

		urls := lib.FormatFeedUrls(sites, sections)
		articles := lib.FetchAndParseArticles(urls, body)

		if output {
			printArticleBrief(articles)
		}

		if mongoUri != "" {
			/*err := lib.RemoveArticles(mongoUri)
			if err != nil {
				panic(err)
			}*/

			err := lib.SaveArticles(mongoUri, articles)
			if err != nil {
				panic(err)
			}
		}

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}

var cmdRemoveArticles = &cobra.Command{
	Use:   "rm",
	Short: "Removes news articles from mongodb",
	Run: func(cmd *cobra.Command, args []string) {
		if !noprompt {
			resp := "n"
			fmt.Printf("Are you sure you want to remove all articles from Snapshot collection? [y/N]: ")
			fmt.Scanf("%s", &resp)

			if strings.ToLower(resp) != "y" {
				return
			}
		}

		if timeit {
			startTime = time.Now()
		}

		if verbose {
			logging.SetLevel(logging.DEBUG, "newsfetch")
		}

		err := lib.RemoveArticles(mongoUri)
		if err != nil {
			panic(err)
		}

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}
