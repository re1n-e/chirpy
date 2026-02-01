package main

import (
	"chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) revokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err, "Failed to get bearer token")
		return
	}

	if err := cfg.query.UpdateRefreshToken(r.Context(), token); err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Failed to update table")
		return
	}

	respondWithJSON(w, 204, nil)
}
