package model

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type RecentSnapshot struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Created_at time.Time     `bson:"created_at"`
	Recents    []*RecentResp `bson:"recents"`
}

func (r RecentSnapshot) Save(session *mgo.Session) {
	col := session.DB("").C("Recent")
	err := col.Insert(r)

	if err != nil {
		debugger.Printf("Failed to insert Recent Snapshot: %v", err)
	}
	removeOldSnapshots(col)
}

type RecentResp struct {
	Source  string
	Recents []Recent
}

type Recent struct {
	Lat      float32 `json:"lat" bson:"lat"`
	Lng      float32 `json:"lng" bson:"lng"`
	Title    string  `json:"title" bson:"title"`
	Url      string  `json:"path" bson"url"`
	Host     string  `json:"domain" bson:"host"`
	Platform string
}
