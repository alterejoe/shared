package env

import "os"

type S3ENV struct {
}

// GetBucket() string
// GetEndpoint() string
// GetRegion() string
// GetAccessKey() string
// GetSecretKey() string
func (auth *S3ENV) GetBucket() string {
	return os.Getenv("AWS_BUCKET")
}

func (auth *S3ENV) GetEndpoint() string {
	return os.Getenv("AWS_ENDPOINT")
}

func (auth *S3ENV) GetRegion() string {
	return os.Getenv("AWS_REGION")
}

func (auth *S3ENV) GetAccessKey() string {
	return os.Getenv("AWS_ACCESS_KEY_ID")
}

func (auth *S3ENV) GetSecretKey() string {
	return os.Getenv("AWS_SECRET_ACCESS_KEY")
}
