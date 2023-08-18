package myminio

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"time"
)

type client struct {
	*s3.S3
}

func NewClient(endpoint, accessID, accessKey, token string) (Client, error) {
	S3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessID, accessKey, token),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession, err := session.NewSession(S3Config)
	if err != nil {
		return nil, err
	}

	return &client{s3.New(newSession)}, err
}

func (self *client) PutFile(file io.ReadSeeker, bucketStr string, keyStr string, contentType string, timeout time.Duration) (filename string, err error) {
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}

	if cancelFn != nil {
		defer cancelFn()
	}

	var out *s3.PutObjectOutput
	if contentType == "" {
		out, err = self.S3.PutObjectWithContext(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketStr),
			Key:    aws.String(keyStr),
			Body:   file,
		})
	} else {
		out, err = self.S3.PutObjectWithContext(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(bucketStr),
			Key:         aws.String(keyStr),
			ContentType: aws.String(contentType),
			Body:        file,
		})
	}

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			return "", fmt.Errorf("upload canceled due to timeout, %v\n", err)
		}
		return "", err
	}

	return out.String(), nil
}

func (self *client) GetFile(bucketStr string, keyStr string, timeout time.Duration) (filename string, err error) {
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}

	if cancelFn != nil {
		defer cancelFn()
	}

	out, err := self.S3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketStr),
		Key:    aws.String(keyStr),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			return "", fmt.Errorf("upload canceled due to timeout, %v\n", err)
		}
		return "", err
	}

	return out.String(), nil
}
