package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigLoadsMinIOBucket(t *testing.T) {
	tempDir := t.TempDir()
	if err := os.WriteFile(
		filepath.Join(tempDir, ".env"),
		[]byte("MINIO_BUCKET=test-bucket\n"),
		0o600,
	); err != nil {
		t.Fatalf("write test .env: %v", err)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("change to temp directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(workingDir); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	if got := LoadConfig().MinIO().Bucket(); got != "test-bucket" {
		t.Errorf("Bucket() = %q, want %q", got, "test-bucket")
	}
}
