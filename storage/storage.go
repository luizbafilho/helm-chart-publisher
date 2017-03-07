package storage

import (
	"errors"
	"fmt"
)

// GetResponse ...
type GetResponse struct {
	Hash string
	Body []byte
}

// PutResponse ...
type PutResponse struct {
	Hash string
}

// Storage ...
type Storage interface {
	Name() string

	Put(bucket string, path string, content []byte) (*PutResponse, error)
	Get(bucket string, path string, hash string) (*GetResponse, error)

	GetURL(bucket string, directory string) string
}

// PathNotFoundErr is returned when operating on a nonexistent path.
type PathNotFoundErr struct {
	Path string
}

func (err PathNotFoundErr) Error() string {
	return fmt.Sprintf("Path not found: %s", err.Path)
}

// NotModifiedErr  is returned when operating on a nonexistent path.
type NotModifiedErr struct {
	Path string
}

func (err NotModifiedErr) Error() string {
	return fmt.Sprintf("Path not modified: %s", err.Path)
}

// New ...
func New(name string, conf map[string]interface{}) (Storage, error) {
	factory, ok := storageFactories[name]
	if !ok {
		return nil, errors.New("Invalid Datastore name.")
	}

	// Run the factory with the configuration.
	return factory(conf)
}
