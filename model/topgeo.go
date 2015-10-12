package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type TopGeoSnapshot struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Created_at time.Time     `bson:"created_at"`
	Cities     []*TopGeo     `bson:"cities"`
}

type TopGeoResp struct {
	Geo TopGeo `bson:"geo:"`
}

type TopGeo struct {
	Source string `bson:"source"`
	Cities bson.M `bson:"cities"`
}
