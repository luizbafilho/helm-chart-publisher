package swift

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/luizbafilho/helm-chart-publisher/storage"
	"github.com/ncw/swift"
)

type Config struct {
	Username           string
	Password           string
	AuthURL            string
	Tenant             string
	TenantID           string
	Domain             string
	DomainID           string
	TenantDomain       string
	TenantDomainID     string
	TrustID            string
	Region             string
	Container          string
	EndpointType       string
	InsecureSkipVerify bool
}

type SwiftStore struct {
	name   string
	config *Config
	swift  *swift.Connection
}

// Name ...
func (s *SwiftStore) Name() string {
	return s.name
}

func parseError(err error) error {
	if sErr, ok := err.(*swift.Error); ok {
		switch sErr.StatusCode {
		case 304:
			return storage.NotModifiedErr{}
		case 404:
			return storage.PathNotFoundErr{}
		}
	}

	return err
}

// Get ...
func (s *SwiftStore) Get(bucket string, path string, hash string) (*storage.GetResponse, error) {
	var b bytes.Buffer
	content := bufio.NewWriter(&b)

	headers, err := s.swift.ObjectGet(bucket, path, content, false, swift.Headers{"If-None-Match": hash})
	if err != nil {
		return nil, parseError(err)
	}

	return &storage.GetResponse{
		Hash: headers["Etag"],
		Body: b.Bytes(),
	}, nil
}

// Put stores the content
func (s *SwiftStore) Put(bucket string, path string, content []byte) (*storage.PutResponse, error) {
	r := bytes.NewReader(content)
	headers, err := s.swift.ObjectPut(bucket, path, r, false, "", "application/gzip", swift.Headers{})
	if err != nil {
		return nil, parseError(err)
	}

	return &storage.PutResponse{
		Hash: headers["Etag"],
	}, nil
}

// GetURL ...
func (s *SwiftStore) GetURL(bucket string, path string) string {
	return fmt.Sprintf("%s/%s/%s", s.swift.Auth.StorageUrl(true), bucket, path)
}
