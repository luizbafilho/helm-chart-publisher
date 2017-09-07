package gcs

import (
	"fmt"

	gcsStorage "cloud.google.com/go/storage"
	"context"
	"github.com/luizbafilho/helm-chart-publisher/storage"
	"io/ioutil"
)

type Config struct {
	GoogleApplicationCredentials string
	Project                      string
	Bucket                       string
}

type GcsStore struct {
	name   string
	config *Config
	gcs    *gcsStorage.Client
}

// Name ...
func (s *GcsStore) Name() string {
	return s.name
}

func parseError(err error) error {
	if err == gcsStorage.ErrObjectNotExist {
		return storage.PathNotFoundErr{}
	}
	return err
}

// Get ...
func (s *GcsStore) Get(bucket string, path string) (*storage.GetResponse, error) {
	gcsBucket := s.gcs.Bucket(bucket)
	obj := gcsBucket.Object(path)

	ctx := context.Background()
	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, parseError(err)
	}
	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return &storage.GetResponse{
		Body: body,
	}, nil
}

// Put stores the content
func (s *GcsStore) Put(bucket string, path string, content []byte) (*storage.PutResponse, error) {
	// This authentication is done because there is some weird bug when authentication retry
	// has to happen, making the swift package send a request with zero byte body.
	// To avoid that, I'm authenticating before the PUT call, to make sure no retry will be need.
	gcsBucket := s.gcs.Bucket(bucket)
	obj := gcsBucket.Object(path)

	ctx := context.Background()
	w := obj.NewWriter(ctx)

	// Write some text to obj. This will overwrite whatever is there.
	if _, err := w.Write(content); err != nil {
		// TODO: Handle error.
	}
	// Close, just like writing a file.
	if err := w.Close(); err != nil {
		// TODO: Handle error.
	}

	return &storage.PutResponse{}, nil
}

// GetURL ...
func (s *GcsStore) GetURL(bucket string, path string) string {
	domain := "storage.cloud.google.com"
	return fmt.Sprintf("%s/%s/%s", domain, bucket, path)
}
