package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockFileHandler struct {
	uploadCalled   bool
	downloadCalled bool
	deleteCalled   bool
	signURLCalled  bool
}

func (m *mockFileHandler) UploadFile(c *gin.Context) {
	m.uploadCalled = true
	c.Status(http.StatusOK)
}

func (m *mockFileHandler) DownlaodFileMiniO(c *gin.Context) {
	m.downloadCalled = true
	c.Status(http.StatusOK)
}

func (m *mockFileHandler) DownloadFileMiniO(c *gin.Context) {
	m.DownlaodFileMiniO(c)
}

func (m *mockFileHandler) DeleteFileMinio(c *gin.Context) {
	m.deleteCalled = true
	c.Status(http.StatusOK)
}

func (m *mockFileHandler) DeleteFileMiniO(c *gin.Context) {
	m.DeleteFileMinio(c)
}

func (m *mockFileHandler) SignURLExpired(c *gin.Context) {
	m.signURLCalled = true
	c.Status(http.StatusOK)
}

func TestInitFileRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	r := NewRoute(engine)
	handler := &mockFileHandler{}

	r.InitFile(handler)

	// Test GET /file
	wGet := httptest.NewRecorder()
	reqGet, _ := http.NewRequest(http.MethodGet, "/file?key=assets/test.jpg", nil)
	engine.ServeHTTP(wGet, reqGet)

	if !handler.downloadCalled {
		t.Error("expected DownlaodFileMiniO to be called for GET /file")
	}

	// Test POST /file
	wPost := httptest.NewRecorder()
	reqPost, _ := http.NewRequest(http.MethodPost, "/file", nil)
	engine.ServeHTTP(wPost, reqPost)

	if !handler.uploadCalled {
		t.Error("expected UploadFile to be called for POST /file")
	}

	// Test DELETE /file
	wDelete := httptest.NewRecorder()
	reqDelete, _ := http.NewRequest(http.MethodDelete, "/file?key=assets/test.jpg", nil)
	engine.ServeHTTP(wDelete, reqDelete)

	if !handler.deleteCalled {
		t.Error("expected DeleteFileMinio to be called for DELETE /file")
	}

	// Test GET /file/sign-url
	wSignURL := httptest.NewRecorder()
	reqSignURL, _ := http.NewRequest(http.MethodGet, "/file/sign-url?key=assets/test.jpg", nil)
	engine.ServeHTTP(wSignURL, reqSignURL)

	if !handler.signURLCalled {
		t.Error("expected SignURLExpired to be called for GET /file/sign-url")
	}

	// Test GET /file/sign-url-expired
	handler.signURLCalled = false
	wSignURLExpired := httptest.NewRecorder()
	reqSignURLExpired, _ := http.NewRequest(http.MethodGet, "/file/sign-url-expired?key=assets/test.jpg", nil)
	engine.ServeHTTP(wSignURLExpired, reqSignURLExpired)

	if !handler.signURLCalled {
		t.Error("expected SignURLExpired to be called for GET /file/sign-url-expired")
	}
}

