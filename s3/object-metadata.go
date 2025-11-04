package s3

import "time"

type ObjectMetadata struct {
	ETag         string
	LastModified time.Time
	Size         int64
	ContentType  string
	Metadata     map[string]string
}
