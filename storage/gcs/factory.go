package gcs

import (
	gcsStorage "cloud.google.com/go/storage"
	"fmt"
	"github.com/luizbafilho/helm-chart-publisher/storage"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func init() {
	storage.Register("gcs", NewGcsStorage)
}

func NewGcsStorage(conf map[string]interface{}) (storage.Storage, error) {
	config, err := decodeAndValidateConfig(conf)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	client, err := gcsStorage.NewClient(ctx)
	fmt.Println("ctx: ", ctx)
	fmt.Println("client: ", client)
	fmt.Println("err: ", err)

	return &GcsStore{
		name:   "gcs",
		config: config,
		gcs:    client,
	}, nil
}

func decodeAndValidateConfig(c map[string]interface{}) (*Config, error) {
	config := Config{}
	if err := mapstructure.Decode(c, &config); err != nil {
		return nil, err
	}
	fmt.Println("c: ", c)

	mergeEnviromentVariables(&config)
	fmt.Println("config: ", config)

	if config.project == "" {
		return nil, errors.New("Invalid config. project is required.")
	}
	if config.bucket == "" {
		return nil, errors.New("Invalid config. bucket is required.")
	}

	return &config, nil
}

func mergeEnviromentVariables(config *Config) {
	config.project = "staging-165617"
	config.bucket = "helm-charts-staging-165617"
}
