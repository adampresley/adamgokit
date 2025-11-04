package s3

import (
	"io"
	"time"
)

type GetObjectResponse struct {
	Body               io.ReadCloser
	ContentDisposition string
	Size               int64
	ContentType        string
	ETag               string
	LastModified       time.Time
}
