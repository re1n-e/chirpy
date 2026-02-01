package main

import (
	"chirpy/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, "Failed to get refresh token")
		return
	}

	token, err := cfg.query.GetUserFromToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, "Failed to retrieve refresh token")
		return
	}

	if token.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, nil, "The token has been revoked")
		return
	}

	if time.Now().After(token.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, nil, "The refresh token has expired")
		return
	}

	newJWT, err := auth.MakeJWT(token.UserID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Failed to create JWT")
		return
	}

	type Resp struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, Resp{
		Token: newJWT,
	})
}
