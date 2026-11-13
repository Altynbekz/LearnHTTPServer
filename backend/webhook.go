package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"httpServerWithGo/internal/auth"
	"httpServerWithGo/internal/database"
	"net/http"

	"github.com/google/uuid"
)

type WebHook struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) webhookHandler(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Couldn't recognize the api key")
		return
	}

	if apiKey != cfg.PolkaKey {
		respondWithError(w, 401, "Incorrect key try a different key")
		return
	}

	webHook := &WebHook{}

	if err := json.NewDecoder(r.Body).Decode(webHook); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the request")
		return
	}

	if webHook.Event != "user.upgraded" {
		respondWithoutBody(w, 204)
		return
	}

	exists, err := cfg.db.ExistsByUserId(r.Context(), webHook.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, "The user cant be found")
			return
		}
		respondWithError(w, 500, "Sorry error has occured")
		return
	}

	if !exists {
		respondWithError(w, 404, "The user cant be found")
		return
	}

	newStatus := database.UpdateChripStatusParams{
		ID:          webHook.Data.UserID,
		IsChirpyRed: true,
	}

	if err := cfg.db.UpdateChripStatus(r.Context(), newStatus); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Sorry, the user status could not be updated")
		return
	}

	respondWithoutBody(w, http.StatusNoContent)
}
