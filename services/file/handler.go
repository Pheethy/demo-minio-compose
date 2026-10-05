package file

import "github.com/gin-gonic/gin"

type IFileHandler interface {
	UploadFile(c *gin.Context)
	DownlaodFileMiniO(c *gin.Context)
	DeleteFileMinio(c *gin.Context)
	DeleteFileMiniO(c *gin.Context)
	SignURLExpired(c *gin.Context)
}
