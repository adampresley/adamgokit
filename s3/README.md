## Amazon AWS S3

This wrapper provides methods for working with S3. Amazon's library is "obtuse" at best. This makes my head hurt less. Here is a basic example of usage.

```go
var (
	err    error
	object s3.GetObjectResponse
)

awsConfig := &awsconfig.Config{
	Endpoint:        config.AwsEndpointUrl,
	Region:          config.AwsRegion,
	AccessKeyID:     config.AwsAccessKeyId,
	SecretAccessKey: config.AwsSecretAccessKey,
}

if err = awsConfig.Load(); err != nil {
	slog.Error("failed to load AWS config. trying again", "error", err)
	return err
}

if err != nil {
	panic(err)
}

s3Client, err := s3.NewClient(awsConfig)

if err != nil {
	panic(err)
}

object, err = c.s3Client.Get(
	"bucket",
	"key",
)
```

