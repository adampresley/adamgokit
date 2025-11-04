package createbucketoptions

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type CreateBucketOptions struct {
	Context context.Context
	Region  types.BucketLocationConstraint
	Timeout time.Duration
}

type CreateBucketOption func(*CreateBucketOptions)

func WithContext(ctx context.Context) CreateBucketOption {
	return func(o *CreateBucketOptions) {
		o.Context = ctx
	}
}

func WithRegion(region string) CreateBucketOption {
	return func(o *CreateBucketOptions) {
		o.Region = types.BucketLocationConstraint(region)
	}
}

func WithTimeout(d time.Duration) CreateBucketOption {
	return func(o *CreateBucketOptions) {
		o.Timeout = d
	}
}
