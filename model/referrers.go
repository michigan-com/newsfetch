package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type ReferrersSnapshot struct {
	Id         bson.ObjectId `bson:"_id"`
	Created_at time.Time     `bson:"created_at"`
	Referrers  []*Referrers  `bson:"referrers"`
}

type Referrers struct {
	Source    string `json:"source"`
	Referrers bson.M `bson:"referrers" json:"referrers"`
}
