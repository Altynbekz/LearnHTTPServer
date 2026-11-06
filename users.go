package main

import (
	"encoding/json"
	"httpServerWithGo/internal/auth"
	"httpServerWithGo/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type UserRequest struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) dbHandler(w http.ResponseWriter, r *http.Request) {
	var request UserRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Sorry invalid request")
		return
	}

	if request.Password == "" {
		respondWithError(w, 401, "The password is invalid")
		return
	}

	request.Password, err = auth.HashPassword(request.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "sorry could save the password")
		return
	}

	created := database.CreateUserParams{
		Email:          request.Email,
		HashedPassword: request.Password,
	}
	create, err := cfg.db.CreateUser(r.Context(), created)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not create user")
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
		respondWithError(w, http.StatusInternalServerError, "internal error sorry")
		return
	}

	getUser, err := cfg.db.FindUserByEmail(r.Context(), login.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "user not found")
		return
	}

	entered, err := auth.CheckPasswordHash(login.Password, getUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "soory internal error occured")
		return
	}

	if !entered {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	response := userResponse{
		ID:        getUser.ID,
		CreatedAt: getUser.CreatedAt,
		UpdatedAt: getUser.UpdatedAt,
		Email:     getUser.Email,
	}

	respondWithJSON(w, 200, response)
}
