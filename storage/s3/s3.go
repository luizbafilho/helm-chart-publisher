package s3

import (
	"bytes"
	"fmt"
	"github.com/HotelsDotCom/helm-chart-publisher/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"io/ioutil"
)

type Config struct {
	Bucket   string
	Region   string
	Protocol string
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

// Create a new aws session for new requests
func newS3Session(region string) *s3.S3 {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	return s3.New(sess)
}

// Get ...
func (s *S3Store) Get(bucket string, path string) (*storage.GetResponse, error) {
	s.s3 = newS3Session(s.config.Region)
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
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
		Body: body,
	}, nil
}

// Put stores the content
func (s *S3Store) Put(bucket string, path string, content []byte) (*storage.PutResponse, error) {
	s.s3 = newS3Session(s.config.Region)
	params := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(path),
		Body:        bytes.NewReader(content),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("application/gzip"),
	}

	_, err := s.s3.PutObject(params)
	if err != nil {
		return nil, errors.Wrap(err, "[S3] PutObject failed")
	}

	return &storage.PutResponse{}, nil
}

// GetURL ...
func (s *S3Store) GetURL(bucket string, path string, url string) string {
	if url != "" {
		return fmt.Sprintf(url)

	} else if s.config.Protocol == "s3" {
		return s.getS3URL(bucket, path)
	}

	return s.getHttpsURL(bucket, path)
}

func (s *S3Store) getHttpsURL(bucket, path string) string {
	domain := fmt.Sprintf("s3-%s.amazonaws.com", s.config.Region)
	if s.config.Region == "us-east-1" {
		domain = "s3.amazonaws.com"
	}
	return fmt.Sprintf("https://%s/%s/%s", domain, bucket, path)
}

func (s *S3Store) getS3URL(bucket, path string) string {
	return fmt.Sprintf("s3://%s/%s", bucket, path)
}
