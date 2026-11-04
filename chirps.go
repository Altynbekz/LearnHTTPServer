package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"httpServerWithGo/internal/database"

	"github.com/google/uuid"
)

type chirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type chirpyResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirpsHandler(w http.ResponseWriter, r *http.Request) {
	var request chirpRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON request")
		return
	}

	// Validate and clean the body before writing anything to the database.
	body, err := validateChirp(request.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	params := database.CreateChirpParams{
		Body:   body,
		UserID: request.UserID,
	}

	create, err := cfg.db.CreateChirp(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not create chirp")
		return
	}

	response := chirpyResponse{
		ID:        uuid.UUID(create.ID),
		CreatedAt: create.CreatedAt,
		UpdatedAt: create.UpdatedAt,
		Body:      create.Body,
		UserID:    create.UserID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("encode chirp response: %v", err)
	}
}

func validateChirp(body string) (string, error) {
	if len(body) > 140 {
		return "", fmt.Errorf("chirp is too long")
	}
	return cleanBody(body), nil
}

func cleanBody(word string) string {
	sliced := strings.Split(word, " ")
	prohibitedWords := []string{"kerfuffle", "sharbert", "fornax"}
	newWord := []string{}
	for _, str := range sliced {
		if slices.Contains(prohibitedWords, strings.ToLower(str)) {
			newWord = append(newWord, "****")
			continue
		}
		newWord = append(newWord, str)
	}
	return strings.Join(newWord, " ")
}
