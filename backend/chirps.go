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

	"httpServerWithGo/internal/auth"
	"httpServerWithGo/internal/database"

	"github.com/google/uuid"
)

type chirpRequest struct {
	Body string `json:"body"`
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

	// The token identifies the author; the request body cannot choose a user.
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid Bearer token")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
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
		UserID: userID,
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
	id := r.URL.Query().Get("author_id")
	sortType := r.URL.Query().Get("sort")

	var items []database.Chirp
	if id == "" {
		allChirps, err := cfg.db.GetAllChirps(r.Context(), sortType)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, 400, "No chirps found")
				return
			}
			respondWithError(w, 500, "Error occured in the server")
			return
		}
		items = allChirps
	} else {
		id, err := uuid.Parse(id)
		if err != nil {
			respondWithError(w, 500, "Couldn't parse the id")
			return
		}
		params := database.GetChirpBySortUserIdParams{UserID: id, SortType: sortType}
		oneChirp, err := cfg.db.GetChirpBySortUserId(r.Context(), params)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, 400, "No chirp was found")
				return
			}
			respondWithError(w, 500, "Internal error occured")
			return
		}
		items = oneChirp
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

func (cfg *apiConfig) DeleteChirpHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("chirpID")
	database := cfg.db
	chirpId, err := uuid.Parse(id)

	if err != nil {
		respondWithError(w, 401, "Couldn't find the user with that id")
		return
	}

	chirp, err := database.GetChirps(r.Context(), chirpId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 400, "Chirp not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error occured couldn't find the chirp")
		return
	}

	//check if the user is eglible to delete that chirp
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Couldn't get the token")
		return
	}

	validate, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, 400, "Couldn't validate the token")
		return
	}

	if validate != chirp.UserID {
		respondWithError(w, 403, "Sorry you dont have access")
		return
	}

	//fund if the user with that id exists
	exists, err := database.ExistsByChirpId(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, 401, "Couldn't find the user with that id in the database")
		return
	}

	if !exists {
		respondWithError(w, 404, "The chirp does not exist")
		return
	}

	if err := database.DeleteChirp(r.Context(), chirpId); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Sorry error occured couldn't delete the chirp")
		return
	}

	respondWithoutBody(w, 204)

}
