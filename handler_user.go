package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

const refreshExpiryTime = 30 * 24 * time.Hour

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

	token, err := auth.MakeJWT(resp.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Failed to generate JWT")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Failed to generate refresh token")
		return
	}

	if err := cfg.query.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    resp.ID,
		ExpiresAt: time.Now().Add(refreshExpiryTime),
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Failed to save refresh token to db")
	}

	user := User{
		ID:           resp.ID,
		CreatedAt:    resp.CreatedAt,
		UpdatedAt:    resp.UpdatedAt,
		Email:        resp.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (cfg *apiConfig) updateUsers(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, 401, err, "Failed to get token")
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, err, "Failed to validate JWT")
		return
	}

	param := Params{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Failed to parse params")
		return
	}

	if param.Email == "" {
		respondWithError(w, http.StatusNoContent, nil, "Email fieldl can't be empty")
		return
	}

	if param.Password == "" {
		respondWithError(w, http.StatusNoContent, nil, "Password feild requied")
		return
	}

	hashedPassword, err := auth.HashPassword(param.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Failed to hash password")
		return
	}

	updatedUser, err := cfg.query.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          param.Email,
		HashedPassword: hashedPassword,
		ID:             userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Failed to update users")
		return
	}

	user := User{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	}

	respondWithJSON(w, http.StatusOK, user)
}

