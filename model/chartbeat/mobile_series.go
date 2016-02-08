package model

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MobileSeries struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Date      time.Time     `bson:"date"`
	StartTime time.Time     `bson:"startTime"`
	Series    []int         `bson:"series"`
}

func (m *MobileSeries) AddSeriesValue(value int) {
	//debugger.Printf("Series while in AddSeriesValue: %v", m.Series)
	seriesSlice := m.Series[0:len(m.Series)]
	seriesSlice = append(seriesSlice, value)
	m.Series = seriesSlice
}

func (m MobileSeries) Save(session *mgo.Session) {
	col := session.DB("").C("MobileSeries")
	selector := bson.M{}

	info, err := col.Upsert(selector, m)
	if err != nil {
		debugger.Printf("Error upserting MobileSeries: %v", err)
	}

	debugger.Printf("%v", info)
}
