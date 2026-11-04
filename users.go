package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) dbHandler(w http.ResponseWriter, r *http.Request) {
	var request User
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Sorry invalid request")
		return
	}

	create, err := cfg.db.CreateUser(r.Context(), request.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not create user")
		return
	}
	response := User{
		ID:        uuid.UUID(create.ID),
		CreatedAt: create.CreatedAt,
		UpdatedAt: create.UpdatedAt,
		Email:     create.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		respondWithError(w, http.StatusBadRequest, "Sorry invalid request")
	}
}
