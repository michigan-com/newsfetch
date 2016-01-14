package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"

	m "github.com/michigan-com/newsfetch/model/chartbeat"
)

type TrafficSeries struct {
	Data struct {
		Start int `json:"start"`
		End   int `json:"end"`
		// frequency is the data sample interval in minutes
		Frequency int `json:"frequency"`

		// Michigan.com
		Freep       *m.TrafficSeriesIn `json:"freep.com"`
		DetroitNews *m.TrafficSeriesIn `json:"detroitnews.com"`
		BattleCreek *m.TrafficSeriesIn `json:"battlecreekenquirer.com"`
		Hometown    *m.TrafficSeriesIn `json:"hometownlife.com"`
		Lansing     *m.TrafficSeriesIn `json:"lansingstatejournal.com"`
		Livingston  *m.TrafficSeriesIn `json:"livingstondaily.com"`
		Herald      *m.TrafficSeriesIn `json:"thetimesherald.com"`

		// Usat
		UsaToday *m.TrafficSeriesIn `json:"usatoday.com"`

		// Tennessean
		Tennessean *m.TrafficSeriesIn `json:"tennessean.com"`

		// Central Ohio omg why are there so many sites help
		Mansfield        *m.TrafficSeriesIn `json:"mansfieldnewsjournal.com"`
		Newark           *m.TrafficSeriesIn `json:"newarkadvocate.com"`
		Zanesville       *m.TrafficSeriesIn `json:"zanesvilletimesrecorder.com"`
		Chillicothe      *m.TrafficSeriesIn `json:"chillicothegazette.com"`
		Lancaster        *m.TrafficSeriesIn `json:"lancastereaglegazette.com"`
		Marion           *m.TrafficSeriesIn `json:"marionstar.com"`
		TheNewsMessenger *m.TrafficSeriesIn `json:"thenews-messenger.com"`
		Coshocton        *m.TrafficSeriesIn `json:"coshoctontribune.com"`
		Bucyrus          *m.TrafficSeriesIn `json:"bucyrustelegraphforum.com"`
		PortClinton      *m.TrafficSeriesIn `json:"portclintonnewsherald.com"`
	} `json:"data"`
}

func (h *TrafficSeries) GetSeries() *m.TrafficSeriesIn {
	if h.Data.Freep != nil {
		return h.Data.Freep
	} else if h.Data.DetroitNews != nil {
		return h.Data.DetroitNews
	} else if h.Data.BattleCreek != nil {
		return h.Data.BattleCreek
	} else if h.Data.Hometown != nil {
		return h.Data.Hometown
	} else if h.Data.Lansing != nil {
		return h.Data.Lansing
	} else if h.Data.Livingston != nil {
		return h.Data.Livingston
	} else if h.Data.Herald != nil {
		return h.Data.Herald
	} else if h.Data.UsaToday != nil {
		return h.Data.UsaToday
	} else if h.Data.Tennessean != nil {
		return h.Data.Tennessean
	} else if h.Data.Mansfield != nil {
		return h.Data.Mansfield
	} else if h.Data.Newark != nil {
		return h.Data.Newark
	} else if h.Data.Zanesville != nil {
		return h.Data.Zanesville
	} else if h.Data.Chillicothe != nil {
		return h.Data.Chillicothe
	} else if h.Data.Lancaster != nil {
		return h.Data.Lancaster
	} else if h.Data.Marion != nil {
		return h.Data.Marion
	} else if h.Data.TheNewsMessenger != nil {
		return h.Data.TheNewsMessenger
	} else if h.Data.Coshocton != nil {
		return h.Data.Coshocton
	} else if h.Data.Bucyrus != nil {
		return h.Data.Bucyrus
	} else if h.Data.PortClinton != nil {
		return h.Data.PortClinton
	}

	return nil
}

func (h *TrafficSeries) String() string {
	return fmt.Sprintf("<TrafficSeries %d-%d>", h.Data.Start, h.Data.End)
}

func (h *TrafficSeries) SignalMapi() {
	resp, err := http.Get("https://api.michigan.com/traffic-series/")
	if err != nil {
		chartbeatDebugger.Println(err)
	} else {
		defer resp.Body.Close()
		now := time.Now()
		chartbeatDebugger.Printf("Updated snapshot at Mapi at %s", now)
	}
}

func (h TrafficSeries) Fetch(urls []string, session *mgo.Session) m.Snapshot {
	var wait sync.WaitGroup
	var start int
	var end int
	var freq int
	queue := make(chan *m.Traffic, len(urls))

	for _, url := range urls {
		wait.Add(1)

		go func(url string) {
			chartbeatDebugger.Printf("Fetching %s", url)
			resp, err := http.Get(url)
			if err != nil {
				chartbeatDebugger.Printf("Failed to fetch url: %s: %s", url, err.Error())
				wait.Done()
				return
			}
			defer resp.Body.Close()

			chartbeatDebugger.Printf("DOne fetching %s", url)

			tmpHI := &TrafficSeries{}
			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(tmpHI)

			chartbeatDebugger.Printf("JSON decoded")

			if err != nil {
				chartbeatDebugger.Printf("Failed to parse json body: %s", err.Error())
				wait.Done()
				return
			}

			visits := tmpHI.GetSeries()
			if visits == nil {
				chartbeatDebugger.Printf("Failed to pull a visits from response")
				wait.Done()
				return
			}

			series := &m.Traffic{}
			source, _ := GetHostFromParams(url)

			series.Source = strings.Replace(source, ".com", "", -1)
			series.Visits = visits.Visits()
			queue <- series

			start = tmpHI.Data.Start
			end = tmpHI.Data.End
			freq = tmpHI.Data.Frequency
			wait.Done()

		}(url)
	}

	wait.Wait()
	close(queue)

	// Get the values out of the queue
	trafficSlice := make([]*m.Traffic, 0, len(urls))
	for traffic := range queue {
		trafficSlice = append(trafficSlice, traffic)
	}

	snapshot := m.TrafficSeriesSnapshot{
		Start:      start,
		End:        end,
		Frequency:  freq,
		Traffic:    trafficSlice,
		Created_at: time.Now(),
	}

	return snapshot
}
