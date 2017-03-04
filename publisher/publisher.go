package publisher

import (
	"bytes"
	"io"
	"io/ioutil"
	"sync"

	"github.com/luizbafilho/chart-server/config"
	"github.com/luizbafilho/chart-server/storage"
	"github.com/pkg/errors"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/provenance"
	"k8s.io/helm/pkg/repo"
)

// Publisher ...
type Publisher struct {
	sync.Mutex

	indexHash string
	indexFile *repo.IndexFile

	config config.Config
	store  storage.Storage
}

// New ..
func New() (*Publisher, error) {
	storageType, config := config.GetStorageAndConfig()
	store, err := storage.New(storageType, config)
	if err != nil {
		return nil, errors.Wrap(err, "initialize storage failed")
	}

	return &Publisher{
		indexFile: repo.NewIndexFile(),
		store:     store,
		config:    config,
	}, nil
}

func (p *Publisher) GetIndex() *repo.IndexFile {
	return p.indexFile
}

// Publish ...
func (p *Publisher) Publish(path string, chart io.Reader) error {
	content, err := ioutil.ReadAll(chart)
	if err != nil {
		return err
	}

	if err := p.storeChart(path, content); err != nil {
		return err
	}

	if err := p.updateIndex(path, content); err != nil {
		return err
	}

	return nil
}

func (p *Publisher) storeChart(path string, content []byte) error {
	if err := p.store.Put(path, content); err != nil {
		return err
	}

	return nil
}

func (p *Publisher) updateIndex(path string, content []byte) error {
	index := repo.NewIndexFile()

	chart, err := chartutil.LoadArchive(bytes.NewBuffer(content))
	if err != nil {
		return errors.Wrap(err, "Load helm chart failed")
	}

	hash, err := provenance.Digest(bytes.NewBuffer(content))
	if err != nil {
		return errors.Wrap(err, "Digest helm chart failed")
	}

	index.Add(chart.Metadata, path, "http://charts.videos.globoi.com", hash)

	p.Lock()
	p.indexFile.Merge(index)
	p.indexFile.SortEntries()
	p.Unlock()

	return nil
}

func (p *Publisher) getCurrentIndex() (*repo.IndexFile, error) {
	return p.indexFile, nil
}
