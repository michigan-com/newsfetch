package main

import (
	"runtime"

	//"github.com/davecheney/profile"
	"github.com/michigan-com/newsfetch/commands"
)

// Version number that gets compiled via `make build` or `make install`
var VERSION string

// Git commit hash that gets compiled via `make build` or `make install`
var COMMITHASH string

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()
	runtime.GOMAXPROCS(runtime.NumCPU())

	//VERSION is set in our build step
	commands.Execute(VERSION, COMMITHASH)
}
