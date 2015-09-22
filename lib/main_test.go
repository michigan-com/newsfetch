package lib

import (
	"fmt"
	"os"
	"testing"
)

var (
	TestMongoUri string
)

func TestMain(m *testing.M) {
	TestMongoUri = os.Getenv("MONGO_URI")
	if TestMongoUri == "" {
		panic("Mongo URI not specified, please set the MONGOURI environment variable and try again")
	}

	val := m.Run()

	session := DBConnect(TestMongoUri)
	defer DBClose(session)
	err := session.DB("").DropDatabase()

	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	os.Exit(val)
}
