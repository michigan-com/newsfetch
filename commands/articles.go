package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"

	a "github.com/michigan-com/newsfetch/fetch/article"
	"github.com/michigan-com/newsfetch/lib"
	"github.com/spf13/cobra"
)

var artDebugger = lib.NewCondLogger("command-article")

type SummaryResponse struct {
	Skipped    int `json:"skipped"`
	Summarized int `json:"summarized"`
}

func processSummaries() (*SummaryResponse, error) {
	url := "http://brevity.detroitnow.io/newsfetch-summarize/"
	artDebugger.Println("Fetching: ", url)

	resp, err := http.Get(url)
	defer resp.Body.Close()

	var jso []byte
	summResp := SummaryResponse{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&summResp)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(jso, summResp)

	return &summResp, nil
}

func processArticle(articleUrl string, session *mgo.Session) bool {
	processor := a.ParseArticleAtURL(articleUrl, body /* global flag */)
	if processor.Err != nil {
		artDebugger.Println("Failed to process article: ", processor.Err)
		return false
	}

	var isNew bool
	var err error
	if globalConfig.MongoUrl != "" {
		if session == nil {
			session = lib.DBConnect(globalConfig.MongoUrl)
			defer lib.DBClose(session)
		}

		artDebugger.Println("Attempting to save article: ", processor.Article)
		isNew, err = processor.Article.Save(session)
		if err != nil {
			lib.Logger.Println(err)
		}
	}

	artDebugger.Println(processor.Article)
	return isNew
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

		// create one session for all saves bruh
		var session *mgo.Session
		if globalConfig.MongoUrl != "" {
			session = lib.DBConnect(globalConfig.MongoUrl)
			defer lib.DBClose(session)
		}

		newArticles := 0
		updatedArticles := 0

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
					isNew := processArticle(url, session)
					if isNew {
						newArticles++
					} else {
						updatedArticles++
					}
				}(fmt.Sprintf("http://%s.com%s", host, aurl))
			}
		}
		close(ach)
		wg.Wait()

		lib.Logger.Println("Sending request to brevity to process summaries")
		sumRes, err := processSummaries()
		if err != nil {
			lib.Logger.Println("Summarizer failed: ", err)
		}

		lib.Logger.Println("New articles: ", newArticles)
		lib.Logger.Println("Updated articles: ", updatedArticles)
		lib.Logger.Println("Skipped article summaries: ", sumRes.Skipped)
		lib.Logger.Println("Summarized articles: ", sumRes.Summarized)

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
			_, err := art.Save(session)
			if err != nil {
				panic(err)
			}
		}

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}
