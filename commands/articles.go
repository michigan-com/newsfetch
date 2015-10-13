package commands

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	a "github.com/michigan-com/newsfetch/fetch/article"
	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model"
	"github.com/spf13/cobra"
)

var artDebugger = lib.NewCondLogger("command-article")

var w = new(tabwriter.Writer)

func printArticleBrief(w *tabwriter.Writer, article *m.Article) {
	fmt.Fprintf(
		w, "%s\t%s\t%s\t%s\t%s\n", article.Source, article.Section,
		article.Headline, article.Url, article.Timestamp,
	)
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

		if output {
			w.Init(os.Stdout, 0, 8, 0, '\t', 0)
			fmt.Fprintln(w, "Source\tSection\tHeadline\tURL\tTimestamp")
		}

		feedUrls := a.FormatFeedUrls(sites, sections)

		var wg sync.WaitGroup
		ach := make(chan *a.ArticleUrlsChan)
		for _, url := range feedUrls {
			go a.GetArticleUrlsFromFeed(url, ach)
			aurls := <-ach
			for _, aurl := range aurls.Urls {
				host, _ := lib.GetHost(url)
				wg.Add(1)
				go func(url string) {
					defer wg.Done()
					ProcessArticle(url)
				}(fmt.Sprintf("http://%s.com%s", host, aurl))
			}
		}
		close(ach)
		wg.Wait()

		artDebugger.Println("Sending request to brevity to process summaries")
		go a.ProcessSummaries(nil)

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}

func ProcessArticle(articleUrl string) {
	article, _, _, err := a.ParseArticleAtURL(articleUrl, body /* global flag */)
	if err != nil {
		artDebugger.Println("Failed to process article: ", err)
		return
	}

	if globalConfig.MongoUrl != "" {
		artDebugger.Println("Attempting to save article ...")
		session := lib.DBConnect(globalConfig.MongoUrl)
		defer lib.DBClose(session)
		article.Save(session)
	}

	artDebugger.Println(article)
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

		articles, err := a.LoadRemoteArticles(url)
		if err != nil {
			panic(err)
		}

		artDebugger.Printf("Saving %d articles...\n", len(articles))
		err = a.SaveArticles(globalConfig.MongoUrl, articles)
		if err != nil {
			panic(err)
		}

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}
