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
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/** Sorting stuff */
type ByVisits []*m.TopArticle

func (a ByVisits) Len() int           { return len(a) }
func (a ByVisits) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByVisits) Less(i, j int) bool { return a[i].Visits > a[j].Visits }

/*
	Fetch the top pages data for each url in the urls parameter. Url expected
	to be http://api.chartbeat.com/live/toppages/v3
*/
func FetchTopPages(urls []string) []*m.TopArticle {
	chartbeatDebugger.Println("Fetching chartbeat top pages")
	topArticles := make([]*m.TopArticle, 0, 100*len(urls))
	articleQueue := make(chan *m.TopArticle, 100*len(urls))

	var wg sync.WaitGroup

	for i := 0; i < len(urls); i++ {
		wg.Add(1)

		go func(url string) {
			pages, err := GetTopPages(url)
			host, _ := GetHostFromParams(url)

			if err != nil {
				chartbeatDebugger.Println("%v", err)
				wg.Done()
				return
			}

			for i := 0; i < len(pages.Pages); i++ {
				page := pages.Pages[i]
				articleUrl := page.Path
				articleId := lib.GetArticleId(articleUrl)
				article := &m.TopArticle{}

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

	return SortTopArticles(topArticles)
}

func CalculateTimeInterval(articles []*m.TopArticle, session *mgo.Session) {
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
func GetTopPages(url string) (*m.TopPages, error) {
	chartbeatDebugger.Println("Fetching %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	chartbeatDebugger.Println("Successfully fetched %s", url)

	topPages := &m.TopPages{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&topPages)

	return topPages, err
}

/*
	Save the toppages snapshot

	mongoUri - connection string to Mongodb
	toppages - Sorted array of top articles
*/
func SaveTopPagesSnapshot(toppages []*m.TopArticle, session *mgo.Session) error {
	chartbeatDebugger.Println("Saving snapshot ...")

	// Save the current snapshot
	snapshotCollection := session.DB("").C("Toppages")
	snapshot := m.TopPagesSnapshot{}
	snapshot.Articles = toppages
	snapshot.Created_at = time.Now()
	snapshotCollection.Insert(snapshot)

	// remove all the other snapshots
	snapshotCollection.Find(bson.M{}).
		Select(bson.M{"_id": 1}).
		Sort("-_id").
		One(&snapshot)

	_, err := snapshotCollection.RemoveAll(bson.M{
		"_id": bson.M{
			"$ne": snapshot.Id,
		},
	})

	if err != nil {
		chartbeatDebugger.Println("Error when removing older snapshots: %v", err)
	}

	return nil
}

func SortTopArticles(articles []*m.TopArticle) []*m.TopArticle {
	chartbeatDebugger.Println("Sorting articles ...")
	sort.Sort(ByVisits(articles))
	return articles
}
