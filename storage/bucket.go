package storage

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Bucket struct {
	session *session.Session
	service *s3.S3
}

func (bucket *Bucket) CreateSession(id, secret string) error {
	var err error
	credentials := credentials.NewStaticCredentials(id, secret, "")

	bucket.session, err = session.NewSession(&aws.Config{
		Credentials: credentials,
		Region:      aws.String("us-east-1"),
	})

	if err != nil {
		return err
	}

	bucket.service = s3.New(bucket.session)

	return nil
}

func (bucket *Bucket) Store(key string, contents string) error {
	params := &s3.PutObjectInput{
		Bucket: aws.String("thegreyjoy-historical-data"),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(contents)),
	}

	_, err := bucket.service.PutObject(params)

	if err != nil {
		return err
	}

	return nil
}
