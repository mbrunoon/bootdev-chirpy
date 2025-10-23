package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mbrunoon/bootdev-chirpy/helpers"
	"github.com/mbrunoon/bootdev-chirpy/internal/auth"
)

type loginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) LoginAuthController(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	params := loginParams{}
	err := decoder.Decode(&params)

	if err != nil {
		helpers.RespondWithError(res, http.StatusBadRequest, fmt.Sprintf("Error coding params: %v", err))
		return
	}

	user, err := cfg.DB.FindUserByEmail(req.Context(), sql.NullString{String: params.Email, Valid: true})

	if err == sql.ErrNoRows || err != nil {
		returnAuthError(res)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)

	if err != nil {
		log.Fatalf(`[auth.CheckPasswordHash(params.Password, user.HashedPassword)]: %v`, err)
		returnAuthError(res)
		return
	}

	if !match {
		returnAuthError(res)
		return
	}

	helpers.RespondWithJson(res, http.StatusOK, helpers.MapUser(user))
}

func returnAuthError(res http.ResponseWriter) {
	helpers.RespondWithError(res, http.StatusUnauthorized, "Incorrect email or password")
}
