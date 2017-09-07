package gcs

import (
	gcsStorage "cloud.google.com/go/storage"
	"fmt"
	"github.com/luizbafilho/helm-chart-publisher/storage"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"os"
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

	mergeEnviromentVariables(&config)

	if config.GoogleApplicationCredentials == "" {
		return nil, errors.New("Invalid config. GoogleApplicationCredentials is required.")
	}
	if config.Project == "" {
		return nil, errors.New("Invalid config. Project is required.")
	}
	if config.Bucket == "" {
		return nil, errors.New("Invalid config. Bucket is required.")
	}

	return &config, nil
}

func mergeEnviromentVariables(config *Config) {
	if googleApplicationCredentials := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); googleApplicationCredentials != "" {
		config.GoogleApplicationCredentials = googleApplicationCredentials
	}

}
