package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

type HistoricalSnapshot struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Created_at time.Time     `bson:"created_at"`
	Start      int           `bson:"start"`
	End        int           `bson:"end"`
	Frequency  int           `bson:"frequency"`
	Traffic    []*Traffic    `bson:"sites"`
}

func (h HistoricalSnapshot) Save(session *mgo.Session) {
	collection := session.DB("").C("HistoricalTraffic")
	err := collection.Insert(h)

	if err != nil {
		debugger.Printf("Failed to insert Historical snapshot: %v", err)
		return
	}
	removeOldSnapshots(collection)
}

type Traffic struct {
	Source string `bson:"site"`
	Visits []int  `bson:"visits"`
}

type HistoricalInSeries struct {
	Series *struct {
		People []int `json:"people"`
	} `json:"series"`
}

func (his *HistoricalInSeries) Visits() []int {
	return his.Series.People
}