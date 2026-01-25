package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handleMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %v", cfg.fileserverHits.Load())
}

func (cfg *apiConfig) handleReset(_ http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
}
