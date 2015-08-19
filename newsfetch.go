package main

import (
	"runtime"

	"github.com/michigan-com/newsfetch/commands"
	"github.com/op/go-logging"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	logging.SetLevel(logging.CRITICAL, "newsfetch")

	commands.Execute()
}
