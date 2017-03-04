package s3

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

type Config struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
}

type S3Store struct {
	name   string
	config *Config
	s3     *s3.S3
}

// Name ...
func (s *S3Store) Name() string {
	return s.name
}

// Put stores the content
func (s *S3Store) Put(path string, content []byte) error {

	params := &s3.PutObjectInput{
		Bucket:      aws.String(s.config.Bucket),
		Key:         aws.String(path),
		Body:        bytes.NewReader(content),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("application/gzip"),
	}

	_, err := s.s3.PutObject(params)
	if err != nil {
		return errors.Wrap(err, "[S3] PutObject failed")
	}

	return nil
}
