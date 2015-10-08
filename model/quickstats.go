package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type QuickStatsSnapshot struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Created_at time.Time     `bson:"created_at"`
	Stats      []*QuickStats `bson:"stats"`
}

type QuickStatsResp struct {
	Data *QuickStatsRespStats `bson:"data"`
}

type QuickStatsRespStats struct {
	Stats *QuickStats `bson:"stats"`
}

type QuickStats struct {
	Source   string        `bson:source`
	Visits   int           `bson:"visits"`
	Links    int           `bson:"links"`
	Direct   int           `bson:"direct"`
	Search   int           `bson:"search"`
	Social   int           `bson:"social"`
	Platform PlatformStats `bson:"platform"`
}

type PlatformStats struct {
	M int `bson:"m"`
	T int `bson:"t"`
	D int `bson:"d"`
	A int `bson:"a"`
}
