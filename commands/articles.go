package commands

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/michigan-com/newsfetch/lib"
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

		var err error

		foodArticles := lib.FilterArticlesForRecipeExtraction(articles)
		if len(foodArticles) > 0 {
			err = lib.DownloadAndSaveRecipesForArticles(globalConfig.MongoUrl, foodArticles)
			if err != nil {
				panic(err)
			}
		}

		if globalConfig.MongoUrl != "" {
			/*err := lib.RemoveArticles(globalConfig.MongoUrl)
			if err != nil {
				panic(err)
			}*/

			err = lib.SaveArticles(globalConfig.MongoUrl, articles)
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

		err := lib.RemoveArticles(globalConfig.MongoUrl)
		if err != nil {
			panic(err)
		}

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}

var cmdCopyArticles = &cobra.Command{
	Use:   "copy-from",
	Short: "Copies articles from a mapi JSON URL",
	Run: func(cmd *cobra.Command, args []string) {
		if timeit {
			startTime = time.Now()
		}

		if len(args) != 1 {
			panic("Required argument: URL")
		}
		url := args[0]

		articles, err := lib.LoadRemoteArticles(url)
		if err != nil {
			panic(err)
		}

		printArticleBrief(articles)

		fmt.Printf("Saving %d articles...\n", len(articles))
		err = lib.SaveArticles(globalConfig.MongoUrl, articles)
		if err != nil {
			panic(err)
		}

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}
