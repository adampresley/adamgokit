package s3

import "io"

// PutStreamResponse contains the writer to stream data to S3 and a Wait function
// that blocks until the upload completes, returning the final PutObjectResponse.
type PutStreamResponse struct {
	Writer io.WriteCloser
	Wait   func() (PutObjectResponse, error)
}
