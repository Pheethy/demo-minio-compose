package usecase

import (
	"context"
	"io"

	"github.com/Pheethy/demo-minio-compose/services/file"
	"github.com/minio/minio-go/v7"
)

type fileUsecase struct {
	fileRepo file.IFileRepository
}

func NewFileUseCase(fileRepo file.IFileRepository) file.IFileUseCase {
	return &fileUsecase{fileRepo: fileRepo}
}

func (u *fileUsecase) UploadFile(ctx context.Context, filePath string, fileName string, fileContent io.Reader, fileSize int64) (*minio.UploadInfo, error) {
	return u.fileRepo.UploadFileMinio(ctx, filePath, fileName, fileContent, fileSize)
}
