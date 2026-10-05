package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Pheethy/demo-minio-compose/models"
	"github.com/Pheethy/demo-minio-compose/services/file"
	"github.com/gin-gonic/gin"
)

type fileHandler struct {
	fileUsecase file.IFileUseCase
}

func NewFileHandler(fileUsecase file.IFileUseCase) file.IFileHandler {
	return &fileHandler{
		fileUsecase: fileUsecase,
	}
}

func (f *fileHandler) UploadFile(c *gin.Context) {
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

	uploadInfo, err := f.fileUsecase.UploadFile(ctx, fileReq.FilePath, fileReq.FileName, file, fileReq.FileContent.Size)
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

func (f *fileHandler) DownlaodFileMiniO(c *gin.Context) {
	ctx := c.Request.Context()
	key := c.Query("key")

	file, err := f.fileUsecase.DownloadFileMiniO(ctx, key)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, stat.Key))
	c.Header("Content-Type", stat.ContentType)

	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}

func (f *fileHandler) DeleteFileMinio(c *gin.Context) {
	ctx := c.Request.Context()
	key := c.Query("key")
	if key == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "key is required",
		})
		return
	}

	err := f.fileUsecase.DeleteFileMinio(ctx, key)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successful",
	})
}

func (f *fileHandler) DeleteFileMiniO(c *gin.Context) {
	f.DeleteFileMinio(c)
}

func (f *fileHandler) SignURLExpired(c *gin.Context) {
	ctx := c.Request.Context()
	key := c.Query("key")
	if key == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "key is required",
		})
		return
	}

	url, err := f.fileUsecase.SignURLExpired(ctx, key)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successful",
		"url":     url,
	})
}
