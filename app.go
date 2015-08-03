package main

import (
	"github.com/michigan-com/newsFetch/newsFetch"
	"log"
	"time"
)

func run() {
	log.Print("Running loop")

	newsFetch.FetchArticles()

	log.Print("Sleeping again")
	time.Sleep(10 * time.Minute) // Sleep for 10 minutes
}

func main() {
	run()
}
