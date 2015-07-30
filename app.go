package main

import (
	"log"
	"time"
)

func run() {

	for {
		log.Print("Running loop")

		FetchArticles()

		log.Print("Sleeping, don't bother me")
		time.Sleep(10 * time.Minute) // Sleep for 10 minutes
	}
}

func main() {
	run()
}
