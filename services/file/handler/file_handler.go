package handler

import (
	"net/http"

	"github.com/Pheethy/demo-minio-compose/models"
	"github.com/Pheethy/demo-minio-compose/services/file"
	"github.com/gin-gonic/gin"
)

type fileHandler struct {
	fileUs file.IFileUseCase
}

func NewFileHandler(fileUs file.IFileUseCase) file.IFileHandler {
	return fileHandler{
		fileUs: fileUs,
	}
}

func (h fileHandler) UploadFile(c *gin.Context) {
	ctx := c.Request.Context()
	fileReq := new(models.FileUploadRequest)
	form, err := c.MultipartForm()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if val, ok := form.Value["file_path"]; ok {
		fileReq.FilePath = val[0]
	}

	if val, ok := form.File["file"]; ok {
		if len(val) > 0 {
			fileReq.FileContent = val[0]
		}
	}
	fileReq.FileName = fileReq.FileContent.Filename

	file, err := fileReq.FileContent.Open()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer file.Close()

	uploadInfo, err := h.fileUs.UploadFile(ctx, fileReq.FilePath, fileReq.FileName, file, fileReq.FileContent.Size)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp := map[string]interface{}{
		"message":     "successful",
		"upload_info": uploadInfo,
	}

	c.JSON(http.StatusOK, resp)
}
