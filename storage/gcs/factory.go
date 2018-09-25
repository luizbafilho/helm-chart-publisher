package gcs

import (
	gcsStorage "cloud.google.com/go/storage"
	"github.com/luizbafilho/helm-chart-publisher/storage"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func init() {
	storage.Register("gcs", NewGcsStorage)
}

func NewGcsStorage(conf map[string]interface{}) (storage.Storage, error) {
	ctx := context.Background()
	// requires GCloud "Application Default Credentials" set via
	//  `gcloud auth application-default login` command or
	//  GOOGLE_APPLICATION_CREDENTIALS environment variable
	client, err := gcsStorage.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "[GCS] Client failure")
	}

	return &GcsStore{
		name: "gcs",
		gcs:  client,
	}, nil
}
