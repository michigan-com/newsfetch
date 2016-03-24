package model

import (
	"time"
	"errors"
  "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ReferrersSnapshot struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Created_at time.Time     `bson:"created_at"`
	Referrers  []*Referrers  `bson:"referrers"`
	Expire_at  time.Time     `bson:"expire_at"`
}

func (r ReferrersSnapshot) Save(session *mgo.Session) {
  collection := session.DB("").C("Referrers")
  historyCollection := session.DB("").C("ReferrerHistory")
	r.Expire_at = r.Created_at.Add(time.Duration(30)* time.Second)
  err := collection.Insert(r)
  if err != nil {
    debugger.Printf("Failed to insert Referrers snapshot: %v", err)
    return
  }

	latest := &ReferrersSnapshot{};

	err = historyCollection.Find(bson.M{}).Sort("-created_at").One(latest)
	fiveMinutesAgo := time.Now().Add(-time.Duration(5)* time.Minute)
	if err == errors.New("not found") || latest.Created_at.Before(fiveMinutesAgo) {
    debugger.Printf("Saved a Snapshot to ReferrerHistory Collection")
		r.Expire_at = r.Expire_at.Add(time.Duration(24*7)* time.Hour)
		historyCollection.Insert(r)
	}
	// remove once indexes are inforced
  removeOldSnapshots(collection)
}

type Referrers struct {
	Source    string `json:"source"`
	Referrers bson.M `bson:"referrers" json:"referrers"`
}
