package main

import (
	"context"
	"log"

	"github.com/Pheethy/demo-minio-compose/config"
	"github.com/Pheethy/demo-minio-compose/router"
	_file_handler "github.com/Pheethy/demo-minio-compose/services/file/handler"
	_file_repository "github.com/Pheethy/demo-minio-compose/services/file/repository"
	_file_usecase "github.com/Pheethy/demo-minio-compose/services/file/usecase"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	ctx := context.Background()
	_ = ctx
	cfg := config.LoadConfig()

	minioClient, err := minio.New(cfg.MinIO().Endpoint(), &minio.Options{
		Creds: credentials.NewStaticV4(
			cfg.MinIO().AccessKey(),
			cfg.MinIO().SecretKey(),
			"",
		),
		Secure: cfg.MinIO().UseSSL(),
	})
	if err != nil {
		log.Fatal(err)
	}

	/* Init Repository */
	fileRepo := _file_repository.NewFileRepository(minioClient, cfg)

	/* Init Usecase */
	fileUs := _file_usecase.NewFileUseCase(fileRepo)

	/* Init Handler */
	fileHandler := _file_handler.NewFileHandler(fileUs)

	/* New Web Server */
	server := gin.Default()

	/* Register Route */
	route := router.NewRoute(server)
	route.InitFile(fileHandler)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
