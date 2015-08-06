package main

import (
	"bufio"
	"fmt"
	"github.com/michigan-com/newsfetch/lib"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/urandom/text-summary/summarize"
	"io/ioutil"
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
		mongoUri     string
		articleUrl   string
		siteStr      string
		sectionStr   string
		title        string
		output       bool
		body         bool
		verbose      bool
		includeTitle bool
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

			body, err := lib.ExtractBodyFromURL(articleUrl, includeTitle)
			if err != nil {
				panic(err)
			}

			fmt.Println(body)
		},
	}

	url := "http://www.freep.com/story/news/local/michigan/2015/08/06/farid-fata-cancer-sentencing/31213475/"
	cmdBody.Flags().StringVarP(&articleUrl, "url", "u", url, "URL of Gannett article")
	cmdBody.Flags().BoolVarP(&includeTitle, "title", "t", false, "Place title of article on the first line of output")

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

	var cmdSummary = &cobra.Command{
		Use:   "summary",
		Short: "Attempts to generate a summary based on an article body",
		Run: func(cmd *cobra.Command, args []string) {
			var summary summarize.Summarize

			if title == "" {
				reader := bufio.NewReader(os.Stdin)
				contents, _ := ioutil.ReadAll(reader)
				lines := strings.Split(string(contents), "\n")

				summary = summarize.NewFromString(lines[0], strings.Join(lines[1:], "\n"))
			} else {
				summary = summarize.New(title, os.Stdin)
			}

			for _, point := range summary.KeyPoints() {
				fmt.Println(point)
			}
		},
	}

	cmdSummary.Flags().StringVarP(&title, "title", "t", "", "Title for article summarizer, if not supplied then the summarizer assumes first line is title")

	cmdNewsfetch.AddCommand(cmdBody, cmdArticles, cmdSummary)
	cmdNewsfetch.Execute()
}
