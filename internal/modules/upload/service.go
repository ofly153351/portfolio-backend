package upload

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"portfolio-backend/internal/config"
)

type Service struct {
	client        *minio.Client
	bucket        string
	publicBaseURL string
}

func NewService(cfg config.Config) (*Service, error) {
	client, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIORootUser, cfg.MinIORootPassword, ""),
		Secure: cfg.MinIOUseSSL,
	})
	if err != nil {
		return nil, err
	}
	s := &Service{
		client:        client,
		bucket:        strings.TrimSpace(cfg.MinIOBucket),
		publicBaseURL: strings.TrimRight(strings.TrimSpace(cfg.MinIOPublicBaseURL), "/"),
	}
	if s.bucket == "" {
		s.bucket = "portfolio"
	}
	if s.publicBaseURL == "" {
		scheme := "http"
		if cfg.MinIOUseSSL {
			scheme = "https"
		}
		s.publicBaseURL = fmt.Sprintf("%s://%s", scheme, cfg.MinIOEndpoint)
	}
	if err := s.ensureBucket(context.Background()); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Service) UploadImage(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (string, error) {
	if s == nil || s.client == nil {
		return "", ErrUploadFailed
	}
	if fileHeader == nil {
		return "", ErrFileRequired
	}
	if err := s.ensureBucket(ctx); err != nil {
		return "", err
	}
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	if !isAllowedImageType(fileHeader) {
		return "", ErrInvalidFileType
	}

	key, err := buildObjectKey(fileHeader.Filename, folder)
	if err != nil {
		return "", err
	}

	contentType := strings.TrimSpace(fileHeader.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	_, err = s.client.PutObject(ctx, s.bucket, key, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", s.publicBaseURL, s.bucket, key), nil
}

func buildObjectKey(filename string, folder string) (string, error) {
	name := strings.TrimSpace(filename)
	if name == "" {
		return "", ErrFileRequired
	}
	ext := strings.ToLower(filepath.Ext(name))
	if ext == "" {
		ext = ".jpg"
	}
	id := uuid.NewString()
	normalizedFolder := strings.Trim(strings.TrimSpace(folder), "/")
	if normalizedFolder == "" {
		normalizedFolder = "projects"
	}
	return fmt.Sprintf("%s/%s_%s%s", normalizedFolder, time.Now().UTC().Format("20060102"), id, ext), nil
}

func isAllowedImageType(fileHeader *multipart.FileHeader) bool {
	if fileHeader == nil {
		return false
	}
	ct := strings.ToLower(strings.TrimSpace(fileHeader.Header.Get("Content-Type")))
	switch ct {
	case "image/jpeg", "image/jpg", "image/png", "image/webp", "image/gif", "image/svg+xml":
		return true
	}
	ext := strings.ToLower(filepath.Ext(strings.TrimSpace(fileHeader.Filename)))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif", ".svg":
		return true
	default:
		return false
	}
}

func (s *Service) ensureBucket(ctx context.Context) error {
	if s == nil || s.client == nil {
		return ErrUploadFailed
	}
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	if err := s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{}); err != nil {
		return err
	}
	policy := fmt.Sprintf(`{
	  "Version":"2012-10-17",
	  "Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]
	}`, s.bucket)
	_ = s.client.SetBucketPolicy(ctx, s.bucket, policy)
	return nil
}
