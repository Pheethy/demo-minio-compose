package file

import (
	"context"
	"io"

	"github.com/Pheethy/demo-minio-compose/models"
	"github.com/minio/minio-go/v7"
)

type IFileRepository interface {
	UploadFileMinio(ctx context.Context, filePath string, fileName string, fileContent io.Reader, fileSize int64) (*minio.UploadInfo, error)
	DownloadFileMinio(ctx context.Context, filePath string) (models.MinioObject, error)
	DeleteFileMinio(ctx context.Context, filePath string) error
	SignURLExpired(ctx context.Context, key string) (string, error)
}
