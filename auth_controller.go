package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mbrunoon/bootdev-chirpy/helpers"
	"github.com/mbrunoon/bootdev-chirpy/internal/auth"
)

type loginParams struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
}

type loginResponse struct {
	User         helpers.UserMap
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
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

	expirationTime := time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.SecretToken,
		expirationTime,
	)

	if err != nil {
		helpers.RespondWithError(res, http.StatusInternalServerError, fmt.Sprintf("Fail to create JWT: %v", err))
		return
	}

	if !match {
		returnAuthError(res)
		return
	}

	helpers.RespondWithJson(res, http.StatusOK, loginResponse{
		User:  helpers.MapUser(user),
		Token: accessToken,
	})
}

func returnAuthError(res http.ResponseWriter) {
	helpers.RespondWithError(res, http.StatusUnauthorized, "Incorrect email or password")
}
