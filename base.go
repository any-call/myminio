package myminio

import (
	"io"
	"time"
)

type (
	Client interface {
		PutFile(file io.ReadSeeker, bucketStr string, keyStr string, contentType string, timeout time.Duration) (filename string, err error)
		GetFile(bucketStr string, keyStr string, timeout time.Duration) (filename string, err error)
	}
)
