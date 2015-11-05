package model

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"

  "github.com/michigan-com/newsfetch/lib"
)

var debugger = lib.NewCondLogger("newsfetch:model:chartbeat")

type Snapshot interface {
  Save(session *mgo.Session)
}

func removeOldSnapshots(col *mgo.Collection) {
  var snapshot = bson.M{
    "_id": -1,
  }
  // Remove old snapshots
  col.Find(bson.M{}).
    Select(bson.M{"_id": 1}).
    Sort("-_id").
    One(&snapshot)

  _, err := col.RemoveAll(bson.M{
    "_id": bson.M{
      "$ne": snapshot["_id"],
    },
  })

  if err != nil {
    debugger.Printf("Error while removing old quickstats snapshots %v", err)
  }
}