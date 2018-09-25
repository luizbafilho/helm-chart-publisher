package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
)

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

func ReadConfigFile(configFile string) error {
	if configFile == "" {
		return errors.New("no configuration file specified")
	}

	file, err := os.Open(configFile)
	if err != nil {
		return errors.Wrap(err, "open config file failed")
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(file)
	if err != nil {
		return errors.Wrap(err, "read config file failed")
	}

	return nil
}
