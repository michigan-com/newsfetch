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
	Source          string        `bson:source`
	Visits          int           `bson:"visits"`
	Links           int           `bson:"links"`
	Direct          int           `bson:"direct"`
	Search          int           `bson:"search"`
	Social          int           `bson:"social"`
	Recirc          int           `bson:"recirc"`
	Article         int           `bson:"article"`
	PlatformEngaged PlatformStats `json:"platform_engaged"bson:"platform_engaged"`
	Loyalty         LoyaltyStats  `bson:"loyalty"`
}

type PlatformStats struct {
	M int `bson:"m"`
	T int `bson:"t"`
	D int `bson:"d"`
	A int `bson:"a"`
}

type LoyaltyStats struct {
	New       int `bson:"new"`
	Loyal     int `bson:"loyal"`
	Returning int `bson:"returning"`
}
