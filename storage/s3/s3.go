package s3

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/luizbafilho/chart-server/storage"
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

func parseError(err error) error {
	if s3Err, ok := err.(awserr.Error); ok {
		switch s3Err.Code() {
		case "NotModified":
			return storage.NotModifiedErr{}
		case "NoSuchKey":
			return storage.PathNotFoundErr{}
		}
	}

	return err
}

// Get ...
func (s *S3Store) Get(bucket string, path string, hash string) (*storage.GetResponse, error) {
	params := &s3.GetObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(path),
		IfNoneMatch: aws.String(hash),
	}

	resp, err := s.s3.GetObject(params)
	if err != nil {
		return nil, parseError(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &storage.GetResponse{
		Hash: *resp.ETag,
		Body: body,
	}, nil
}

// Put stores the content
func (s *S3Store) Put(bucket string, path string, content []byte) (*storage.PutResponse, error) {
	params := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(path),
		Body:        bytes.NewReader(content),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("application/gzip"),
	}

	resp, err := s.s3.PutObject(params)
	if err != nil {
		return nil, errors.Wrap(err, "[S3] PutObject failed")
	}

	return &storage.PutResponse{
		Hash: *resp.ETag,
	}, nil
}

// GetURL ...
func (s *S3Store) GetURL(bucket string, path string) string {
	return fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", s.config.Region, bucket, path)
}
