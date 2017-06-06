package swift

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"

	"github.com/luizbafilho/helm-chart-publisher/storage"
	"github.com/ncw/swift"
	"github.com/pkg/errors"
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
func (s *SwiftStore) Get(bucket string, path string) (*storage.GetResponse, error) {
	var b bytes.Buffer
	content := bufio.NewWriter(&b)

	_, err := s.swift.ObjectGet(bucket, path, content, false, swift.Headers{})
	if err != nil {
		return nil, parseError(err)
	}

	return &storage.GetResponse{
		Body: b.Bytes(),
	}, nil
}

// Put stores the content
func (s *SwiftStore) Put(bucket string, path string, content []byte) (*storage.PutResponse, error) {
	// This authentication is done because there is some weird bug when authentication retry
	// has to happen, making the swift package send a request with zero byte body.
	// To avoid that, I'm authenticating before the PUT call, to make sure no retry will be need.
	if err := s.swift.Authenticate(); err != nil {
		return nil, parseError(errors.Wrap(err, "authentication failed"))
	}

	r := bytes.NewReader(content)
	h := swift.Headers{"Content-Length": strconv.Itoa(len(content))}
	_, err := s.swift.ObjectPut(bucket, path, r, true, "", "", h)
	if err != nil {
		return nil, parseError(err)
	}

	return &storage.PutResponse{}, nil
}

// GetURL ...
func (s *SwiftStore) GetURL(bucket string, path string) string {
	return fmt.Sprintf("%s/%s/%s", s.swift.Auth.StorageUrl(true), bucket, path)
}
