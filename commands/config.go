package commands

import (
	"github.com/kelseyhightower/envconfig"
)

type GlobalConfig struct {
	MongoUrl string `envconfig:"mongo_uri" required:"true"`
}

var globalConfig GlobalConfig

func loadConfig() {
	err := envconfig.Process("newsfetch", &globalConfig)
	if err != nil {
		panic(err)
	}

	println("MongoUrl (env) =", globalConfig.MongoUrl)
}
