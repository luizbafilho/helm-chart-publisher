package s3

import (
	"github.com/aws/aws-sdk-go/aws"
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

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
	}))
	svc := s3.New(sess)

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

	if config.Bucket == "" {
		return nil, errors.New("Invalid config. Bucket is required.")
	}
	if config.Region == "" {
		return nil, errors.New("Invalid config. Region is required.")
	}

	return &config, nil
}
