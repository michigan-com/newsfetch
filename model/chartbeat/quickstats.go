package model

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type QuickStatsSnapshot struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Created_at time.Time     `bson:"created_at"`
	Stats      []*QuickStats `bson:"stats"`
}

// Snapshot interface
func (q QuickStatsSnapshot) Save(session *mgo.Session) {
	quickStatsCol := session.DB("").C("Quickstats")
	err := quickStatsCol.Insert(q)

	if err != nil {
		debugger.Printf("ERROR: %v", err)
		return
	}
	removeOldSnapshots(quickStatsCol)

	q.saveMobileSeries(session)
}

func (q QuickStatsSnapshot) saveMobileSeries(session *mgo.Session) {
	// These times are right now
	_today := time.Now()
	_tomorrow := _today.Add(24 * time.Hour)
	estLocation, _ := time.LoadLocation("EST")

	// These times are at midnight
	today := time.Date(_today.Year(), _today.Month(), _today.Day(), 0, 0, 0, 0, estLocation)
	tomorrow := time.Date(_tomorrow.Year(), _tomorrow.Month(), _tomorrow.Day(), 0, 0, 0, 0, estLocation)

	debugger.Printf("today: %v, tomorrow: %v", today, tomorrow)
	mobileSeriesCol := session.DB("").C("MobileSeries")
	mobileSeries := &MobileSeries{}

	// Find one if it exists
	query := mobileSeriesCol.Find(bson.M{
		"date": bson.M{
			"$gte": today,
			"$lt":  tomorrow,
		},
	}).
		Sort("-_id")

	//debugger.Printf("%v", query)
	query.One(&mobileSeries)

	if !mobileSeries.Id.Valid() {
		// This means there's no mobile series for today
		mobileSeries.Date = today
		mobileSeries.StartTime = _today
	}

	// Compile the mobile's total
	mobileTotal := 0
	for _, stat := range q.Stats {
		mobileTotal += stat.PlatformEngaged.M
	}
	mobileSeries.AddSeriesValue(mobileTotal)
	mobileSeries.Save(session)
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
