package file

import (
	"context"
	"io"

	"github.com/Pheethy/demo-minio-compose/models"
	"github.com/minio/minio-go/v7"
)

type IFileUseCase interface {
	UploadFile(ctx context.Context, filePath string, fileName string, fileContent io.Reader, fileSize int64) (*minio.UploadInfo, error)
	DownloadFileMiniO(ctx context.Context, key string) (models.MinioObject, error)
	DownloadFileMinio(ctx context.Context, key string) (models.MinioObject, error)
	DeleteFileMinio(ctx context.Context, key string) error
	DeleteFileMiniO(ctx context.Context, key string) error
	SignURLExpired(ctx context.Context, key string) (string, error)
}
