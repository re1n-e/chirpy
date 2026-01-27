package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		Body string `json:"body"`
	}

	chirp := Chirp{}
	if err := json.NewDecoder(r.Body).Decode(&chirp); err != nil {
		msg := fmt.Sprintf("Failed to decode chirp: %v", err)
		respondWithError(w, 500, msg)
		return
	}

	if len(chirp.Body) > 140 {
		msg := "Chirp is too long"
		respondWithError(w, 400, msg)
		return
	}

	type Resp struct {
		Valid        bool   `json:"valid"`
		Cleaned_Body string `json:"cleaned_body"`
	}

	respondWithJSON(w, 200, Resp{
		Valid:        true,
		Cleaned_Body: cleanse_chirp(chirp.Body),
	})
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
