package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// --- Handlers ---

func rootHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprintln(w, "File Server API")
    fmt.Fprintln(w, "")
    fmt.Fprintln(w, "Endpoints:")
    fmt.Fprintln(w, "  POST   /upload        (auth required) - Upload an image or PDF")
    fmt.Fprintln(w, "  GET    /files/{id}    (public)        - Download/View a file")
    fmt.Fprintln(w, "  GET    /files         (auth required) - List uploaded files")
    fmt.Fprintln(w, "  DELETE /delete/{id}   (auth required) - Delete a file")
    fmt.Fprintln(w, "")
    fmt.Fprintln(w, "Authentication:")
    fmt.Fprintln(w, "  Use Authorization: Bearer <API_KEY>")
}



// Upload a file (image/pdf only)
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(100 << 20); err != nil { // 100 MB
		http.Error(w, "file too big or invalid form", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Peek at file header for content-type
	buff := make([]byte, 512)
	n, _ := file.Read(buff)
	contentType := http.DetectContentType(buff[:n])
	if !(strings.HasPrefix(contentType, "image/") || contentType == "application/pdf") {
		http.Error(w, "only images and pdf allowed", http.StatusBadRequest)
		return
	}
	// rewind so we can save file properly
	file.Seek(0, io.SeekStart)

	// Generate unique ID from filename + timestamp
	id := newID(header.Filename)

	// Ensure extension is kept
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		if strings.HasPrefix(contentType, "image/") {
			ext = ".jpg" // fallback
		} else if contentType == "application/pdf" {
			ext = ".pdf"
		}
	}
	savePath := filepath.Join(dataDir, id+ext)

	// Save file to disk
	out, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "cannot save file", http.StatusInternalServerError)
		return
	}
	size, err := io.Copy(out, file)
	out.Close()
	if err != nil {
		http.Error(w, "failed writing file", http.StatusInternalServerError)
		return
	}

	meta := FileMeta{
		ID:   id,
		Name: header.Filename,
		Size: size,
		Type: contentType,
		URL:  fmt.Sprintf("/files/%s", id),
		Path: savePath,
	}

	// Save metadata
	store.Lock()
	store.Files[id] = meta
	saveIndex()
	store.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meta)
}

// List all uploaded files
func listHandler(w http.ResponseWriter, r *http.Request) {
	store.Lock()
	defer store.Unlock()

	var list []FileMeta
	for _, f := range store.Files {
		list = append(list, f)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// Delete a file by ID
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/delete/")

	store.Lock()
	defer store.Unlock()

	meta, ok := store.Files[id]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// Remove file from disk
	os.Remove(meta.Path)
	// Remove metadata
	delete(store.Files, id)
	saveIndex()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// Download/View a file (no auth required)
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/files/")

	store.Lock()
	meta, ok := store.Files[id]
	store.Unlock()

	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, meta.Path)
}

// Update a file's name
func updateHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/update/")

	var reqBody struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if reqBody.Name == "" {
		http.Error(w, "name cannot be empty", http.StatusBadRequest)
		return
	}

	store.Lock()
	defer store.Unlock()

	meta, ok := store.Files[id]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// Keep old path for renaming
	oldPath := meta.Path
	newPath := filepath.Join(dataDir, id+filepath.Ext(reqBody.Name))

	// Rename file on disk
	if err := os.Rename(oldPath, newPath); err != nil {
		http.Error(w, "failed to rename file", http.StatusInternalServerError)
		return
	}

	// Update metadata
	meta.Name = reqBody.Name
	meta.Path = newPath
	store.Files[id] = meta
	saveIndex()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meta)
}
