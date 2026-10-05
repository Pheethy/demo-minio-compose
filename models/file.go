package models

import "mime/multipart"

type FileUploadRequest struct {
	FileName    string `json:"file_name" form:"file_name"`
	FilePath    string `json:"file_path" form:"file_path"`
	FileContent *multipart.FileHeader
	FileBytes   []byte
}
