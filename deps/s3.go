package create

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"github.com/alterejoe/structs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// GetS3Manager builds an S3Manager based on environment and returns it.
// func CreateS3Manager(bucket string, endpoint string, region string, accessKey string, secretKey string) (*S3Manager, error) {
func CreateS3Manager(auth *structs.S3Auth) (*S3Manager, error) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(auth.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if auth.Endpoint != "" {
			o.BaseEndpoint = aws.String(auth.Endpoint)
			o.UsePathStyle = true
		}
	})

	return &S3Manager{
		client:     client,
		uploader:   s3manager.NewUploader(client),
		downloader: s3manager.NewDownloader(client),
		bucket:     auth.Bucket,
	}, nil
}

type S3Manager struct {
	client     *s3.Client
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader

	bucket string
}

func (m *S3Manager) FileURL(ctx context.Context, key string) string {
	// if os.Getenv("APP_ENV") == "dev" {
	endpoint := strings.TrimSuffix(*m.client.Options().BaseEndpoint, "/")
	return fmt.Sprintf("%s/%s/%s", endpoint, m.bucket, key)
	// }

	// url, err := m.PresignGetURL(ctx, key)
	// if err != nil {
	// 	return ""
	// }
	// return url
}

// UploadFile uploads a local file path → S3 key.
func (m *S3Manager) UploadFile(ctx context.Context, key string, localPath string) error {
	f, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	_, err = m.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(key),
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("upload: %w", err)
	}
	return nil
}
func groupPrefix(groupID string) string {
	return fmt.Sprintf("groups/%s/", groupID)
}

func (m *S3Manager) GenerateKey(ctx context.Context, groupID string, filename string) string {
	return groupPrefix(groupID) + filename
}

func (m *S3Manager) AddFileToGroup(ctx context.Context, groupID, fileName string, r io.Reader) error {
	key := m.GenerateKey(ctx, groupID, fileName)

	return m.UploadReader(ctx, key, r)
}
func (m *S3Manager) ListGroupFiles(ctx context.Context, groupID string) ([]string, error) {
	prefix := groupPrefix(groupID)

	out, err := m.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(m.bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, fmt.Errorf("list group files: %w", err)
	}

	keys := make([]string, 0, len(out.Contents))
	for _, obj := range out.Contents {
		keys = append(keys, aws.ToString(obj.Key))
	}

	return keys, nil
}
func (m *S3Manager) DeleteFileFromGroup(ctx context.Context, groupID, fileName string) error {
	key := m.GenerateKey(ctx, groupID, fileName)

	_, err := m.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("delete group file: %w", err)
	}

	return nil
}

// DownloadFile downloads key → local file path.
func (m *S3Manager) DownloadFile(ctx context.Context, key string, localPath string) error {
	f, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	_, err = m.downloader.Download(ctx, f, &s3.GetObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	return nil
}

// UploadReader accepts any io.Reader (e.g. uploaded file/stream).
func (m *S3Manager) UploadReader(ctx context.Context, key string, r io.Reader) error {
	_, err := m.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(key),
		Body:   r,
	})
	if err != nil {
		return fmt.Errorf("upload reader: %w", err)
	}
	return nil
}
func (s *S3Manager) FileExists(ctx context.Context, groupID, filename string) (bool, error) {
	key := s.GenerateKey(ctx, groupID, filename)

	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err == nil {
		return true, nil
	}

	var nsk *types.NotFound
	if errors.As(err, &nsk) {
		return false, nil
	}

	return false, err
}
