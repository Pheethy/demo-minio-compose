package repository

import (
	"context"
	"io"
	"mime"
	"path/filepath"

	"github.com/Pheethy/demo-minio-compose/config"
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
	info, err := r.minioClient.PutObject(ctx, r.cfg.MinIO().Bucket(), fileName, fileContent, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return nil, err
	}

	return &info, nil
}
