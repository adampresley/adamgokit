package s3

import "time"

type Object struct {
	ETag         string
	Key          string
	LastModified time.Time
	OwnerID      string
	OwnerName    string
	Size         int64
	Url          string
}
