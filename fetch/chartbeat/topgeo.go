package fetch

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	m "github.com/michigan-com/newsfetch/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func FetchTopGeo(urls []string) []*m.TopGeo {
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

	return topGeos
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

func SaveTopGeo(topGeos []*m.TopGeo, session *mgo.Session) {
	topGeoCol := session.DB("").C("Topgeo")

	snapshot := m.TopGeoSnapshot{}
	snapshot.Created_at = time.Now()
	snapshot.Cities = topGeos
	err := topGeoCol.Insert(snapshot)

	if err != nil {
		chartbeatDebugger.Printf("ERROR: %v", err)
	}

	topGeoCol.Find(bson.M{}).
		Select(bson.M{"_id": 1}).
		Sort("-_id").
		One(&snapshot)

	_, err = topGeoCol.RemoveAll(bson.M{
		"_id": bson.M{
			"$ne": snapshot.Id,
		},
	})

	if err != nil {
		chartbeatDebugger.Printf("Error while removing old topgeo snapshots: %v", err)
	}
}
