package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/adampresley/adamgokit/awsconfig"
	"github.com/adampresley/adamgokit/s3/createbucketoptions"
	"github.com/adampresley/adamgokit/s3/deleteoptions"
	"github.com/adampresley/adamgokit/s3/getoptions"
	"github.com/adampresley/adamgokit/s3/geturloptions"
	"github.com/adampresley/adamgokit/s3/listoptions"
	"github.com/adampresley/adamgokit/s3/putoptions"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Client interface {
	BucketExists(bucket string) (bool, error)
	CreateBucket(bucket string, options ...createbucketoptions.CreateBucketOption) error
	Delete(bucket string, keys []string, options ...deleteoptions.DeleteOption) (DeleteResponse, error)
	Get(bucket, key string, options ...getoptions.GetOption) (GetObjectResponse, error)
	GetUrl(bucket, key string, options ...geturloptions.GetUrlOption) (string, error)
	List(bucket, path string, options ...listoptions.ListOption) (ListResponse, error)
	Put(bucket, key string, body io.Reader, options ...putoptions.PutOption) (PutObjectResponse, error)
	PutStream(bucket, key string, options ...putoptions.PutOption) (PutStreamResponse, error)
	StatObject(bucket, key string) (*ObjectMetadata, error)
}

type Client struct {
	cfg       awsconfig.Configer
	client    *s3.Client
	presigner *s3.PresignClient
}

func NewClient(config awsconfig.Configer) (*Client, error) {
	var (
		err error
	)

	if err = config.Load(); err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	result := &Client{
		cfg: config,
	}

	cfg := config.GetConfigValues().(aws.Config)
	result.client = s3.NewFromConfig(cfg)
	result.presigner = s3.NewPresignClient(result.client)

	return result, nil
}

func (c *Client) BucketExists(bucket string) (bool, error) {
	_, err := c.client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		var nfe *types.NotFound
		if errors.As(err, &nfe) {
			return false, nil
		}

		return false, fmt.Errorf("failed to check if bucket '%s' exists: %w", bucket, err)
	}

	return true, nil
}

func (c *Client) CreateBucket(bucket string, options ...createbucketoptions.CreateBucketOption) error {
	opts := &createbucketoptions.CreateBucketOptions{
		Context: context.Background(),
		Timeout: time.Second * 20,
	}

	for _, opt := range options {
		opt(opts)
	}

	ctx, cancel := context.WithTimeout(opts.Context, opts.Timeout)
	defer cancel()

	createBucketInput := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}

	region := opts.Region
	if region == "" {
		awsRegion := c.cfg.GetConfigValues().(aws.Config).Region
		if awsRegion != "" && awsRegion != "us-east-1" {
			region = types.BucketLocationConstraint(awsRegion)
		}
	}

	if region != "" {
		createBucketInput.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: region,
		}
	}

	_, err := c.client.CreateBucket(ctx, createBucketInput)
	if err != nil {
		return fmt.Errorf("failed to create bucket '%s': %w", bucket, err)
	}

	return nil
}

func (c *Client) Delete(bucket string, keys []string, options ...deleteoptions.DeleteOption) (DeleteResponse, error) {
	var (
		err      error
		output   *s3.DeleteObjectsOutput
		response = DeleteResponse{
			DeletedKeys: []string{},
			Errors:      []ErrorResponse{},
		}
	)

	opts := &deleteoptions.DeleteOptions{
		Context: context.Background(),
		Timeout: time.Second * 20,
	}

	for _, opt := range options {
		opt(opts)
	}

	deleteInput := &s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &types.Delete{
			Objects: []types.ObjectIdentifier{},
		},
	}

	for _, k := range keys {
		deleteInput.Delete.Objects = append(deleteInput.Delete.Objects, types.ObjectIdentifier{
			Key: aws.String(k),
		})
	}

	if output, err = c.client.DeleteObjects(opts.Context, deleteInput); err != nil {
		return response, fmt.Errorf("failed to delete objects from bucket '%s': %w", bucket, err)
	}

	for _, o := range output.Deleted {
		response.DeletedKeys = append(response.DeletedKeys, aws.ToString(o.Key))
	}

	for _, e := range output.Errors {
		response.Errors = append(response.Errors, ErrorResponse{
			Key:     aws.ToString(e.Key),
			Code:    aws.ToString(e.Code),
			Message: aws.ToString(e.Message),
		})
	}

	return response, nil
}

func (c *Client) Get(bucket, key string, options ...getoptions.GetOption) (GetObjectResponse, error) {
	var (
		err    error
		object *v4.PresignedHTTPRequest
		result = GetObjectResponse{}
	)

	opts := &getoptions.GetOptions{
		Context:    context.Background(),
		Expiration: time.Minute * 10,
		Timeout:    time.Second * 5,
	}

	for _, opt := range options {
		opt(opts)
	}

	/*
	 * Get a presigned URL for the object.
	 */
	ctx, cancel := context.WithTimeout(opts.Context, opts.Timeout)
	defer cancel()

	object, err = c.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(opts.Expiration))

	if err != nil {
		return result, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	/*
	 * Make an HTTP GET request to the presigned URL.
	 */
	r, err := http.NewRequest(http.MethodGet, object.URL, nil)

	if err != nil {
		return result, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := http.DefaultClient.Do(r)

	if err != nil {
		return result, fmt.Errorf("failed to make HTTP request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return result, fmt.Errorf("received HTTP status code %d", resp.StatusCode)
	}

	lastModified, _ := http.ParseTime(resp.Header.Get("Last-Modified"))

	result = GetObjectResponse{
		Body:               resp.Body,
		ContentDisposition: object.SignedHeader.Get("Content-Disposition"),
		Size:               resp.ContentLength,
		ContentType:        object.SignedHeader.Get("Content-Type"),
		ETag:               strings.Trim(resp.Header.Get("ETag"), `"`),
		LastModified:       lastModified,
	}

	return result, nil
}

func (c *Client) GetUrl(bucket, key string, options ...geturloptions.GetUrlOption) (string, error) {
	opts := &geturloptions.GetUrlOptions{
		Context:    context.Background(),
		Expiration: time.Hour * 1,
	}

	for _, opt := range options {
		opt(opts)
	}

	return c.getUrl(bucket, key, opts)
}

func (c *Client) getUrl(bucket, key string, options *geturloptions.GetUrlOptions) (string, error) {
	var (
		err error
		req *v4.PresignedHTTPRequest
	)

	ctx, cancel := context.WithTimeout(options.Context, options.Expiration)
	defer cancel()

	req, err = c.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(options.Expiration))

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return req.URL, nil
}

func (c *Client) List(bucket, path string, options ...listoptions.ListOption) (ListResponse, error) {
	var (
		err    error
		result = ListResponse{}
		output *s3.ListObjectsV2Output
	)

	opts := &listoptions.ListOptions{
		Context:           context.Background(),
		ContinuationToken: "",
		GetAll:            false,
		GetUrls:           false,
		GetUrlOptions: &geturloptions.GetUrlOptions{
			Context:    context.Background(),
			Expiration: time.Hour * 1,
		},
		Timeout: time.Second * 10,
	}

	for _, opt := range options {
		opt(opts)
	}

	keepGoing := true

	for keepGoing {
		ctx, cancel := context.WithTimeout(opts.Context, opts.Timeout)
		defer cancel()

		listOptions := &s3.ListObjectsV2Input{
			Bucket:            aws.String(bucket),
			ContinuationToken: aws.String(opts.ContinuationToken),
		}

		if path != "" {
			listOptions.Prefix = aws.String(path)
		}

		if output, err = c.client.ListObjectsV2(ctx, listOptions); err != nil {
			return result, fmt.Errorf("failed to list objects in bucket '%s': %w", bucket, err)
		}

		result.NumObjects = result.NumObjects + int(aws.ToInt32(output.KeyCount))
		result.ContinuationToken = aws.ToString(output.ContinuationToken)

		for _, item := range output.Contents {
			if item.Key == nil {
				continue
			}

			if filepath.Clean(aws.ToString(item.Key)) == filepath.Clean(path) {
				continue
			}

			if opts.Filter != nil && !opts.Filter(item) {
				continue
			}

			newObject := Object{
				ETag:         aws.ToString(item.ETag),
				Key:          aws.ToString(item.Key),
				LastModified: aws.ToTime(item.LastModified),
				Size:         aws.ToInt64(item.Size),
			}

			if item.Owner != nil {
				newObject.OwnerID = aws.ToString(item.Owner.ID)
				newObject.OwnerName = aws.ToString(item.Owner.DisplayName)
			}

			if opts.GetUrls {
				if newObject.Url, err = c.getUrl(bucket, newObject.Key, opts.GetUrlOptions); err != nil {
					return result, fmt.Errorf("failed to generate URL for object '%s': %w", newObject.Key, err)
				}
			}

			result.Objects = append(result.Objects, newObject)
		}

		keepGoing = opts.GetAll && result.ContinuationToken != ""
	}

	return result, nil
}

func (c *Client) Put(bucket, key string, body io.Reader, options ...putoptions.PutOption) (PutObjectResponse, error) {
	var (
		err    error
		resp   *s3.PutObjectOutput
		result = PutObjectResponse{}
	)

	opts := &putoptions.PutOptions{
		Context: context.Background(),
		Timeout: time.Second * 5,
	}

	for _, opt := range options {
		opt(opts)
	}

	ctx, cancel := context.WithTimeout(opts.Context, opts.Timeout)
	defer cancel()

	putObjectInput := &s3.PutObjectInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		Body:     body,
		Metadata: opts.Metadata,
	}

	if opts.ContentType != "" {
		putObjectInput.ContentType = aws.String(opts.ContentType)
	}

	if resp, err = c.client.PutObject(ctx, putObjectInput); err != nil {
		return result, fmt.Errorf("failed to put object in bucket '%s': %w", bucket, err)
	}

	result.Size = aws.ToInt64(resp.Size)
	return result, nil
}

/*
Sets up a stream writer for streaming large files to S3. Use this when you want to upload
large files, or do not know the size of the file in advance.

	stream, _ := client.PutStream(bucket, key, putoptions.WithContentType("application/zip"))
	zipWriter := zip.NewWriter(stream.Writer)

	// Write files to the zip writer...

	_ = zipWriter.Close()
	_ = stream.Writer.Close()

	response, _ := stream.Wait()
*/
func (c *Client) PutStream(bucket, key string, options ...putoptions.PutOption) (PutStreamResponse, error) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	result := PutStreamResponse{}

	opts := &putoptions.PutOptions{
		Context: context.Background(),
		Timeout: 0,
	}

	for _, opt := range options {
		opt(opts)
	}

	ctx = opts.Context
	cancel = func() {}

	if opts.Timeout > 0 {
		ctx, cancel = context.WithTimeout(opts.Context, opts.Timeout)
	} else if ctx == nil {
		ctx = context.Background()
	}

	reader, writer := io.Pipe()
	streamWriter := newPutStreamWriter(writer)

	uploader := manager.NewUploader(c.client)
	resultCh := make(chan putStreamResult, 1)

	go func() {
		defer cancel()

		putObjectInput := &s3.PutObjectInput{
			Bucket:   aws.String(bucket),
			Key:      aws.String(key),
			Body:     reader,
			Metadata: opts.Metadata,
		}

		if opts.ContentType != "" {
			putObjectInput.ContentType = aws.String(opts.ContentType)
		}

		_, err := uploader.Upload(ctx, putObjectInput)
		if err != nil {
			_ = streamWriter.closeWithError(err)
			resultCh <- putStreamResult{
				err: fmt.Errorf("failed to stream object '%s' to bucket '%s': %w", key, bucket, err),
			}
			return
		}

		resultCh <- putStreamResult{
			response: PutObjectResponse{
				Size: streamWriter.bytesWritten(),
			},
		}
	}()

	var (
		once     sync.Once
		waitResp PutObjectResponse
		waitErr  error
	)

	result.Wait = func() (PutObjectResponse, error) {
		once.Do(func() {
			res := <-resultCh
			waitResp = res.response
			waitErr = res.err
		})

		return waitResp, waitErr
	}

	result.Writer = streamWriter

	return result, nil
}

func (c *Client) StatObject(bucket, key string) (*ObjectMetadata, error) {
	var (
		err    error
		output *s3.HeadObjectOutput
	)

	if output, err = c.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		var nfe *types.NotFound
		if errors.As(err, &nfe) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to stat object '%s' in bucket '%s': %w", key, bucket, err)
	}

	result := &ObjectMetadata{
		ETag:         aws.ToString(output.ETag),
		LastModified: aws.ToTime(output.LastModified),
		Size:         aws.ToInt64(output.ContentLength),
		ContentType:  aws.ToString(output.ContentType),
		Metadata:     output.Metadata,
	}

	return result, nil
}

type putStreamResult struct {
	response PutObjectResponse
	err      error
}

type putStreamWriter struct {
	writer  *io.PipeWriter
	written int64
}

func newPutStreamWriter(writer *io.PipeWriter) *putStreamWriter {
	return &putStreamWriter{
		writer: writer,
	}
}

func (w *putStreamWriter) Write(data []byte) (int, error) {
	n, err := w.writer.Write(data)
	if n > 0 {
		atomic.AddInt64(&w.written, int64(n))
	}

	return n, err
}

func (w *putStreamWriter) Close() error {
	return w.writer.Close()
}

func (w *putStreamWriter) closeWithError(err error) error {
	return w.writer.CloseWithError(err)
}

func (w *putStreamWriter) bytesWritten() int64 {
	return atomic.LoadInt64(&w.written)
}
