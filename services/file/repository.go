package file

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type IFileRepository interface {
	UploadFileMinio(ctx context.Context, filePath string, fileName string, fileContent io.Reader, fileSize int64) (*minio.UploadInfo, error)
}
