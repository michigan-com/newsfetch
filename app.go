package main

import (
	"github.com/michigan-com/newsfetch/lib"
	"time"
)

var log = lib.GetLogger()

func run() {

	for {
		log.Info("Running loop")

		FetchArticles()

		log.Info("Sleeping, don't bother me")
		time.Sleep(10 * time.Minute) // Sleep for 10 minutes
	}
}

func main() {
	run()
}
