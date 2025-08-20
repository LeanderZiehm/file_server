package main

import (
    "crypto/sha1"
    "encoding/hex"
    "fmt"
    "net/http"
    "os"
    "strings"
    "time"
)

// newID generates a unique ID from filename + timestamp
func newID(filename string) string {
    data := fmt.Sprintf("%s-%d", filename, time.Now().UnixNano())
    sum := sha1.Sum([]byte(data))
    return hex.EncodeToString(sum[:])
}

// requireAPIKey middleware
func requireAPIKey(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        expected := os.Getenv("API_KEY")
        auth := r.Header.Get("Authorization")
        if !strings.HasPrefix(auth, "Bearer ") || strings.TrimPrefix(auth, "Bearer ") != expected {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next(w, r)
    }
}
