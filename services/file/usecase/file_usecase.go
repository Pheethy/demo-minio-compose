package usecase

import (
	"context"
	"io"

	"github.com/Pheethy/demo-minio-compose/models"
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

func (u *fileUsecase) DownloadFileMiniO(ctx context.Context, key string) (models.MinioObject, error) {
	return u.fileRepo.DownloadFileMinio(ctx, key)
}

func (u *fileUsecase) DownloadFileMinio(ctx context.Context, key string) (models.MinioObject, error) {
	return u.DownloadFileMiniO(ctx, key)
}

func (u *fileUsecase) DeleteFileMinio(ctx context.Context, key string) error {
	return u.fileRepo.DeleteFileMinio(ctx, key)
}

func (u *fileUsecase) DeleteFileMiniO(ctx context.Context, key string) error {
	return u.DeleteFileMinio(ctx, key)
}

func (u *fileUsecase) SignURLExpired(ctx context.Context, key string) (string, error) {
	return u.fileRepo.SignURLExpired(ctx, key)
}
