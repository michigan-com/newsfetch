package lib

import (
	"gopkg.in/mgo.v2"
)

var DB = DBConnect()

func DBConnect() *mgo.Database {
	// TODO read this from config
	session, err := mgo.Dial("mongodb://localhost:27017/")
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	return session.DB("mapi")
}

func DBClose(session *mgo.Session) {
	session.Close()
}
