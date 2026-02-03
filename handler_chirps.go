package main

import (
	"chirpy/internal/auth"
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
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 404, err, "Failed to get bearer token")
		return
	}

	validatedId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, "Failed to validate jwt")
		return
	}

	type Params struct {
		Body string `json:"body"`
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

	dbChirp, err := cfg.query.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   param.Body,
		UserID: validatedId,
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
	id := r.URL.Query().Get("author_id")
	toSort := r.URL.Query().Get("sort")
	chirps := []Chirp{}
	var err error
	var resps []database.Chirp

	if id != "" {
		var authorId uuid.UUID
		authorId, err = uuid.Parse(id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Failed to parse uuid")
			return
		}
		if toSort == "desc" {
			resps, err = cfg.query.GetChirpsByAuthorId(r.Context(), authorId)
		} else {
			resps, err = cfg.query.GetChirpsByAuthorId(r.Context(), authorId)
		}
	} else {
		if toSort == "desc" {
			resps, err = cfg.query.GetAllChirpsDesc(r.Context())
		} else {
			resps, err = cfg.query.GetAllChirps(r.Context())
		}
	}
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

func (cfg *apiConfig) getChirpByChirpId(w http.ResponseWriter, r *http.Request) {
	value := r.PathValue("chirpID")
	if value == "" {
		respondWithError(w, http.StatusBadRequest, nil, "No chirp id provided")
		return
	}
	chirpID, err := uuid.Parse(value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Faield to parse chirp id")
		return
	}

	resp, err := cfg.query.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err, "Failed to retrive chirp from db")
		return
	}

	chirp := Chirp{
		ID:        resp.ID,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
		Body:      resp.Body,
		UserId:    resp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err, "Failed to parse jwt token")
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, err, "Failed to validate JWT")
	}

	value := r.PathValue("chirpID")
	if value == "" {
		respondWithError(w, http.StatusBadRequest, nil, "No chirp id provided")
		return
	}
	chirpID, err := uuid.Parse(value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Faield to parse chirp id")
		return
	}

	resp, err := cfg.query.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err, "Failed to retrive chirp from db")
		return
	}

	if resp.UserID != userId {
		respondWithError(w, http.StatusForbidden, nil, "Unauthorized access to a resource")
		return
	}

	if err := cfg.query.DeleteChirpById(r.Context(), chirpID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Failed to delete chirp")
		return
	}

	respondWithJSON(w, 204, nil)
}
