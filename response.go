package main

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(`{\"error\":\"Something went wrong\"}`))
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}

func respondWithoutBody(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
}
