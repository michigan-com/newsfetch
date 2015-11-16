package fetch

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"

	m "github.com/michigan-com/newsfetch/model/chartbeat"
)

type Quickstats struct{}

type QuickStatsSort []*m.QuickStats

func (q QuickStatsSort) Len() int           { return len(q) }
func (q QuickStatsSort) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q QuickStatsSort) Less(i, j int) bool { return q[i].Visits > q[j].Visits }

func (q Quickstats) Fetch(urls []string, session *mgo.Session) m.Snapshot {

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

	snapshot := m.QuickStatsSnapshot{}
	snapshot.Created_at = time.Now()
	snapshot.Stats = SortQuickStats(quickStats)
	return snapshot
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
		chartbeatError.Printf("Failed to fetch url %s: %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	chartbeatDebugger.Println("Successfully fetched %s", url)

	quickStatsResp := &m.QuickStatsResp{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&quickStatsResp)

	if err != nil {
		chartbeatError.Printf("Failed to parse json body from url %s: %v", url, err)
		return nil, err
	}

	quickStats := quickStatsResp.Data.Stats
	quickStats.Source = strings.Replace(host, ".com", "", -1)

	return quickStats, err
}

func SortQuickStats(quickStats []*m.QuickStats) []*m.QuickStats {
	sort.Sort(QuickStatsSort(quickStats))
	return quickStats
}
