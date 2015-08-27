package lib

import (
	"github.com/op/go-logging"
	"os"
)

var Format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfile} ▶ %{level:.4s} %{color:reset} %{message}",
)

var Backend = logging.NewLogBackend(os.Stderr, "", 0)

func GetLogger() *logging.Logger {
	log := logging.MustGetLogger("newsfetch")
	backendFormatter := logging.NewBackendFormatter(Backend, Format)

	logging.SetBackend(backendFormatter)

	//env var trumps everything
	level := logging.DEBUG
	levelEnv := os.Getenv("LOGLEVEL")
	var err error
	if levelEnv != "" {
		level, err = logging.LogLevel(levelEnv)
		if err != nil {
			level = logging.DEBUG
		}
	}

	logging.SetLevel(level, "newsfetch")

	return log
}
