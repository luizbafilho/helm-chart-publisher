package storage

import "errors"

type Storage interface {
	Name() string

	Put(path string, content []byte) error
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
