package storage

import "github.com/labstack/gommon/log"

type storageFactory func(conf map[string]interface{}) (Storage, error)

var storageFactories = make(map[string]storageFactory)

// Register ..
func Register(name string, factory storageFactory) {
	if factory == nil {
		log.Panicf("Datastore factory %s does not exist.", name)
	}

	_, registered := storageFactories[name]
	if registered {
		log.Errorf("Datastore factory %s already registered. Ignoring.", name)
		return
	}

	storageFactories[name] = factory
}
