package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/luizbafilho/helm-chart-publisher/storage"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func init() {
	storage.Register("s3", NewS3Storage)
}

// NewS3Storage returns a new storage
func NewS3Storage(conf map[string]interface{}) (storage.Storage, error) {
	config, err := decodeAndValidateConfig(conf)
	if err != nil {
		return nil, err
	}

	creds := credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, "")
	_, err = creds.Get()
	if err != nil {
		return nil, errors.Wrap(err, "bad credentials")
	}
	cfg := aws.NewConfig().WithRegion(config.Region).WithCredentials(creds)
	svc := s3.New(session.New(), cfg)

	return &S3Store{
		name:   "s3",
		config: config,
		s3:     svc,
	}, nil
}

func decodeAndValidateConfig(c map[string]interface{}) (*Config, error) {
	config := Config{}
	if err := mapstructure.Decode(c, &config); err != nil {
		return nil, err
	}

	if config.AccessKey == "" {
		return nil, errors.New("Invalid config. Access Key is required.")
	}
	if config.SecretKey == "" {
		return nil, errors.New("Invalid config. Secret Key is required.")
	}
	if config.Bucket == "" {
		return nil, errors.New("Invalid config. Bucket is required.")
	}
	if config.Region == "" {
		return nil, errors.New("Invalid config. Region is required.")
	}

	return &config, nil
}
