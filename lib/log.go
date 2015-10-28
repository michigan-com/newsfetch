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
	name    string
	enabled bool
	*log.Logger
}

func (c *CondLogger) Enable() {
	c.enabled = true
	c.Logger.SetOutput(os.Stdout)
}

func (c *CondLogger) Disable() {
	c.enabled = false
	c.Logger.SetOutput(ioutil.Discard)
}

func (c *CondLogger) IsEnabled() bool {
	return c.enabled
}

func NewCondLogger(name string) *CondLogger {
	prefix := fmt.Sprintf("(newsfetch:%s) ", name)
	logger := &CondLogger{name, false, log.New(ioutil.Discard, prefix, log.Lshortfile)}

	if shouldEnableDebugger(name) {
		logger.Enable()
	}

	return logger
}

func shouldEnableDebugger(name string) bool {
	name = strings.ToLower(name)

	for _, bugger := range debuggers {
		if bugger == "*" || bugger == name {
			return true
		}
	}

	return false
}

// This logger should be used for essential information that will be viewable by
// our centralized logging service. Use this for production only.
var Logger = NewCondLogger("logger")

// This logger should be use for general purpose debugging.  Use this for development only.
var Debugger = NewCondLogger("debugger")
