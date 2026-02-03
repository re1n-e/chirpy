package main

import (
	"chirpy/internal/auth"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (cfg *apiConfig) registerChirpyRed(w http.ResponseWriter, r *http.Request) {
	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, "Failed to get api key")
		return
	}

	if !strings.EqualFold(key, cfg.polkaKey) {
		respondWithError(w, http.StatusUnauthorized, nil, "Unauthorized access")
		return
	}

	type Params struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}

	param := Params{}

	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Failed to parse params")
		return
	}

	if param.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	id, err := uuid.Parse(param.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusNotAcceptable, err, "Failed tp parse uuid")
		return
	}

	if err := cfg.query.UpdatePackageById(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, err, "User not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err, "Failed to update package")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
