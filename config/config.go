package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const APP = "helm-charts-publisher"

func init() {
	initViper()
}

func initViper() {
	viper.SetConfigName("config")        // name of config file (without extension)
	viper.AddConfigPath("/etc/" + APP)   // path to look for the config file in
	viper.AddConfigPath("$HOME/." + APP) // call multiple times to add many search paths
	viper.AddConfigPath(".")             // optionally look for config in the working directory
	err := viper.ReadInConfig()          // Find and read the config file
	if err != nil {                      // Handle errors reading the config file
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

type Config map[string]interface{}

func GetStorage() (string, Config) {
	storageConfig := viper.GetStringMap("storage")

	var storages []string
	// Return only key in this map
	for k := range storageConfig {
		storages = append(storages, k)
	}

	if len(storages) != 1 {
		fmt.Println("Error: multiple or none storage specified in configuration")
		os.Exit(1)
	}

	return storages[0], viper.GetStringMap("storage." + storages[0])
}

func GetRepos() interface{} {
	return viper.Get("repos")
}
