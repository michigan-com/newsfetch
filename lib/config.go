package lib

import (
	"fmt"
	"github.com/spf13/viper"
	//"io/ioutil"
	"path"
	"runtime"
	"runtime/debug"
)

func ErrorCheck(err error) {
	if err != nil {
		fmt.Println(debug.Stack())
		panic(err)
	}
}

func Config(prop string) {

	_, filename, _, _ := runtime.Caller(1)
	baseDir := path.Join(path.Dir(filename), "..")

	// Set up config
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(baseDir)
	fmt.Println(baseDir)

	// Read the file
	fmt.Println("reading in config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}
