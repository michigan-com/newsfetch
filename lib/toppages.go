package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const chartbeatApiUrlFormat = "http://api.chartbeat.com/%s/?apikey=%s&host=%s&limit=100"

/** Sorting stuff */
type ByVisits []*TopArticle

func (a ByVisits) Len() int           { return len(a) }
func (a ByVisits) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByVisits) Less(i, j int) bool { return a[i].Visits > a[j].Visits }

type TopPagesSnapshot struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Created_at time.Time     `bson:"created_at"`
	Articles   []*TopArticle `bson:"articles"`
}
type TopArticle struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	ArticleId int           `bson:"article_id"`
	Headline  string        `bson:"headline"`
	Url       string        `bson:"url"`
	Sections  []string      `bson:"sections"`
	Visits    int           `bson:"visits"`
}

type TopPages struct {
	Site  string
	Pages []*ArticleContent `json:"pages"`
}

type ArticleContent struct {
	Path     string        `json:"path"`
	Sections []string      `json:"sections"`
	Stats    *ArticleStats `json: "stats"`
	Title    string        `json:"title"`
}

type ArticleStats struct {
	Visits int `json:"visits"`
}

/*
	Fetch the top pages data for each url in the urls parameter. Url expected
	to be http://api.chartbeat.com/live/toppages/v3
*/
func FetchTopPages(urls []string) []*TopArticle {
	Debugger.Println("Fetching chartbeat top pages")
	topArticles := make([]*TopArticle, 0, 100*len(urls))
	articleQueue := make(chan *TopArticle, 100*len(urls))

	var wg sync.WaitGroup

	for i := 0; i < len(urls); i++ {
		wg.Add(1)

		go func(url string) {
			pages, err := GetTopPages(url)

			if err != nil {
				Debugger.Println("%v", err)
				wg.Done()
				return
			}

			for i := 0; i < len(pages.Pages); i++ {
				page := pages.Pages[i]
				articleUrl := page.Path
				articleId := GetArticleId(articleUrl)
				article := &TopArticle{}

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

				articleQueue <- article
			}

			wg.Done()
		}(urls[i])
	}

	wg.Wait()
	Debugger.Println("Done")
	close(articleQueue)

	for article := range articleQueue {
		topArticles = append(topArticles, article)
	}

	Debugger.Printf("Num article: %d", len(topArticles))

	Debugger.Println("Done fetching and parsing URLs...")

	return SortTopArticles(topArticles)
}

func CalculateTimeInterval(articles []*TopArticle, mongoUri string) {
	if mongoUri == "" {
		Debugger.Printf("No mongoUri, cannot calculate time intervals")
		return
	}

	Debugger.Printf("Updating numbers for articles for this given time interval")
	articleIds := make([]int, len(articles), len(articles))
	articleVisits := map[int]int{}
	savedArticles := make([]*Article, len(articles), len(articles))

	for _, article := range articles {
		articleIds = append(articleIds, article.ArticleId)
		articleVisits[article.ArticleId] = article.Visits
	}

	session := DBConnect(mongoUri)
	defer DBClose(session)
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
		Debugger.Printf("ERROR: %v", err)
		return
	}

	calculateTimeInterval(savedArticles, articleVisits)

	saveTimeInterval(savedArticles, session)
}

// In-memory adjusting of time intervals. Easier for testing since it doesnt
// hit mongo
func calculateTimeInterval(savedArticles []*Article, articleVisits map[int]int) {
	now := time.Now()

	for _, article := range savedArticles {
		visits, ok := articleVisits[article.ArticleId]
		if !ok {
			continue
		}
		CheckHourlyMax(article, now, visits)
	}
}

func saveTimeInterval(articles []*Article, session *mgo.Session) {
	articleCol := session.DB("").C("Article")
	for _, article := range articles {
		err := articleCol.Update(bson.M{"_id": article.Id}, bson.M{
			"$set": bson.M{
				"visits": article.Visits,
			},
		})
		if err != nil {
			Debugger.Printf("ERROR: %v", err)
		}
	}
}

/*
	Format chartbeat URLs based on a chartbeat API endpoint

	Format: http://api.chartbeat.com/<endPoint>/?apikey=<key>&host=<site[i]>

	Example endPoint (NOTE no starting or ending slashes): live/toppages/v3
*/
func FormatChartbeatUrls(endPoint string, sites []string, apiKey string) ([]string, error) {
	urls := make([]string, 0, len(sites))

	if apiKey == "" {
		return urls, errors.New(fmt.Sprintf("No API key specified. Use the -k flag to specify (Run ./newsfetch chartbeat --help for more info)"))
	}

	for i := 0; i < len(sites); i++ {
		site := sites[i]

		url := fmt.Sprintf(chartbeatApiUrlFormat, endPoint, apiKey, site)

		urls = append(urls, url)
	}

	return urls, nil
}

/*
	Given a URL for the api.chartbeat.com/live/toppages/v3 API, get the data and
	read the response
*/
func GetTopPages(url string) (*TopPages, error) {
	Debugger.Println("Fetching %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	Debugger.Println("Successfully fetched %s", url)

	topPages := &TopPages{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&topPages)

	return topPages, err
}

/*
	Save the toppages snapshot

	mongoUri - connection string to Mongodb
	toppages - Sorted array of top articles
*/
func SaveTopPagesSnapshot(toppages []*TopArticle, mongoUri string) error {
	Debugger.Println("Saving snapshot ...")

	// Save the current snapshot
	session := DBConnect(mongoUri)
	defer DBClose(session)
	snapshotCollection := session.DB("").C("Toppages")
	snapshot := TopPagesSnapshot{}
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
		Debugger.Println("Error when removing older snapshots: %v", err)
	}

	return nil
}

func SortTopArticles(articles []*TopArticle) []*TopArticle {
	Debugger.Println("Sorting articles ...")
	sort.Sort(ByVisits(articles))
	return articles
}
