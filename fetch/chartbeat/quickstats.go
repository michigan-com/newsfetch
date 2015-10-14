package fetch

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	m "github.com/michigan-com/newsfetch/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type QuickStatsSort []*m.QuickStats

func (q QuickStatsSort) Len() int           { return len(q) }
func (q QuickStatsSort) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q QuickStatsSort) Less(i, j int) bool { return q[i].Visits > q[j].Visits }

func FetchQuickStats(urls []string) []*m.QuickStats {

	var urlWait sync.WaitGroup
	statQueue := make(chan *m.QuickStats, len(urls))

	for _, url := range urls {
		urlWait.Add(1)

		go func(url string) {
			stats, err := GetQuickStats(url)

			if err != nil {
				chartbeatDebugger.Println("ERROR: %v", err)
			} else {
				statQueue <- stats
			}
			urlWait.Done()
		}(url)
	}

	urlWait.Wait()
	close(statQueue)

	quickStats := make([]*m.QuickStats, 0, len(urls))
	for stats := range statQueue {
		quickStats = append(quickStats, stats)
	}

	return SortQuickStats(quickStats)
}

func GetQuickStats(url string) (*m.QuickStats, error) {
	chartbeatDebugger.Println("Fetching %s", url)

	// Parse out the host we're getting data on
	host, err := GetHostFromParams(url)
	if err != nil {
		chartbeatDebugger.Printf("ERROR: %v", err)
		chartbeatDebugger.Printf("Host will be \"\"")
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	chartbeatDebugger.Println("Successfully fetched %s", url)

	quickStatsResp := &m.QuickStatsResp{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&quickStatsResp)

	if err != nil {
		return nil, err
	}

	quickStats := quickStatsResp.Data.Stats
	quickStats.Source = strings.Replace(host, ".com", "", -1)

	return quickStats, err
}

func SaveQuickStats(quickStats []*m.QuickStats, session *mgo.Session) {
	quickStatsCol := session.DB("").C("Quickstats")

	// Insert this snapshot
	snapshot := m.QuickStatsSnapshot{}
	snapshot.Created_at = time.Now()
	snapshot.Stats = quickStats
	err := quickStatsCol.Insert(snapshot)

	if err != nil {
		chartbeatDebugger.Printf("ERROR: %v", err)
	}

	// Remove old snapshots
	quickStatsCol.Find(bson.M{}).
		Select(bson.M{"_id": 1}).
		Sort("-_id").
		One(&snapshot)

	_, err = quickStatsCol.RemoveAll(bson.M{
		"_id": bson.M{
			"$ne": snapshot.Id,
		},
	})

	if err != nil {
		chartbeatDebugger.Printf("Error while removing old quickstats snapshots %v", err)
	}
}

// Chartbeat queries have a GET parameter "host", which represents the host
// we're getting data on. Pull the host from the url and return it.
// Return host (e.g. freep.com)
// Return "" if we don't find one
func GetHostFromParams(inputUrl string) (string, error) {
	var host string
	var err error

	parsed, err := url.Parse(inputUrl)
	if err != nil {
		return host, err
	}

	hosts := parsed.Query()["host"]
	if len(hosts) > 0 {
		host = hosts[0]
	}

	return host, err
}

func SortQuickStats(quickStats []*m.QuickStats) []*m.QuickStats {
	chartbeatDebugger.Println("Sorting quickstats ...")
	sort.Sort(QuickStatsSort(quickStats))
	return quickStats
}
