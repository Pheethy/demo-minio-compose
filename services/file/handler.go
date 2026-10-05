package file

import "github.com/gin-gonic/gin"

type IFileHandler interface {
	UploadFile(c *gin.Context)
}
