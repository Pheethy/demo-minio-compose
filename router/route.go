package router

import (
	"github.com/Pheethy/demo-minio-compose/services/file"
	"github.com/gin-gonic/gin"
)

type Route struct {
	c *gin.Engine
}

func NewRoute(c *gin.Engine) *Route {
	return &Route{
		c: c,
	}
}

func (r *Route) InitFile(handler file.IFileHandler) {
	r.c.POST("/file", handler.UploadFile)
}
