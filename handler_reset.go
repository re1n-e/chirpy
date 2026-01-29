package main

import (
	"net/http"
)

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, nil, "Forbidden Request")
		return
	}
	cfg.fileserverHits.Store(0)
	if err := cfg.query.ResetDb(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Failed to reset db")
		return
	}

	respondWithJSON(w, 200, "db sucessfully reseted")
}
