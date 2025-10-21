package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mbrunoon/bootdev-chirpy/helpers"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) CreateUserController(res http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		helpers.RespondWithError(res, http.StatusBadRequest, fmt.Sprintf("Error decoding: %v", err))
		return
	}

	user, err := cfg.DB.CreateUser(req.Context(), sql.NullString{String: params.Email, Valid: true})
	if err != nil {
		helpers.RespondWithError(res, http.StatusUnprocessableEntity, fmt.Sprintf("Error create user: %v", err))
		return
	}

	helpers.RespondWithJson(res, http.StatusCreated, UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email.String,
	})

}
