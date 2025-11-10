package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mbrunoon/bootdev-chirpy/helpers"
	"github.com/mbrunoon/bootdev-chirpy/internal/auth"
	"github.com/mbrunoon/bootdev-chirpy/internal/database"
)

type parameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) CreateUserController(res http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		helpers.RespondWithError(res, http.StatusBadRequest, fmt.Sprintf("Error decoding: %v", err))
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		helpers.RespondWithError(res, http.StatusUnprocessableEntity, fmt.Sprintf("Password Error: %v", err))
		return
	}

	userParam := database.CreateUserParams{
		Email:          sql.NullString{String: params.Email, Valid: true},
		HashedPassword: hashedPass,
	}

	user, err := cfg.DB.CreateUser(req.Context(), userParam)
	if err != nil {
		helpers.RespondWithError(res, http.StatusUnprocessableEntity, fmt.Sprintf("Error create user: %v", err))
		return
	}

	helpers.RespondWithJson(res, http.StatusCreated, helpers.MapUser(user))

}

func (cfg *apiConfig) UpdateUserController(res http.ResponseWriter, req *http.Request) {

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		helpers.RespondWithError(res, http.StatusUnauthorized, "Invalid token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.SecretToken)
	if err != nil {
		helpers.RespondWithError(res, http.StatusUnauthorized, "Invalid token")
		return
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(res, http.StatusBadRequest, "error on decode params")
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		helpers.RespondWithError(res, http.StatusInternalServerError, "error on hash password")
		return
	}

	user, err := cfg.DB.UpdateUser(req.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          sql.NullString{String: params.Email, Valid: true},
		HashedPassword: hashedPass,
	})
	if err != nil {
		helpers.RespondWithError(res, http.StatusInternalServerError, "error on update user")
		return
	}

	helpers.RespondWithJson(res, http.StatusOK, helpers.MapUser(user))
}
