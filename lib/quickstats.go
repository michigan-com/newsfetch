package lib

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type QuickStatsSnapshot struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Created_at time.Time     `bson:"created_at"`
	Stats      []*QuickStats `bson:"stats"`
}

type QuickStatsResp struct {
	Data *QuickStatsRespStats `bson:"data"`
}

type QuickStatsRespStats struct {
	Stats *QuickStats `bson:"stats"`
}

type QuickStats struct {
	Source   string        `bson:source`
	Visits   int           `bson:"visits"`
	Links    int           `bson:"links"`
	Direct   int           `bson:"direct"`
	Search   int           `bson:"search"`
	Social   int           `bson:"social"`
	Platform PlatformStats `bson:"platform"`
}

type PlatformStats struct {
	M int `bson:"m"`
	D int `bson:"d"`
}

type QuickStatsSort []*QuickStats

func (q QuickStatsSort) Len() int           { return len(q) }
func (q QuickStatsSort) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q QuickStatsSort) Less(i, j int) bool { return q[i].Visits > q[j].Visits }

func FetchQuickStats(urls []string) []*QuickStats {

	var urlWait sync.WaitGroup
	statQueue := make(chan *QuickStats, len(urls))

	for _, url := range urls {
		urlWait.Add(1)

		go func(url string) {
			stats, err := GetQuickStats(url)

			if err != nil {
				Debugger.Println("ERROR: %v", err)
			} else {
				statQueue <- stats
			}
			urlWait.Done()
		}(url)
	}

	urlWait.Wait()
	close(statQueue)

	quickStats := make([]*QuickStats, 0, len(urls))
	for stats := range statQueue {
		quickStats = append(quickStats, stats)
	}

	return SortQuickStats(quickStats)
}

func GetQuickStats(url string) (*QuickStats, error) {
	Debugger.Println("Fetching %s", url)

	// Parse out the host we're getting data on
	host, err := GetHostFromParams(url)
	if err != nil {
		Debugger.Printf("ERROR: %v", err)
		Debugger.Printf("Host will be \"\"")
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	Debugger.Println("Successfully fetched %s", url)

	quickStatsResp := &QuickStatsResp{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&quickStatsResp)

	if err != nil {
		return nil, err
	}

	quickStats := quickStatsResp.Data.Stats
	quickStats.Source = strings.Replace(host, ".com", "", -1)

	return quickStats, err
}

func SaveQuickStats(quickStats []*QuickStats, mongoUri string) {

	session := DBConnect(mongoUri)
	defer DBClose(session)

	quickStatsCol := session.DB("").C("Quickstats")

	// Insert this snapshot
	snapshot := QuickStatsSnapshot{}
	snapshot.Created_at = time.Now()
	snapshot.Stats = quickStats
	err := quickStatsCol.Insert(snapshot)

	if err != nil {
		Debugger.Printf("ERROR: %v", err)
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
		Debugger.Printf("Error while removing old quickstats snapshots %v", err)
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

func SortQuickStats(quickStats []*QuickStats) []*QuickStats {
	Debugger.Println("Sorting quickstats ...")
	sort.Sort(QuickStatsSort(quickStats))
	return quickStats
}
