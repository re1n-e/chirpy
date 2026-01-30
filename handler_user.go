package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	param := Params{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		respondWithError(w, 500, err, "failed to decode params")
		return
	}

	if param.Password == "" {
		respondWithError(w, http.StatusNotAcceptable, nil, "password feild can't be empty")
		return
	}

	hashedPasswd, err := auth.HashPassword(param.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "failed to hash password")
		return
	}

	userDb, err := cfg.query.CreateUser(r.Context(), database.CreateUserParams{
		Email:          param.Email,
		HashedPassword: hashedPasswd,
	})

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

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	param := Params{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		respondWithError(w, 500, err, "failed to decode params")
		return
	}

	if param.Password == "" {
		respondWithError(w, http.StatusNotAcceptable, nil, "password feild can't be empty")
		return
	}

	resp, err := cfg.query.LoginUser(r.Context(), param.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err, "Failed to retrive user")
		return
	}

	match, err := auth.CheckPasswordHash(param.Password, resp.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Failed to check password hash")
		return
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, nil, "Unauthorized login")
		return
	}

	user := User{
		ID:        resp.ID,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
		Email:     resp.Email,
	}

	respondWithJSON(w, http.StatusOK, user)
}
