package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/michigan-com/newsfetch/lib"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var debugger = lib.NewCondLogger("HistoricalTraffic")

type HistoricalIn struct {
	Data struct {
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
	return fmt.Sprintf("<HistoricalIn %d-%d>", h.Data.Start, h.Data.End)
}

func (h *HistoricalIn) SignalMapi() {
	resp, err := http.Get("https://api.michigan.com/historical-traffic/")
	if err != nil {
		debugger.Println(err)
	} else {
		defer resp.Body.Close()
		now := time.Now()
		debugger.Printf("Updated snapshot at Mapi at %s", now)
	}
}

func (h *HistoricalIn) Run(session *mgo.Session, apiKey string) {
	debugger.Println("RUNNING HISTORICAL TRAFFIC")

	urls, err := FormatChartbeatUrls("historical/traffic/series", lib.Sites, apiKey)
	if err != nil {
		debugger.Println(err)
		return
	}

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			debugger.Printf("Failed to fetch url: %s: %s", url, err.Error())
		}

		tmpHI := &HistoricalIn{}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(tmpHI)

		if err != nil {
			debugger.Printf("Failed to parse json body: %s", err.Error())
		}

		h.Data.Start = tmpHI.Data.Start
		h.Data.End = tmpHI.Data.End
		h.Data.Frequency = tmpHI.Data.Frequency
		h.CombineSeries(tmpHI)
	}

	debugger.Println(h)

	if session == nil {
		debugger.Println("No mongo session found, skipping save")
		return
	}

	snapshot := &HistoricalSnapshot{
		Start:     h.Data.Start,
		End:       h.Data.End,
		Frequency: h.Data.Frequency,
	}

	// merge all data into mongo model
	snapshot.Merge(h)
	// save snapshot data to mongo
	snapshot.Save(session)
	// send signal to mapi that there's new data
	h.SignalMapi()
}

func (h *HistoricalIn) CombineSeries(hi *HistoricalIn) {
	if hi.Data.Freep != nil {
		h.Data.Freep = hi.Data.Freep
	} else if hi.Data.DetroitNews != nil {
		h.Data.DetroitNews = hi.Data.DetroitNews
	} else if hi.Data.BattleCreek != nil {
		h.Data.BattleCreek = hi.Data.BattleCreek
	} else if hi.Data.Hometown != nil {
		h.Data.Hometown = hi.Data.Hometown
	} else if hi.Data.Lansing != nil {
		h.Data.Lansing = hi.Data.Lansing
	} else if hi.Data.Livingston != nil {
		h.Data.Livingston = hi.Data.Livingston
	} else if hi.Data.Herald != nil {
		h.Data.Herald = hi.Data.Herald
	}
}

type HistoricalInSeries struct {
	Series *struct {
		People []int `json:"people"`
	} `json:"series"`
}

func (his *HistoricalInSeries) Visits() []int {
	return his.Series.People
}

type HistoricalSnapshot struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Created_at time.Time     `bson:"created_at"`
	Start      int           `bson:"start"`
	End        int           `bson:"end"`
	Frequency  int           `bson:"frequency"`
	Traffic    []*Traffic    `bson:"sites"`
}

type Traffic struct {
	Site   string `bson:"site"`
	Visits []int  `bson:"visits"`
}

func (h *HistoricalSnapshot) Merge(hi *HistoricalIn) {
	h.Traffic = []*Traffic{
		&Traffic{"freep", hi.Data.Freep.Visits()},
		&Traffic{"detroitnews", hi.Data.DetroitNews.Visits()},
		&Traffic{"battlecreekenquirer", hi.Data.BattleCreek.Visits()},
		&Traffic{"hometownlife", hi.Data.Hometown.Visits()},
		&Traffic{"lansingstatejournal", hi.Data.Lansing.Visits()},
		&Traffic{"livingstondaily", hi.Data.Livingston.Visits()},
		&Traffic{"thetimesherald", hi.Data.Herald.Visits()},
	}
}

func (h *HistoricalSnapshot) Save(session *mgo.Session) {
	if session == nil {
		return
	}

	collection := session.DB("").C("HistoricalTraffic")
	_, err := collection.RemoveAll(bson.M{})
	if err != nil {
		debugger.Println(err)
	}

	h.Created_at = time.Now()

	err = collection.Insert(h)
	if err != nil {
		debugger.Println(err)
	}
}
