package usecase

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/minio/minio-go/v7"
)

type fileRepositoryStub struct {
	ctx         context.Context
	filePath    string
	fileName    string
	fileContent io.Reader
	fileSize    int64
	uploadInfo  *minio.UploadInfo
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
