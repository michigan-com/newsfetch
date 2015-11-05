package model

import (
	"time"

  "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ReferrersSnapshot struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Created_at time.Time     `bson:"created_at"`
	Referrers  []*Referrers  `bson:"referrers"`
}

func (r ReferrersSnapshot) Save(session *mgo.Session) {
  collection := session.DB("").C("Referrers")
  err := collection.Insert(r)

  if err != nil {
    debugger.Printf("Failed to insert Referrers snapshot: %v", err)
    return
  }

  removeOldSnapshots(collection)
}

type Referrers struct {
	Source    string `json:"source"`
	Referrers bson.M `bson:"referrers" json:"referrers"`
}
