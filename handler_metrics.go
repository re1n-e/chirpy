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
