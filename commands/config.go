package commands

import (
	"github.com/kelseyhightower/envconfig"
)

type GlobalConfig struct {
	MongoUrl        string `envconfig:"mongo_uri"`
	ChartbeatApiKey string `envconfig:"chartbeat_api_key"`
	GnapiDomain     string `envconfig:"gnapi_domain"`
}

var globalConfig GlobalConfig

func loadConfig() {
	err := envconfig.Process("newsfetch", &globalConfig)
	if err != nil {
		panic(err)
	}
}
