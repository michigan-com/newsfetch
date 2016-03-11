package model

import (
  "fmt"
  "time"

	"gopkg.in/mgo.v2/bson"
)

type Referrer struct {
  Site string `bson:"site"`
  TotalViewers float64 `bson:"visitors"`
  PublicationsCount []bson.M `bson:"publicationsCount"`
}

type HistoricalEntry struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Timestamp time.Time     `bson:"timestamp"`
	Referrers    []Referrer   `bson:"referrers"`
}
