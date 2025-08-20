package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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

// Reset the store and create test data folder
func setupTestEnv(t *testing.T) {
	os.Setenv("API_KEY", "testkey")
	os.MkdirAll(dataDir, 0755)
	store.Files = make(map[string]FileMeta)
}

// Upload a file and return its ID
func uploadFile(t *testing.T, path string) string {
	req, err := newFileUploadRequestFromFile("/upload", "file", path)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	requireAPIKey(uploadHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}

	// response is just the file ID
	return strings.TrimSpace(rr.Body.String())
}

func TestFullFlow(t *testing.T) {
	setupTestEnv(t)

	// 1. Upload PDF
	pdfID := uploadFile(t, "testdata/sample.pdf")
	if pdfID == "" {
		t.Fatal("pdf upload returned empty ID")
	}

	// 2. Upload Image
	imgID := uploadFile(t, "testdata/sample.jpg")
	if imgID == "" {
		t.Fatal("image upload returned empty ID")
	}

	// 3. List files
	listReq := httptest.NewRequest("GET", "/files", nil)
	listReq.Header.Set("Authorization", "Bearer testkey")
	listRR := httptest.NewRecorder()
	requireAPIKey(listHandler).ServeHTTP(listRR, listReq)

	body := listRR.Body.String()
	if !strings.Contains(body, "sample.pdf") || !strings.Contains(body, "sample.jpg") {
		t.Fatal("uploaded files not found in list")
	}

	// 4. Download files
	for _, id := range []string{pdfID, imgID} {
		downloadReq := httptest.NewRequest("GET", "/files/"+id, nil)
		downloadRR := httptest.NewRecorder()
		downloadHandler(downloadRR, downloadReq)

		if downloadRR.Code != http.StatusOK {
			t.Fatalf("download failed for %s with status %d", id, downloadRR.Code)
		}
		if downloadRR.Body.Len() == 0 {
			t.Fatalf("download returned empty content for %s", id)
		}
	}

	// 5. Delete files
	for _, id := range []string{pdfID, imgID} {
		delReq := httptest.NewRequest("DELETE", "/delete/"+id, nil)
		delReq.Header.Set("Authorization", "Bearer testkey")
		delRR := httptest.NewRecorder()
		requireAPIKey(deleteHandler).ServeHTTP(delRR, delReq)

		if delRR.Code != http.StatusOK {
			t.Fatalf("delete failed for %s with status %d", id, delRR.Code)
		}
	}

	// 6. Confirm deletion
	listRR2 := httptest.NewRecorder()
	requireAPIKey(listHandler).ServeHTTP(listRR2, listReq)
	if strings.Contains(listRR2.Body.String(), "sample.pdf") || strings.Contains(listRR2.Body.String(), "sample.jpg") {
		t.Fatal("deleted files still appear in list")
	}

	// 7. Unauthorized access check
	unauthReq := httptest.NewRequest("GET", "/files", nil)
	unauthRR := httptest.NewRecorder()
	requireAPIKey(listHandler).ServeHTTP(unauthRR, unauthReq)
	if unauthRR.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for unauthorized access, got %d", unauthRR.Code)
	}
}
