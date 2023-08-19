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
	endpoint  string
	accessKey string
	accessID  string
}

func NewClient(endpoint, accessID, accessKey string) (Client, error) {
	S3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessID, accessKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String("us-east-1"), //us-east-1
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession, err := session.NewSession(S3Config)
	if err != nil {
		return nil, err
	}

	return &client{
		S3:        s3.New(newSession),
		endpoint:  endpoint,
		accessKey: accessKey,
		accessID:  accessID,
	}, err
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

	if contentType == "" {
		_, err = self.S3.PutObjectWithContext(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketStr),
			Key:    aws.String(keyStr),
			Body:   file,
		})
	} else {
		_, err = self.S3.PutObjectWithContext(ctx, &s3.PutObjectInput{
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

	return fmt.Sprintf("%s/%s/%s", self.endpoint, bucketStr, keyStr), nil
}

func (self *client) GetFile(bucketStr string, keyStr string, timeout time.Duration) (ret []byte, err error) {
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
			return nil, fmt.Errorf("upload canceled due to timeout, %v\n", err)
		}
		return nil, err
	}

	defer out.Body.Close()

	b, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (self *client) ListFile(bucketStr string, timeout time.Duration) (ret []string, err error) {
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}

	if cancelFn != nil {
		defer cancelFn()
	}

	ret = []string{}
	err = self.S3.ListObjectsPagesWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(bucketStr),
	}, func(p *s3.ListObjectsOutput, b bool) bool {
		for _, o := range p.Contents {
			ret = append(ret, aws.StringValue(o.Key))
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	return ret, nil
}
