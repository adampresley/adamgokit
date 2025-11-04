package awsconfig

import (
	"context"
	"fmt"
	"os"

	awsv2config "github.com/aws/aws-sdk-go-v2/config"
)

type Configer interface {
	GetConfigValues() awsv2config.Config
	Load() error
}

type Config struct {
	Endpoint        string
	Region          string
	AccessKeyID     string
	SecretAccessKey string

	cfg awsv2config.Config
}

func (c *Config) GetConfigValues() awsv2config.Config {
	return c.cfg
}

func (c *Config) Load() error {
	var (
		err error
	)

	// Set os ENV vars because the AWS SDK library is dumb
	os.Setenv("AWS_ACCESS_KEY_ID", c.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", c.SecretAccessKey)
	os.Setenv("AWS_REGION", c.Region)
	os.Setenv("AWS_ENDPOINT_URL", c.Endpoint)

	if c.cfg, err = awsv2config.LoadDefaultConfig(context.TODO()); err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	return nil
}
