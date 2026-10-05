package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Pheethy/demo-minio-compose/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	ctx := context.Background()
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

	bucket := cfg.MinIO().Bucket()
	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("bucket %q exists: %v\n", bucket, exists)
}
