package fetch

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	e "github.com/michigan-com/newsfetch/extraction"
	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model"
	mc "github.com/michigan-com/newsfetch/model/chartbeat"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TopPages struct{}

/** Sorting stuff */
type ByVisits []*mc.TopArticle

func (a ByVisits) Len() int           { return len(a) }
func (a ByVisits) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByVisits) Less(i, j int) bool { return a[i].Visits > a[j].Visits }

/*
	Fetch the top pages data for each url in the urls parameter. Url expected
	to be http://api.chartbeat.com/live/toppages/v3
*/
func (t TopPages) Fetch(urls []string) mc.Snapshot {
	chartbeatDebugger.Println("Fetching chartbeat top pages")
	topArticles := make([]*mc.TopArticle, 0, 100*len(urls))
	articleQueue := make(chan *mc.TopArticle, 100*len(urls))

	var wg sync.WaitGroup

	for i := 0; i < len(urls); i++ {
		wg.Add(1)

		go func(url string) {
			pages, err := GetTopPages(url)
			host, _ := GetHostFromParams(url)

			if err != nil {
				chartbeatError.Println("Failed to json parse url %s: %v", url, err)
				wg.Done()
				return
			}

			for i := 0; i < len(pages.Pages); i++ {
				page := pages.Pages[i]
				articleUrl := page.Path
				articleId := lib.GetArticleId(articleUrl)
				article := &mc.TopArticle{}

				// this means we can't find an article ID. It's probably a section front,
				// so ignore
				if articleId < 0 {
					continue
				}

				article.ArticleId = articleId
				article.Headline = page.Title
				article.Url = page.Path
				article.Sections = page.Sections
				article.Visits = page.Stats.Visits
				article.Loyalty = page.Stats.Loyalty
				article.Authors = e.ParseAuthors(page.Authors)
				article.Source = strings.Replace(host, ".com", "", -1)

				articleQueue <- article
			}

			wg.Done()
		}(urls[i])
	}

	wg.Wait()
	chartbeatDebugger.Println("Done")
	close(articleQueue)

	for article := range articleQueue {
		topArticles = append(topArticles, article)
	}

	chartbeatDebugger.Printf("Num article: %d", len(topArticles))

	chartbeatDebugger.Println("Done fetching and parsing URLs...")

	snapshot := mc.TopPagesSnapshot{}
	snapshot.Articles = SortTopArticles(topArticles)
	snapshot.Created_at = time.Now()
	return snapshot
}

func CalculateTimeInterval(articles []*mc.TopArticle, session *mgo.Session) {
	if session == nil {
		chartbeatDebugger.Printf("No session, cannot calculate time intervals")
		return
	}

	chartbeatDebugger.Printf("Updating numbers for articles for this given time interval")
	articleIds := make([]int, len(articles), len(articles))
	articleVisits := map[int]int{}
	savedArticles := make([]*m.Article, len(articles), len(articles))

	for _, article := range articles {
		articleIds = append(articleIds, article.ArticleId)
		articleVisits[article.ArticleId] = article.Visits
	}

	db := session.DB("")
	articleCol := db.C("Article")

	err := articleCol.
		Find(bson.M{
		"article_id": bson.M{
			"$in": articleIds,
		},
	}).
		Select(bson.M{
		"article_id": 1,
		"visits":     1,
	}).
		All(&savedArticles)

	if err != nil {
		chartbeatDebugger.Printf("ERROR: %v", err)
		return
	}

	calculateTimeInterval(savedArticles, articleVisits)

	saveTimeInterval(savedArticles, session)
}

// In-memory adjusting of time intervals. Easier for testing since it doesnt
// hit mongo
func calculateTimeInterval(savedArticles []*m.Article, articleVisits map[int]int) {
	now := time.Now()

	for _, article := range savedArticles {
		visits, ok := articleVisits[article.ArticleId]
		if !ok {
			continue
		}
		e.CheckHourlyMax(article, now, visits)
	}
}

func saveTimeInterval(articles []*m.Article, session *mgo.Session) {
	articleCol := session.DB("").C("Article")
	for _, article := range articles {
		err := articleCol.Update(bson.M{"_id": article.Id}, bson.M{
			"$set": bson.M{
				"visits": article.Visits,
			},
		})
		if err != nil {
			chartbeatDebugger.Printf("ERROR: %v", err)
		}
	}
}

/*
	Given a URL for the api.chartbeat.com/live/toppages/v3 API, get the data and
	read the response
*/
func GetTopPages(url string) (*mc.TopPages, error) {
	chartbeatDebugger.Println("Fetching %s", url)

	resp, err := http.Get(url)
	if err != nil {
		chartbeatError.Printf("Failed to get url %s: %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	chartbeatDebugger.Println("Successfully fetched %s", url)

	topPages := &mc.TopPages{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&topPages)

	return topPages, err
}

func SortTopArticles(articles []*mc.TopArticle) []*mc.TopArticle {
	chartbeatDebugger.Println("Sorting articles ...")
	sort.Sort(ByVisits(articles))
	return articles
}
