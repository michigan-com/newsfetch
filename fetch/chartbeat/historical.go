package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/michigan-com/newsfetch/lib"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var debugger = lib.NewCondLogger("HistoricalTraffic")

type HistoricalIn struct {
	Data *struct {
		Start int `json:"start"`
		End   int `json:"end"`
		// frequency is the data sample interval in minutes
		Frequency   int                 `json:"frequency"`
		Freep       *HistoricalInSeries `json:"freep.com"`
		DetroitNews *HistoricalInSeries `json:"detroitnews.com"`
		BattleCreek *HistoricalInSeries `json:"battlecreekenquirer.com"`
		Hometown    *HistoricalInSeries `json:"hometownlife.com"`
		Lansing     *HistoricalInSeries `json:"lansingstatejournal.com"`
		Livingston  *HistoricalInSeries `json:"livingstondaily.com"`
		Herald      *HistoricalInSeries `json:"thetimesherald.com"`
	} `json:"data"`
}

func NewHistoricalIn() *HistoricalIn {
	return &HistoricalIn{}
}

func (h *HistoricalIn) String() string {
	return fmt.Sprintf("<HistoricalIn %d:%d>", h.Data.Start, h.Data.End)
}

func (h *HistoricalIn) Run(session *mgo.Session, apiKey string) {
	debugger.Println("RUNNING HISTORICAL TRAFFIC")

	urls, err := FormatChartbeatUrls("historical/traffic/series", lib.Sites, apiKey)
	if err != nil {
		debugger.Println(err)
		return
	}

	debugger.Println(urls)

	var rWait sync.WaitGroup
	for _, url := range urls {
		rWait.Add(1)

		go func() {
			defer rWait.Done()
			resp, err := http.Get(url)
			if err != nil {
				debugger.Printf("Failed to fetch url: %s: %s", url, err.Error())
			}

			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(h)

			if err != nil {
				debugger.Printf("Failed to parse json body: %s", err.Error())
			}
		}()
	}

	rWait.Wait()
	debugger.Println(h)
}

type HistoricalInSeries struct {
	Series *struct {
		People []int `json:"people"`
	} `json:"series"`
}

type HistoricalSnapshot struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Created_at time.Time     `bson:"created_at"`
	Start      int           `bson:"start"`
	End        int           `bson:"end"`
	Frequency  int           `bson:"frequency"`
	Traffic    []*struct {
		Site   string `bson:"site"`
		Visits []int  `bson:"visits"`
	} `bson:"sites"`
}

func (h *HistoricalSnapshot) Save(session *mgo.Session) {}
