package main

import (
	"chirpy/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	query          *database.Queries
	platform       string
	jwtSecret      string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("Failed to load db url")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to load db: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("failed to load jwt secret")
	}

	platform := os.Getenv("PLATFORM")

	dbQueries := database.New(db)

	mux := http.NewServeMux()

	filepathRoot := "."
	assetsPath := "./assets"
	port := "8080"

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		query:          dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
	}

	fileServer := http.FileServer(http.Dir(filepathRoot))
	wrappedFS := cfg.middlewareMetricsInc(fileServer)

	mux.Handle("/app/", http.StripPrefix("/app/", wrappedFS))

	assetsFS := http.FileServer(http.Dir(assetsPath))
	wrappedAssetFs := cfg.middlewareMetricsInc(assetsFS)
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets/", wrappedAssetFs))

	mux.HandleFunc("GET /api/healthz", readiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handleMetric)
	mux.HandleFunc("POST /admin/reset", cfg.handleReset)
	mux.HandleFunc("POST /api/chirps", cfg.createChirps)
	mux.HandleFunc("POST /api/users", cfg.createUser)
	mux.HandleFunc("GET /api/chirps", cfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirpByChirpId)
	mux.HandleFunc("POST /api/login", cfg.loginUser)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
