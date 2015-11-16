package fetch

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model/chartbeat"
)

type Recent struct {}

func (r Recent) Fetch(urls []string, session *mgo.Session) m.Snapshot {
	var wait sync.WaitGroup
	queue := make(chan *m.RecentResp, len(urls))

	for _, url := range urls {
		wait.Add(1)

		go func(url string) {

			recent, err := GetRecents(url)
			if err != nil {
				chartbeatDebugger.Printf("Failed to get %s: %v", url, err)
			} else {
				parsed_articles := make([]m.Recent, 0, 100)
				for _, article := range recent.Recents {
					articleId := lib.GetArticleId(article.Url)

					if articleId > 0 {
						article.Host = strings.Replace(article.Host, ".com", "", -1)
						parsed_articles = append(parsed_articles, article)
					}
				}

				recent.Recents = parsed_articles
				queue <- recent
			}
			wait.Done()
		}(url)
	}

	wait.Wait()
	close(queue)

	recents := make([]*m.RecentResp, 0, len(urls))
	for recent := range queue {
		recents = append(recents, recent)
	}

	snapshot := m.RecentSnapshot{}
	snapshot.Created_at = time.Now()
	snapshot.Recents = recents
	return snapshot
}

func GetRecents(url string) (*m.RecentResp, error) {
	chartbeatDebugger.Printf("Getting %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	host, _ := GetHostFromParams(url)

	recentArray := make([]m.Recent, 0, 100)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&recentArray)
	if err != nil {
		return nil, err
	}

	recent := &m.RecentResp{}
	recent.Recents = recentArray
	recent.Source = strings.Replace(host, ".com", "", -1)

	return recent, nil
}
