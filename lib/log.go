package lib

import (
	"github.com/op/go-logging"
	"os"
)

var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfile} ▶ %{level:.4s} %{color:reset} %{message}",
)

func GetLogger() *logging.Logger {
	log := logging.MustGetLogger("newsfetch")
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	logging.SetBackend(backendFormatter)

	return log
}
