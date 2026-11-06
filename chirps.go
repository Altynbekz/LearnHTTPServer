package main

import (
	"database/sql"
	"encoding/json"
	"errors"
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

// func to return all the chirps
func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	items, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "sorry internal erro occured try again later")
		return
	}
	var AllChirps []chirpyResponse

	for _, chirp := range items {
		oneChirp := chirpyResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		AllChirps = append(AllChirps, oneChirp)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(AllChirps); err != nil {
		respondWithError(w, http.StatusBadRequest, "Sorry something went wrong try again later")
		return
	}

}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}
	chirp, err := cfg.db.GetChirps(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Could not get chirp")
		return
	}

	response := struct {
		Body string `json:"body"`
	}{}

	response.Body = chirp.Body

	respondWithJSON(w, http.StatusOK, response)
}
