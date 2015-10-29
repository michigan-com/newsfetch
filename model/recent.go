package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type RecentSnapshot struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Created_at time.Time     `bson:"created_at"`
	Recents    []*RecentResp `bson"recents"`
}

type RecentResp struct {
	Source  string
	Recents []Recent
}

type Recent struct {
	Lat   float32 `json:"lat" bson:"lat"`
	Lng   float32 `json:"lng" bson:"lng"`
	Title string  `json:"title" bson:"title"`
	Url   string  `json:"path" bson"url"`
	Host  string  `json:"domain" bson:"host"`
}
