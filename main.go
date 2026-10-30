package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {

	mux := http.NewServeMux()
	apiCfg := &apiConfig{}
	mux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(
			http.StripPrefix("/app/", http.FileServer(http.Dir("."))),
		),
	)
	log.Println("The endpoint Altynbek is ready to user")

	mux.Handle("/assets/logo.png", http.FileServer(http.Dir(".")))
	log.Println("The endpoint assets is ready to user")

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	log.Println("The endpoint healthz is ready to user")

	mux.HandleFunc("/about-me", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test/test.html")
	})
	log.Println("The test url is ready to test")

	mux.HandleFunc("/metrics", apiCfg.handler)
	mux.HandleFunc("/reset", apiCfg.resetHandler)
	log.Println("/metrics and /reset are ready ")

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Server is ready to use")
	log.Fatal(server.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Swap(0)
	w.WriteHeader(200)
}

func (cfg *apiConfig) handler(w http.ResponseWriter, req *http.Request) {
	response := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(response))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
