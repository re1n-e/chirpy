package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	filepathRoot := "."
	assetsPath := "./assets"
	port := "8080"

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	fileServer := http.FileServer(http.Dir(filepathRoot))
	wrappedFS := cfg.middlewareMetricsInc(fileServer)

	mux.Handle("/app/", http.StripPrefix("/app/", wrappedFS))

	assetsFS := http.FileServer(http.Dir(assetsPath))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets/", assetsFS))

	mux.HandleFunc("GET /healthz", readiness)
	mux.HandleFunc("GET /metrics", cfg.handleMetric)
	mux.HandleFunc("POST /reset", cfg.handleReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
