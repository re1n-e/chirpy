package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, err error, msg string) {
	type Error struct {
		Err string `json:"error"`
	}
	errMsg := Error{
		Err: fmt.Sprintf("%s: %v", msg, err),
	}
	respondWithJSON(w, code, errMsg)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resp, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to marshall the payload"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}
