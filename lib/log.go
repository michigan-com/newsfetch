package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var Debug string = strings.ToLower(os.Getenv("DEBUG"))
var debuggers = strings.Split(Debug, ",")

type CondLogger struct {
	name string
	*log.Logger
}

func (c *CondLogger) Enable() {
	c.Logger.SetOutput(os.Stdout)
}

func (c *CondLogger) Disable() {
	c.Logger.SetOutput(ioutil.Discard)
}

func NewCondLogger(name string) *CondLogger {
	out := ioutil.Discard
	name = strings.ToLower(name)

	for _, bugger := range debuggers {
		if bugger == "*" || bugger == name {
			out = os.Stdout
			break
		}
	}

	prefix := fmt.Sprintf("(newsfetch:%s) ", name)
	return &CondLogger{name, log.New(out, prefix, log.Lshortfile)}
}

// This logger should be used for essential information that will be viewable by
// our centralized logging service. Use this for production only.
var Logger = NewCondLogger("logger")

// This logger should be use for general purpose debugging.  Use this for development only.
var Debugger = NewCondLogger("debugger")
