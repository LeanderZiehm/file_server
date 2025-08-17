package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// Helper: create a file upload request from a real file
func newFileUploadRequestFromFile(uri, paramName, filePath string) (*http.Request, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	io.Copy(part, file)
	writer.Close()

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer testkey")
	return req, nil
}

func TestUploadPDF(t *testing.T) {
	os.Setenv("API_KEY", "testkey")
	os.MkdirAll(dataDir, 0755)
	store.Files = make(map[string]FileMeta)

	req, err := newFileUploadRequestFromFile("/upload", "file", "testdata/sample.pdf")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	requireAPIKey(uploadHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}
}

func TestUploadImage(t *testing.T) {
	os.Setenv("API_KEY", "testkey")
	os.MkdirAll(dataDir, 0755)
	store.Files = make(map[string]FileMeta)

	req, err := newFileUploadRequestFromFile("/upload", "file", "testdata/sample.jpg")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	requireAPIKey(uploadHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}
}
