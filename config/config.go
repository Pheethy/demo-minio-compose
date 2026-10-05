package config

import (
	"strconv"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		minio *minioConfig
	}

	IConfig interface {
		MinIO() IMinIOConfig
	}
)

func LoadConfig() IConfig {
	envMap, err := godotenv.Read(".env")
	if err != nil {
		envMap = map[string]string{}
	}

	return &Config{
		minio: &minioConfig{
			endpoint:  envMap["MINIO_ENDPOINT"],
			accessKey: envMap["MINIO_ACCESS_KEY"],
			secretKey: envMap["MINIO_SECRET_KEY"],
			bucket:    envMap["MINIO_BUCKET"],
			useSSL: func() bool {
				useSSL, _ := strconv.ParseBool(envMap["MINIO_USE_SSL"])
				return useSSL
			}(),
		},
	}
}

func (c *Config) MinIO() IMinIOConfig {
	return c.minio
}

type (
	minioConfig struct {
		endpoint  string
		accessKey string
		secretKey string
		bucket    string
		useSSL    bool
	}

	IMinIOConfig interface {
		Endpoint() string
		AccessKey() string
		SecretKey() string
		Bucket() string
		UseSSL() bool
	}
)

func (m *minioConfig) Endpoint() string  { return m.endpoint }
func (m *minioConfig) AccessKey() string { return m.accessKey }
func (m *minioConfig) SecretKey() string { return m.secretKey }
func (m *minioConfig) Bucket() string    { return m.bucket }
func (m *minioConfig) UseSSL() bool      { return m.useSSL }
