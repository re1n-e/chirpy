package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Email string `json:"email"`
	}

	param := Params{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		respondWithError(w, 500, err, "failed to decode params")
		return
	}

	userDb, err := cfg.query.CreateUser(r.Context(), param.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Failed to create user db")
		return
	}

	user := User{
		ID:        userDb.ID,
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
		Email:     userDb.Email,
	}

	respondWithJSON(w, 201, user)
}
