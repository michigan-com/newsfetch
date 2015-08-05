package lib

import (
	"gopkg.in/mgo.v2"
)

func DBConnect(uri string) *mgo.Session {
	// TODO read this from config
	session, err := mgo.Dial(uri)
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	return session
}

func DBClose(session *mgo.Session) {
	session.Close()
}
