package storage

import (
	"context"
	"fmt"
	"micro-warehouse/warehouse-service/configs"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOInterface interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*UploadResult, error)
	DeleteFile(ctx context.Context, filePath string) error
}

// UploadResult represents the stored file metadata returned after upload.
type UploadResult struct {
	URL      string `json:"url"`
	Path     string `json:"path"`
	Filename string `json:"filename"`
}

type MinIOStorage struct {
	client *minio.Client
	cfg    configs.Config
}

// UploadFile implements MinIOInterface.
func (m *MinIOStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*UploadResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%s_%d%s", strings.TrimSuffix(file.Filename, ext), timestamp, ext)

	// Create file path
	filePath := fmt.Sprintf("%s/%s", folder, filename)

	// Detect content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".webp":
			contentType = "image/webp"
		case ".svg":
			contentType = "image/svg+xml"
		case ".gif":
			contentType = "image/gif"
		case ".pdf":
			contentType = "application/pdf"
		default:
			contentType = "application/octet-stream"
		}
	}

	// Ensure bucket exists
	exists, err := m.client.BucketExists(ctx, m.cfg.MinIO.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		err = m.client.MakeBucket(ctx, m.cfg.MinIO.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		// Set bucket policy to public read
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, m.cfg.MinIO.Bucket)
		err = m.client.SetBucketPolicy(ctx, m.cfg.MinIO.Bucket, policy)
		if err != nil {
			return nil, fmt.Errorf("failed to set bucket policy: %w", err)
		}
	}

	// Upload file
	_, err = m.client.PutObject(ctx, m.cfg.MinIO.Bucket, filePath, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to minio: %w", err)
	}

	// Generate public URL
	protocol := "http"
	if m.cfg.MinIO.UseSSL {
		protocol = "https"
	}
	publicURL := fmt.Sprintf("%s://%s/%s/%s", protocol, m.cfg.MinIO.Endpoint, m.cfg.MinIO.Bucket, filePath)

	return &UploadResult{
		URL:      publicURL,
		Path:     filePath,
		Filename: filename,
	}, nil
}

// DeleteFile implements MinIOInterface.
func (m *MinIOStorage) DeleteFile(ctx context.Context, filePath string) error {
	err := m.client.RemoveObject(ctx, m.cfg.MinIO.Bucket, filePath, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from minio: %w", err)
	}
	return nil
}

func NewMinIOStorage(cfg configs.Config) (MinIOInterface, error) {
	client, err := minio.New(cfg.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &MinIOStorage{
		client: client,
		cfg:    cfg,
	}, nil
}
