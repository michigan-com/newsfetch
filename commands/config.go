package commands

import (
	"github.com/kelseyhightower/envconfig"
)

type GlobalConfig struct {
	MongoUrl string
}

var globalConfig GlobalConfig

func loadConfig() {
	err := envconfig.Process("newsfetch", &globalConfig)
	if err != nil {
		panic(err)
	}

	// println("MongoUrl =", globalConfig.MongoUrl)
}
