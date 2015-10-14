package commands

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	a "github.com/michigan-com/newsfetch/fetch/article"
	"github.com/michigan-com/newsfetch/lib"
	"github.com/spf13/cobra"
)

var artDebugger = lib.NewCondLogger("command-article")

func processSummaries(ch chan error) {
	url := "http://brevity.detroitnow.io/newsfetch-summarize/"
	artDebugger.Println("Fetching: ", url)

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		ch <- err
	}

	ch <- nil
}

func processArticle(articleUrl string) {
	processor := a.ParseArticleAtURL(articleUrl, body /* global flag */)
	if processor.Err != nil {
		artDebugger.Println("Failed to process article: ", processor.Err)
		return
	}

	if globalConfig.MongoUrl != "" {
		artDebugger.Println("Attempting to save article ...")
		session := lib.DBConnect(globalConfig.MongoUrl)
		defer lib.DBClose(session)
		processor.Article.Save(session)
	}

	artDebugger.Println(processor.Article)
}

func formatFeedUrls(sites []string, sections []string) []string {
	urls := make([]string, 0, len(sites)*len(sections))

	for i := 0; i < len(sites); i++ {
		site := sites[i]
		for j := 0; j < len(sections); j++ {
			section := sections[j]

			if strings.Contains(site, "detroitnews") && section == "life" {
				section += "-home"
			}
			url := fmt.Sprintf("http://%s/feeds/live/%s/json", site, section)
			urls = append(urls, url)
		}
	}
	return urls
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

		/*if output {
			w.Init(os.Stdout, 0, 8, 0, '\t', 0)
			fmt.Fprintln(w, "Source\tSection\tHeadline\tURL\tTimestamp")
		}*/

		feedUrls := formatFeedUrls(sites, sections)

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
					processArticle(url)
				}(fmt.Sprintf("http://%s.com%s", host, aurl))
			}
		}
		close(ach)
		wg.Wait()

		lib.Logger.Println("Sending request to brevity to process summaries")
		go processSummaries(nil)

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

		articles, err := LoadRemoteArticles(url)
		if err != nil {
			panic(err)
		}

		artDebugger.Printf("Saving %d articles...\n", len(articles))
		session := lib.DBConnect(globalConfig.MongoUrl)
		defer lib.DBClose(session)

		for _, art := range articles {
			err = art.Save(session)
			if err != nil {
				panic(err)
			}
		}

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}
