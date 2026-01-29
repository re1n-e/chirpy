package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirps(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Body   string `json:"body"`
		UserId string `json:"user_id"`
	}

	param := Params{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "failed to decode params")
		return
	}

	if len(param.Body) > 140 {
		respondWithError(w, 400, nil, "Chirp is too long")
		return
	}

	param.Body = cleanse_chirp(param.Body)

	userId, err := uuid.Parse(param.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "failed to parse user id")
		return
	}

	dbChirp, err := cfg.query.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   param.Body,
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "failed to create db for chirp")
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	}

	respondWithJSON(w, 201, chirp)
}

func cleanse_chirp(msg string) string {
	tokens := strings.Fields(msg)
	for i, token := range tokens {
		switch strings.ToLower(token) {
		case "kerfuffle":
			fallthrough
		case "sharbert":
			fallthrough
		case "fornax":
			tokens[i] = "****"
		default:
		}
	}
	return strings.Join(tokens, " ")
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps := []Chirp{}
	resps, err := cfg.query.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "failed to retrive chirps from db")
	}
	for _, resp := range resps {
		chirps = append(chirps, Chirp{
			ID:        resp.ID,
			CreatedAt: resp.CreatedAt,
			UpdatedAt: resp.UpdatedAt,
			Body:      resp.Body,
			UserId:    resp.UserID,
		})
	}
	respondWithJSON(w, 200, chirps)
}
