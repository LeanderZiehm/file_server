package main

import (
    "log"
    "net/http"
    "os"
)

func main() {
    // Load .env file if present
    _ = loadDotEnv(".env")

    os.MkdirAll(dataDir, 0755)

    apiKey := os.Getenv("API_KEY")
    if apiKey == "" {
        log.Fatal("API_KEY must be set in environment or .env file")
    }

    loadIndex()

    mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler) // add this line
    mux.HandleFunc("/upload", requireAPIKey(uploadHandler))
    mux.HandleFunc("/files", requireAPIKey(listHandler))
    mux.HandleFunc("/delete/", requireAPIKey(deleteHandler))
    mux.HandleFunc("/files/", downloadHandler) // public

    log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
