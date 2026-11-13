package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func projectFile(path string) string {
	for _, candidate := range []string{path, filepath.Join("..", path)} {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return path
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" || origin == "http://localhost:5174" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func newMux() *http.ServeMux {
	return newMuxWithConfig(&apiConfig{})
}

func newMuxWithConfig(apiCfg *apiConfig) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(
			http.StripPrefix("/app", http.FileServer(http.Dir("."))),
		),
	)
	log.Println("The endpoint Altynbek is ready to user")

	mux.Handle(
		"GET /app/assets/",
		apiCfg.middlewareMetricsInc(
			http.StripPrefix("/app/assets/", http.FileServer(http.Dir("assets"))),
		),
	)
	log.Println("The endpoint assets is ready to user")

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	log.Println("The endpoint healthz is ready to user")

	mux.HandleFunc("/about-me", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, projectFile("test/test.html"))
	})
	log.Println("The test url is ready to test")

	mux.HandleFunc("GET /api/metrics", apiCfg.handler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	log.Println("/metrics and /reset are ready ")

	mux.HandleFunc("POST /api/users", apiCfg.dbHandler)
	log.Println("The db handler has been created successfully")

	mux.HandleFunc("POST /api/chirps", apiCfg.createChirpsHandler)
	log.Println("The /api/chirps is ready to use")

	mux.HandleFunc("GET /api/chirps", apiCfg.getAllChirpsHandler)
	log.Println("get all chips has been created")

	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirp)
	log.Println("Getting chirp by id has been created")

	mux.HandleFunc("POST /api/login", apiCfg.UserLoginHandler)
	log.Println("The login endpoint is ready")

	mux.HandleFunc("POST /api/refresh", apiCfg.UserRefreshHandler)
	log.Println("The refresh endpoint is ready")

	mux.HandleFunc("POST /api/revoke", apiCfg.revokeHandler)
	log.Println("Revoke endpoint is ready to use")

	mux.HandleFunc("PUT /api/users", apiCfg.UpdateUserHandler)
	log.Println("the updating endpoint is ready")

	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.DeleteChirpHandler)
	log.Println("The delete endpoint has been created")

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.webhookHandler)
	log.Println("The webhook endpoint has been created")

	return mux
}
