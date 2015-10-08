package model

import (
	"gopkg.in/mgo.v2"
)

type Document interface {
	Save(session *mgo.Session) error
}
