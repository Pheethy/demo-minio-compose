package usecase

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/Pheethy/demo-minio-compose/models"
	"github.com/minio/minio-go/v7"
)

type dummyMinioObject struct {
	io.ReadCloser
	info minio.ObjectInfo
	err  error
}

func (d *dummyMinioObject) Stat() (minio.ObjectInfo, error) {
	return d.info, d.err
}

type fileRepositoryStub struct {
	ctx         context.Context
	filePath    string
	fileName    string
	fileContent io.Reader
	fileSize    int64
	uploadInfo  *minio.UploadInfo
	downloadKey string
	minioObj    models.MinioObject
	signedURL   string
	err         error
}

func (s *fileRepositoryStub) UploadFileMinio(ctx context.Context, filePath string, fileName string, fileContent io.Reader, fileSize int64) (*minio.UploadInfo, error) {
	s.ctx = ctx
	s.filePath = filePath
	s.fileName = fileName
	s.fileContent = fileContent
	s.fileSize = fileSize

	return s.uploadInfo, s.err
}

func (s *fileRepositoryStub) DownloadFileMinio(ctx context.Context, filePath string) (models.MinioObject, error) {
	s.ctx = ctx
	s.downloadKey = filePath

	return s.minioObj, s.err
}

func (s *fileRepositoryStub) DeleteFileMinio(ctx context.Context, filePath string) error {
	s.ctx = ctx
	s.downloadKey = filePath

	return s.err
}

func (s *fileRepositoryStub) SignURLExpired(ctx context.Context, key string) (string, error) {
	s.ctx = ctx
	s.downloadKey = key

	return s.signedURL, s.err
}

func TestFileUseCaseUploadFileDelegatesToRepository(t *testing.T) {
	ctx := context.Background()
	content := strings.NewReader("file content")
	want := &minio.UploadInfo{Key: "documents/report.txt"}
	repo := &fileRepositoryStub{uploadInfo: want}
	useCase := NewFileUseCase(repo)

	got, err := useCase.UploadFile(ctx, "/tmp/report.txt", "documents/report.txt", content, content.Size())
	if err != nil {
		t.Fatalf("UploadFile() error = %v", err)
	}
	if got != want {
		t.Fatalf("UploadFile() = %p, want %p", got, want)
	}
	if repo.ctx != ctx || repo.filePath != "/tmp/report.txt" || repo.fileName != "documents/report.txt" || repo.fileContent != content || repo.fileSize != content.Size() {
		t.Fatal("UploadFile() did not forward all arguments to the repository")
	}
}

func TestFileUseCaseDownloadFileMiniODelegatesToRepository(t *testing.T) {
	ctx := context.Background()
	dummyObj := &dummyMinioObject{
		ReadCloser: io.NopCloser(strings.NewReader("download content")),
		info: minio.ObjectInfo{
			Key:         "assets/photo.png",
			ContentType: "image/png",
		},
	}
	repo := &fileRepositoryStub{minioObj: dummyObj}
	useCase := NewFileUseCase(repo)

	got, err := useCase.DownloadFileMiniO(ctx, "assets/photo.png")
	if err != nil {
		t.Fatalf("DownloadFileMiniO() error = %v", err)
	}
	if got != dummyObj {
		t.Fatalf("DownloadFileMiniO() = %p, want %p", got, dummyObj)
	}
	if repo.ctx != ctx || repo.downloadKey != "assets/photo.png" {
		t.Fatalf("DownloadFileMiniO() did not forward key, got: %s", repo.downloadKey)
	}

	// Test DownloadFileMinio alias
	gotAlias, err := useCase.DownloadFileMinio(ctx, "assets/photo.png")
	if err != nil {
		t.Fatalf("DownloadFileMinio() error = %v", err)
	}
	if gotAlias != dummyObj {
		t.Fatalf("DownloadFileMinio() = %p, want %p", gotAlias, dummyObj)
	}
}

func TestFileUseCaseDeleteFileMinioDelegatesToRepository(t *testing.T) {
	ctx := context.Background()
	repo := &fileRepositoryStub{}
	useCase := NewFileUseCase(repo)

	err := useCase.DeleteFileMinio(ctx, "assets/photo.png")
	if err != nil {
		t.Fatalf("DeleteFileMinio() error = %v", err)
	}
	if repo.ctx != ctx || repo.downloadKey != "assets/photo.png" {
		t.Fatalf("DeleteFileMinio() did not forward key, got: %s", repo.downloadKey)
	}

	// Test DeleteFileMiniO alias
	err = useCase.DeleteFileMiniO(ctx, "assets/photo2.png")
	if err != nil {
		t.Fatalf("DeleteFileMiniO() error = %v", err)
	}
	if repo.downloadKey != "assets/photo2.png" {
		t.Fatalf("DeleteFileMiniO() did not forward key, got: %s", repo.downloadKey)
	}
}

func TestFileUseCaseSignURLExpiredDelegatesToRepository(t *testing.T) {
	ctx := context.Background()
	repo := &fileRepositoryStub{
		signedURL: "http://localhost:9000/bucket/file.png?token=123",
	}
	useCase := NewFileUseCase(repo)

	got, err := useCase.SignURLExpired(ctx, "assets/photo.png")
	if err != nil {
		t.Fatalf("SignURLExpired() error = %v", err)
	}
	if got != repo.signedURL {
		t.Fatalf("SignURLExpired() = %s, want %s", got, repo.signedURL)
	}
	if repo.ctx != ctx || repo.downloadKey != "assets/photo.png" {
		t.Fatalf("SignURLExpired() did not forward key, got: %s", repo.downloadKey)
	}
}

