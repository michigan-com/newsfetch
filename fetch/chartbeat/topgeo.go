package fetch

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"

	m "github.com/michigan-com/newsfetch/model/chartbeat"
)

type TopGeo struct{}

func (t TopGeo) Fetch(urls []string, session *mgo.Session) m.Snapshot {
	var urlWait sync.WaitGroup
	geoQueue := make(chan *m.TopGeo, len(urls))

	for _, url := range urls {
		urlWait.Add(1)

		go func(url string) {
			stats, err := GetTopGeo(url)

			if err != nil {
				chartbeatDebugger.Println("ERROR: %v", err)
			} else {
				geoQueue <- stats
			}

			urlWait.Done()
		}(url)
	}

	urlWait.Wait()
	close(geoQueue)

	topGeos := make([]*m.TopGeo, 0, len(urls))
	for geo := range geoQueue {
		topGeos = append(topGeos, geo)
	}

	snapshot := m.TopGeoSnapshot{}
	snapshot.Created_at = time.Now()
	snapshot.Cities = topGeos
	return snapshot
}

func GetTopGeo(url string) (*m.TopGeo, error) {
	chartbeatDebugger.Printf("Fetching %s", url)

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

	topGeoResp := &m.TopGeoResp{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&topGeoResp)
	topGeo := &topGeoResp.Geo

	if err != nil {
		chartbeatError.Printf("Failed to parse json body from url %s: %v", url, err)
		return nil, err
	}

	topGeo.Source = strings.Replace(host, ".com", "", -1)
	return topGeo, nil
}