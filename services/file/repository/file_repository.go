package repository

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"time"

	"github.com/Pheethy/demo-minio-compose/config"
	"github.com/Pheethy/demo-minio-compose/models"
	"github.com/Pheethy/demo-minio-compose/services/file"
	"github.com/minio/minio-go/v7"
)

type fileRepository struct {
	minioClient *minio.Client
	cfg         config.IConfig
}

func NewFileRepository(minioClient *minio.Client, cfg config.IConfig) file.IFileRepository {
	return &fileRepository{
		minioClient: minioClient,
		cfg:         cfg,
	}
}

func (r *fileRepository) UploadFileMinio(ctx context.Context, filePath string, fileName string, fileContent io.Reader, fileSize int64) (*minio.UploadInfo, error) {
	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	if contentType == "" {
		contentType = "application/octet-stream" // ค่า default หากหา Content-Type ไม่ได้
	}
	info, err := r.minioClient.PutObject(ctx, r.cfg.MinIO().Bucket(), fmt.Sprintf("%s/%s", filePath, fileName), fileContent, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (r *fileRepository) DownloadFileMinio(ctx context.Context, filePath string) (models.MinioObject, error) {
	object, err := r.minioClient.GetObject(ctx, r.cfg.MinIO().Bucket(), filePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (r *fileRepository) DeleteFileMinio(ctx context.Context, filePath string) error {
	return r.minioClient.RemoveObject(ctx, r.cfg.MinIO().Bucket(), filePath, minio.RemoveObjectOptions{})
}

func (f *fileRepository) SignURLExpired(ctx context.Context, key string) (string, error) {
	expired := time.Second * time.Duration(30)
	presignedURL, err := f.minioClient.PresignedGetObject(ctx, f.cfg.MinIO().Bucket(), key, expired, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
