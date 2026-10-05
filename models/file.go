package models

import (
	"io"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type FileUploadRequest struct {
	FileName    string `json:"file_name" form:"file_name"`
	FilePath    string `json:"file_path" form:"file_path"`
	FileContent *multipart.FileHeader
	FileBytes   []byte
}

type MinioObject interface {
	io.ReadCloser
	Stat() (minio.ObjectInfo, error)
}
