package commands

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/michigan-com/newsfetch/lib"
	"github.com/spf13/cobra"
)

var w = new(tabwriter.Writer)

func printArticleBrief(w *tabwriter.Writer, article *lib.Article) {
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

		feedUrls := lib.FormatFeedUrls(sites, sections)

		var wg sync.WaitGroup
		ach := make(chan *lib.ArticleUrlsChan)
		for _, url := range feedUrls {
			go lib.GetArticleUrlsFromFeed(url, ach)
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

		lib.Debugger.Println("Sending request to brevity to process summaries")
		go lib.ProcessSummaries(nil)
		/*var err error
		foodArticles := lib.FilterArticlesForRecipeExtraction(articles)
		if len(foodArticles) > 0 {
			err = lib.DownloadAndSaveRecipesForArticles(globalConfig.MongoUrl, foodArticles)
			if err != nil {
				panic(err)
			}
		}*/

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}

func ProcessArticle(articleUrl string) {
	article := &lib.Article{}

	articleIn := lib.NewArticleIn(articleUrl)
	err := articleIn.GetData()

	if err != nil {
		lib.Debugger.Println(err)
		return
	}

	if !articleIn.IsValid() {
		lib.Debugger.Println("Article is not valid: ", article)
		return
	}

	err = articleIn.Process(article)
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
	}

	fmt.Println(article)
	/*if output {
		printArticleBrief(w, article)
	}*/
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

		//printArticleBrief(articles)

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
