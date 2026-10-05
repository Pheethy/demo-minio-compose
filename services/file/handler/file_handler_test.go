package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Pheethy/demo-minio-compose/models"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type mockMinioObject struct {
	io.ReadCloser
	statInfo minio.ObjectInfo
	statErr  error
}

func (m *mockMinioObject) Stat() (minio.ObjectInfo, error) {
	return m.statInfo, m.statErr
}

type fileUseCaseStub struct {
	uploadInfo *minio.UploadInfo
	uploadErr  error
	downloadObj models.MinioObject
	downloadErr error
	deleteErr   error
	signedURL   string
	signURLErr  error
	receivedKey string
}

func (s *fileUseCaseStub) UploadFile(ctx context.Context, filePath string, fileName string, fileContent io.Reader, fileSize int64) (*minio.UploadInfo, error) {
	return s.uploadInfo, s.uploadErr
}

func (s *fileUseCaseStub) DownloadFileMiniO(ctx context.Context, key string) (models.MinioObject, error) {
	s.receivedKey = key
	return s.downloadObj, s.downloadErr
}

func (s *fileUseCaseStub) DownloadFileMinio(ctx context.Context, key string) (models.MinioObject, error) {
	return s.DownloadFileMiniO(ctx, key)
}

func (s *fileUseCaseStub) DeleteFileMinio(ctx context.Context, key string) error {
	s.receivedKey = key
	return s.deleteErr
}

func (s *fileUseCaseStub) DeleteFileMiniO(ctx context.Context, key string) error {
	return s.DeleteFileMinio(ctx, key)
}

func (s *fileUseCaseStub) SignURLExpired(ctx context.Context, key string) (string, error) {
	s.receivedKey = key
	return s.signedURL, s.signURLErr
}

func TestDownlaodFileMiniO_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockObj := &mockMinioObject{
		ReadCloser: io.NopCloser(strings.NewReader("hello world")),
		statInfo: minio.ObjectInfo{
			Key:         "assets/test.txt",
			ContentType: "text/plain",
		},
	}

	usecaseStub := &fileUseCaseStub{
		downloadObj: mockObj,
	}

	h := NewFileHandler(usecaseStub)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/file?key=assets/test.txt", nil)
	c.Request = req

	h.DownlaodFileMiniO(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if got := w.Header().Get("Content-Disposition"); got != `attachment; filename="assets/test.txt"` {
		t.Fatalf("unexpected Content-Disposition: %s", got)
	}

	if got := w.Header().Get("Content-Type"); got != "text/plain" {
		t.Fatalf("unexpected Content-Type: %s", got)
	}

	if body := w.Body.String(); body != "hello world" {
		t.Fatalf("expected body 'hello world', got '%s'", body)
	}

	if usecaseStub.receivedKey != "assets/test.txt" {
		t.Fatalf("expected key 'assets/test.txt', got '%s'", usecaseStub.receivedKey)
	}
}

func TestDownlaodFileMiniO_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	usecaseStub := &fileUseCaseStub{
		downloadErr: errors.New("file not found"),
	}

	h := NewFileHandler(usecaseStub)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/file?key=missing.txt", nil)
	c.Request = req

	h.DownlaodFileMiniO(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestDownlaodFileMiniO_StatError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockObj := &mockMinioObject{
		ReadCloser: io.NopCloser(strings.NewReader("")),
		statErr:    errors.New("stat failed"),
	}

	usecaseStub := &fileUseCaseStub{
		downloadObj: mockObj,
	}

	h := NewFileHandler(usecaseStub)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/file?key=error.txt", nil)
	c.Request = req

	h.DownlaodFileMiniO(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestDeleteFileMinio_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	usecaseStub := &fileUseCaseStub{}
	h := NewFileHandler(usecaseStub)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodDelete, "/file?key=assets/test.txt", nil)
	c.Request = req

	h.DeleteFileMinio(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if usecaseStub.receivedKey != "assets/test.txt" {
		t.Fatalf("expected key 'assets/test.txt', got '%s'", usecaseStub.receivedKey)
	}
}

func TestDeleteFileMinio_MissingKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	usecaseStub := &fileUseCaseStub{}
	h := NewFileHandler(usecaseStub)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodDelete, "/file", nil)
	c.Request = req

	h.DeleteFileMinio(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestDeleteFileMinio_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	usecaseStub := &fileUseCaseStub{
		deleteErr: errors.New("failed to delete"),
	}
	h := NewFileHandler(usecaseStub)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodDelete, "/file?key=assets/test.txt", nil)
	c.Request = req

	h.DeleteFileMinio(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestSignURLExpired_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	usecaseStub := &fileUseCaseStub{
		signedURL: "http://localhost:9000/bucket/file.png?token=123",
	}
	h := NewFileHandler(usecaseStub)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/file/sign-url?key=assets/test.txt", nil)
	c.Request = req

	h.SignURLExpired(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if usecaseStub.receivedKey != "assets/test.txt" {
		t.Fatalf("expected key 'assets/test.txt', got '%s'", usecaseStub.receivedKey)
	}
}

func TestSignURLExpired_MissingKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	usecaseStub := &fileUseCaseStub{}
	h := NewFileHandler(usecaseStub)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/file/sign-url", nil)
	c.Request = req

	h.SignURLExpired(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestSignURLExpired_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	usecaseStub := &fileUseCaseStub{
		signURLErr: errors.New("failed to sign url"),
	}
	h := NewFileHandler(usecaseStub)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/file/sign-url?key=assets/test.txt", nil)
	c.Request = req

	h.SignURLExpired(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}


