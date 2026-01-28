package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Email string `json:"email"`
	}

	param := Params{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		msg := fmt.Sprintf("failed to decode params: %v", err)
		respondWithError(w, 500, msg)
		return
	}

	userDb, err := cfg.query.CreateUser(r.Context(), param.Email)
	if err != nil {
		msg := fmt.Sprintf("Failed to create user db: %v", err)
		respondWithError(w, http.StatusInternalServerError, msg)
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
