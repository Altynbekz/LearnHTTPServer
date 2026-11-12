package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"httpServerWithGo/internal/auth"
	"httpServerWithGo/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) dbHandler(w http.ResponseWriter, r *http.Request) {
	var request UserRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Sorry invalid request")
		return
	}

	if request.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}

	if request.Password == "" {
		respondWithError(w, http.StatusBadRequest, "The password is invalid")
		return
	}

	request.Password, err = auth.HashPassword(request.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not save the password")
		return
	}

	created := database.CreateUserParams{
		Email:          request.Email,
		HashedPassword: request.Password,
	}
	create, err := cfg.db.CreateUser(r.Context(), created)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user")
		return
	}

	response := userResponse{
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

func (cfg *apiConfig) UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	var login UserLogin

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON request")
		return
	}

	getUser, err := cfg.db.FindUserByEmail(r.Context(), login.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Could not look up user")
		return
	}

	entered, err := auth.CheckPasswordHash(login.Password, getUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "soory internal error occured")
		return
	}

	if !entered {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(getUser.ID, cfg.secretKey, accessTokenDuration(3600))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could create the token sorry")
		return
	}

	refreshToken := auth.MakeRefreshToken()
	saveToken := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    getUser.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	}

	err = cfg.db.CreateRefreshToken(r.Context(), saveToken)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn`t save rf try again later")
		return
	}

	response := userResponse{
		ID:           getUser.ID,
		CreatedAt:    getUser.CreatedAt,
		UpdatedAt:    getUser.UpdatedAt,
		Email:        getUser.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, 200, response)
}

func accessTokenDuration(seconds int) time.Duration {
	if seconds <= 0 || seconds > 3600 {
		return time.Hour
	}
	return time.Duration(seconds) * time.Second
}

func (cfg *apiConfig) UserRefreshHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't identify the refresh token")
		return
	}

	isMatch, err := cfg.db.FindRefreshToken(r.Context(), token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 401, "invalid refresh token")
			return
		}
		respondWithError(w, http.StatusUnauthorized, "couldn't identify the refresh token")
		return
	}

	if isMatch.RevokedAt.Valid || time.Now().After(isMatch.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}
	newJwt, err := auth.MakeJWT(isMatch.UserID, cfg.secretKey, accessTokenDuration(3600))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create the jwt token try again later")
		return
	}

	respondWithJSON(w, 200, map[string]any{"token": newJwt})
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Couldn't identify the token")
		return
	}
	is_match, err := cfg.db.ExistsByRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "couldn't identify the token")
		return
	}

	if !is_match {
		respondWithError(w, 401, "Invalid token")
		return
	}

	if err := cfg.db.RevokeRefreshToken(r.Context(), token); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke try again later")
		return
	}

	respondWithoutBody(w, 204)

}

func (cfg *apiConfig) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Couldn't get the token")
		return
	}
	uid, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, 401, "Couldn't verify the token ")
		return
	}

	var updatedUser UserRequest
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the user")
		return
	}

	hashed_password, err := auth.HashPassword(updatedUser.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash the password")
		return
	}

	updated := database.UpdateUserParams{
		Email:          updatedUser.Email,
		HashedPassword: hashed_password,
		ID:             uid,
	}

	if err := cfg.db.UpdateUser(r.Context(), updated); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update the user")
		return
	}

	respondWithJSON(w, 200, map[string]any{"email": updatedUser.Email})

}
