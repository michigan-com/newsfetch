package main

import (
	"fmt"
	"github.com/michigan-com/newsfetch/lib"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
)

func printArticleBrief(articles []*Article) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "Source\tSection\tHeadline\tURL\tTimestamp")
	for _, article := range articles {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", article.Source, article.Section, article.Headline, article.Url, article.Timestamp)
	}
	w.Flush()
}

func main() {
	var (
		mongoUri   string
		articleUrl string
		siteStr    string
		sectionStr string
		output     bool
		body       bool
		verbose    bool
	)

	logging.SetLevel(logging.CRITICAL, "newsfetch")

	var cmdNewsfetch = &cobra.Command{
		Use: "newsfetch",
	}

	cmdNewsfetch.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	var cmdBody = &cobra.Command{
		Use:   "body",
		Short: "Get article body content from Gannett URL",
		Run: func(cmd *cobra.Command, args []string) {
			if verbose {
				logging.SetLevel(logging.DEBUG, "newsfetch")
			}

			if len(args) > 0 && args[0] != "" {
				articleUrl = args[0]
			}

			body, err := lib.ExtractBodyFromURL(articleUrl)
			if err != nil {
				panic(err)
			}

			fmt.Println(body)
		},
	}

	url := "http://detroitnews.com/story/news/local/detroit-city/2015/08/04/female-body-found-possible-hit-run-detroit/31094589/"
	cmdBody.Flags().StringVarP(&articleUrl, "url", "u", url, "URL of Gannett article")

	var cmdArticles = &cobra.Command{
		Use:   "articles",
		Short: "Fetches and parses Gannett news articles",
		Run: func(cmd *cobra.Command, args []string) {
			if verbose {
				logging.SetLevel(logging.DEBUG, "newsfetch")
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

			articles := FetchAndParseArticles(sites, sections, body)

			if output {
				printArticleBrief(articles)
			}

			if mongoUri != "" {
				fmt.Println(mongoUri)
				err := SaveArticles(mongoUri, articles)
				if err != nil {
					panic(err)
				}
			}
		},
	}

	cmdArticles.Flags().StringVarP(&siteStr, "sites", "i", "all", "Comma separated list of Gannett sites to fetch articles from")
	cmdArticles.Flags().StringVarP(&sectionStr, "sections", "e", "all", "Comma separated list of article sections to fetch from")
	cmdArticles.Flags().StringVarP(&mongoUri, "save", "s", "", "Saves articles to mongodb server specified in this option, e.g. mongodb://localhost:27017/")
	cmdArticles.Flags().BoolVarP(&output, "output", "o", true, "Outputs summary article inforation")
	cmdArticles.Flags().BoolVarP(&body, "body", "b", false, "Fetches the article body content")

	cmdNewsfetch.AddCommand(cmdBody, cmdArticles)
	cmdNewsfetch.Execute()
}
