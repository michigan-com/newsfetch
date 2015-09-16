package main

import (
	"runtime"

	//"github.com/davecheney/profile"
	"github.com/michigan-com/newsfetch/commands"
	"github.com/op/go-logging"
)

var VERSION string

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()
	runtime.GOMAXPROCS(runtime.NumCPU())

	logging.SetLevel(logging.CRITICAL, "newsfetch")

	//VERSION is set in our build step
	commands.Execute(VERSION)
}
