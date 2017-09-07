package gcs

import (
	"fmt"

	gcsStorage "cloud.google.com/go/storage"
	"context"
	"github.com/luizbafilho/helm-chart-publisher/storage"
	"github.com/pkg/errors"
	"io/ioutil"

)

type Config struct {
	// for public repos, give allUsers the Storage Object Viewer role for bucket
	// configuration all done inside of 'Application Default Credentials'
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
		return nil, errors.Wrap(err, "[GCS] Read failed")
	}

	return &storage.GetResponse{
		Body: body,
	}, nil
}

// Put stores the content
func (s *GcsStore) Put(bucket string, path string, content []byte) (*storage.PutResponse, error) {
	gcsBucket := s.gcs.Bucket(bucket)
	obj := gcsBucket.Object(path)

	ctx := context.Background()
	w := obj.NewWriter(ctx)

	// Write some text to obj. This will overwrite whatever is there.
	if _, err := w.Write(content); err != nil {
		return nil, errors.Wrap(err, "[GCS] Write failed")
	}
	// Close, just like writing a file.
	if err := w.Close(); err != nil {
		return nil, errors.Wrap(err, "[GCS] Write failed on Close()")
	}

	return &storage.PutResponse{}, nil
}

// GetURL ...
func (s *GcsStore) GetURL(bucket string, path string) string {
	return fmt.Sprintf("https://%s.storage.googleapis.com/%s", bucket, path)
}
