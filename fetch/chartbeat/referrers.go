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

type Referrers struct {}

func (r Referrers) Fetch(urls []string, session *mgo.Session) m.Snapshot {
	var wait sync.WaitGroup
	statQueue := make(chan *m.Referrers, len(urls))

	for _, url := range urls {
		wait.Add(1)

		go func(url string) {
			ref, err := getReferrers(url)
			if err != nil {
				chartbeatDebugger.Printf("Error fetching %s:\n%v", url, err)
			} else {
				statQueue <- ref
			}

			wait.Done()
		}(url)
	}

	wait.Wait()
	close(statQueue)

	referrers := make([]*m.Referrers, 0, len(urls))
	for ref := range statQueue {
		referrers = append(referrers, ref)
	}

	snapshot := m.ReferrersSnapshot{}
	snapshot.Created_at = time.Now()
	snapshot.Referrers = referrers
	return snapshot
}

func getReferrers(url string) (*m.Referrers, error) {
	chartbeatDebugger.Printf("Fetching %s", url)

	resp, err := http.Get(url)
	if err != nil {
		chartbeatError.Printf("Failed to get url %s: %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	host, err := GetHostFromParams(url)
	if err != nil {
		chartbeatDebugger.Printf("ERROR: %v", err)
		chartbeatDebugger.Printf("Host will be \"\"")
	}

	referrers := &m.Referrers{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&referrers)
	if err != nil {
		chartbeatError.Printf("Failed to json parse json body from url %s: %v", url, err)
		return nil, err
	}

	referrers.Source = strings.Replace(host, ".com", "", -1)

	return referrers, nil
}
