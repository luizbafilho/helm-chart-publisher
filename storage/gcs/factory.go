package gcs

import (
	gcsStorage "cloud.google.com/go/storage"
	"github.com/luizbafilho/helm-chart-publisher/storage"
	"github.com/mitchellh/mapstructure"
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
	// requires GCloud "Application Default Credentials" set via
	//  `gcloud auth application-default login` command or
	//  GOOGLE_APPLICATION_CREDENTIALS environment variable
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

	return &config, nil
}

