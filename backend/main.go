package main

import (
	"log"
	"net/http"
	"os"
)

// List of allowed origins
var allowedOrigins = []string{
	"http://localhost:5173",            // React dev server
	"*",                                // another dev server (if you need it)
	"https://your-production-site.com", // your production frontend
}

// CORS middleware
func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Check if the request's origin is in the allowed list
		for _, o := range allowedOrigins {
			if o == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}

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
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/upload", requireAPIKey(uploadHandler))
	mux.HandleFunc("/files", requireAPIKey(listHandler))
	mux.HandleFunc("/delete/", requireAPIKey(deleteHandler))
	mux.HandleFunc("/files/", downloadHandler) // public

	log.Println("Server running on http://localhost:8080")
	// Wrap mux with CORS middleware here
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))
}
