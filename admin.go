package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"httpServerWithGo/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Swap(0)
	if cfg.platform != "dev" {
		w.WriteHeader(403)
		return
	}

	if err := cfg.db.DeleteAll(r.Context()); err != nil {
		respondWithError(w, http.StatusBadRequest, "Sorry something went wrong")
		return
	}
	w.WriteHeader(200)
}

func (cfg *apiConfig) handler(w http.ResponseWriter, req *http.Request) {
	response := fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>
	`, cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(response))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
