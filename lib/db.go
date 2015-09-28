package lib

import (
	"os"

	"gopkg.in/mgo.v2"
)

func DBConnect(uri string) *mgo.Session {
	// TODO read this from config
	session, err := mgo.Dial(uri)
	if err != nil {
		Logger.Printf("Failed to connect to '%s': %v", uri, err)
		os.Exit(1)
	}

	session.SetMode(mgo.Monotonic, true)
	return session
}

func DBClose(session *mgo.Session) {
	session.Close()
}
