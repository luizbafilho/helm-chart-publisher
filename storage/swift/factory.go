package swift

import (
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/luizbafilho/helm-chart-publisher/storage"
	"github.com/mitchellh/mapstructure"
	"github.com/ncw/swift"
)

func init() {
	storage.Register("swift", NewSwiftStorage)
}

// NewSwiftStorage returns a new storage
func NewSwiftStorage(conf map[string]interface{}) (storage.Storage, error) {
	config, err := decodeAndValidateConfig(conf)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConnsPerHost: 2048,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify},
	}

	c := &swift.Connection{
		UserName:       config.Username,
		ApiKey:         config.Password,
		AuthUrl:        config.AuthURL,
		Tenant:         config.Tenant,
		Region:         config.Region,
		EndpointType:   swift.EndpointType(config.EndpointType),
		ConnectTimeout: 60 * time.Second,
		Timeout:        15 * 60 * time.Second,
		Transport:      transport,
	}

	// Authenticate
	err = c.Authenticate()
	if err != nil {
		panic(err)
	}

	return &SwiftStore{
		name:   "swift",
		config: config,
		swift:  c,
	}, nil
}

func decodeAndValidateConfig(c map[string]interface{}) (*Config, error) {
	config := Config{}
	if err := mapstructure.Decode(c, &config); err != nil {
		return nil, err
	}

	mergeEnviromentVariables(&config)

	if config.Username == "" {
		return nil, errors.New("Invalid config. Username is required.")
	}
	if config.Password == "" {
		return nil, errors.New("Invalid config. Password is required.")
	}
	if config.AuthURL == "" {
		return nil, errors.New("Invalid config. AuthURL is required.")
	}
	if config.Tenant == "" {
		return nil, errors.New("Invalid config. Tenant is required.")
	}
	if config.EndpointType == "" {
		return nil, errors.New("Invalid config. EndpointType is required.")
	}

	return &config, nil
}

func mergeEnviromentVariables(config *Config) {
	if username := os.Getenv("SWIFT_USERNAME"); username != "" {
		config.Username = username
	}

	if password := os.Getenv("SWIFT_PASSWORD"); password != "" {
		config.Password = password
	}
}
