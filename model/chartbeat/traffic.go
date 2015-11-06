package model

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TrafficSeriesSnapshot struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Created_at time.Time     `bson:"created_at"`
	Start      int           `bson:"start"`
	End        int           `bson:"end"`
	Frequency  int           `bson:"frequency"`
	Traffic    []*Traffic    `bson:"sites"`
}

func (h TrafficSeriesSnapshot) Save(session *mgo.Session) {
	collection := session.DB("").C("TrafficSeries")
	err := collection.Insert(h)

	debugger.Printf("Saving %v", h)

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

type TrafficSeriesIn struct {
	Series *struct {
		People []int `json:"people"`
	} `json:"series"`
}

func (his *TrafficSeriesIn) Visits() []int {
	return his.Series.People
}
