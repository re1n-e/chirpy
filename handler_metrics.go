package main

import (
	"fmt"
	"net/http"
)

func metric(hits int) string {
	return fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`, hits)
}

func (cfg *apiConfig) handleMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, metric(int(cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden Request")
		return
	}
	cfg.fileserverHits.Store(0)
	if err := cfg.query.ResetDb(r.Context()); err != nil {
		msg := fmt.Sprintf("Failed to reset db: %v", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	respondWithJSON(w, 200, "db sucessfully reseted")
}
