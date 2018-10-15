package publisher

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	yaml "github.com/ghodss/yaml"

	"github.com/HotelsDotCom/helm-chart-publisher/config"
	"github.com/HotelsDotCom/helm-chart-publisher/storage"
	_ "github.com/HotelsDotCom/helm-chart-publisher/storage/s3"
	"github.com/pkg/errors"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/provenance"
	"k8s.io/helm/pkg/repo"
)

type Index struct {
	hash  string
	index *repo.IndexFile
}

// Publisher ...
type Publisher struct {
	sync.RWMutex

	indexes map[string]*Index

	config config.Config
	store  storage.Storage

	repos Repos
}

// New creates a new Publisher instance
func New() (*Publisher, error) {
	storageType, storageConfig := config.GetStorage()
	store, err := storage.New(storageType, storageConfig)
	if err != nil {
		return nil, errors.Wrap(err, "initialize storage failed")
	}

	repos, err := decodeRepos(config.GetRepos())
	if err != nil {
		return nil, errors.Wrap(err, "initialize repositories failed")
	}

	indexes := map[string]*Index{}
	for _, r := range repos {
		indexes[r.Name] = &Index{index: repo.NewIndexFile()}
	}

	return &Publisher{
		indexes: indexes,
		store:   store,
		repos:   repos,
	}, nil
}

// GetIndex ...
func (p *Publisher) GetIndex(repoName string) (*repo.IndexFile, error) {
	repo, err := p.repos.Get(repoName)
	if err != nil {
		return nil, err
	}
	return p.getIndex(repo)
}

// Publish stores the chart in the given repository, updates correspondent index and stores it too.
func (p *Publisher) Publish(repoName string, filename string, chart io.Reader) error {
	// Fetches the repo by name
	repo, err := p.repos.Get(repoName)
	if err != nil {
		return err
	}

	// Send the Chart to the store
	content, err := ioutil.ReadAll(chart)
	if err != nil {
		return err
	}
	if _, err := p.storeFile(repo, filename, content); err != nil {
		return err
	}

	// Updates the index
	if err := p.updateIndex(repo, filename, content); err != nil {
		return err
	}

	return nil
}

func (p *Publisher) storeFile(r *Repo, filename string, content []byte) (*storage.PutResponse, error) {
	resp, err := p.store.Put(r.Bucket, r.Path(filename), content)
	if err != nil {
		return nil, StorageErr{err, fmt.Sprintf("store file %s failed", filename)}
	}
	return resp, nil
}

func (p *Publisher) updateIndex(r *Repo, filename string, chartContent []byte) error {
	// Creating a temporary index with the published chart
	newIndex, err := p.createNewIndex(r, filename, chartContent)
	if err != nil {
		return err
	}

	p.Lock()
	defer p.Unlock()

	// Getting the current index
	currentIndex, err := p.getIndex(r)
	if err != nil {
		return err
	}

	// Merging the current index with the temporary
	currentIndex.Merge(newIndex)
	currentIndex.SortEntries()

	// // Updating the in memory index copy
	// p.indexes[r.Name].index = currentIndex

	// Publishing the updated index to the store
	indexContent, err := yaml.Marshal(currentIndex)
	if err != nil {
		return err
	}
	_, err = p.storeFile(r, "index.yaml", indexContent)
	if err != nil {
		return errors.Wrap(err, "store index.yaml failed")
	}

	// Updating the index hash in memory
	// p.indexes[r.Name].hash = resp.Hash

	return nil
}

// createNewIndex creates temporary index containing a single entrie to be merged with the current index
func (p *Publisher) createNewIndex(r *Repo, filename string, chartContent []byte) (*repo.IndexFile, error) {
	index := repo.NewIndexFile()

	chart, err := chartutil.LoadArchive(bytes.NewBuffer(chartContent))
	if err != nil {
		return nil, HelmErr{err, "Load helm chart failed"}
	}

	hash, err := provenance.Digest(bytes.NewBuffer(chartContent))
	if err != nil {
		return nil, HelmErr{err, "Digest helm chart failed"}
	}
	index.Add(chart.Metadata, filename, p.store.GetURL(r.Bucket, r.Directory, r.Url), hash)

	return index, nil
}

// getIndex gets the index for a given repository. It fetches the index from the store passing the stored in memory hash
// for that index. If the hash hasn't changed, the store should return a NotModifiedErr so we can return the
// current valid index stored in memory.
func (p *Publisher) getIndex(repository *Repo) (*repo.IndexFile, error) {
	index := repo.NewIndexFile()

	resp, err := p.store.Get(repository.Bucket, repository.Path("index.yaml"))
	if err != nil {
		switch err.(type) {
		case storage.PathNotFoundErr:
			return index, nil
		}

		return nil, StorageErr{err, fmt.Sprintf("get index.yaml for %s failed", repository.Name)}
	}

	yaml.Unmarshal(resp.Body, index)

	p.indexes[repository.Name] = &Index{
		index: index,
	}

	return index, nil
}
