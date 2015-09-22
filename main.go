package main

import (
	"runtime"

	//"github.com/davecheney/profile"
	"github.com/michigan-com/newsfetch/commands"
)

var VERSION string

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()
	runtime.GOMAXPROCS(runtime.NumCPU())

	//VERSION is set in our build step
	commands.Execute(VERSION)
}
