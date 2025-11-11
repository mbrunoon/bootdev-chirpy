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
	"github.com/mbrunoon/bootdev-chirpy/internal/database"
)

type loginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	helpers.UserMap
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

	if !match {
		returnAuthError(res)
		return
	}

	if err != nil {
		log.Fatalf(`[auth.CheckPasswordHash(params.Password, user.HashedPassword)]: %v`, err)
		returnAuthError(res)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.SecretToken,
		time.Hour,
	)

	if err != nil {
		helpers.RespondWithError(res, http.StatusInternalServerError, fmt.Sprintf("Fail to create JWT: %v", err))
		return
	}

	token := auth.MakeRefreshToken()

	_, err = cfg.DB.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})

	if err != nil {
		helpers.RespondWithError(res, http.StatusInternalServerError, fmt.Sprintf("Error on create Refresh Token: %v", err))
	}

	helpers.RespondWithJson(res, http.StatusOK, loginResponse{
		UserMap:      helpers.MapUser(user),
		Token:        accessToken,
		RefreshToken: token,
	})
}

func (cfg *apiConfig) RefreshTokenController(res http.ResponseWriter, req *http.Request) {

	token, err := auth.GetBearerToken(req.Header)

	if err != nil {
		helpers.RespondWithError(res, http.StatusBadRequest, fmt.Sprint(err))
		return
	}

	user, err := cfg.DB.GetUserFromRefreshToken(req.Context(), token)
	if err != nil {
		helpers.RespondWithError(res, http.StatusUnauthorized, "Invalid Token")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.SecretToken, time.Hour)
	if err != nil {
		helpers.RespondWithError(res, http.StatusUnauthorized, fmt.Sprintf("Error on create access token: %v", err))
		return
	}

	type refreshTokenResponse struct {
		Token string `json:"token"`
	}

	helpers.RespondWithJson(res, http.StatusOK, refreshTokenResponse{Token: accessToken})
}

func (cfg *apiConfig) revokeController(res http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		helpers.RespondWithError(res, http.StatusUnauthorized, "Token not found")
		return
	}

	_, err = cfg.DB.RevokeRefreshToken(req.Context(), refreshToken)
	if err != nil {
		helpers.RespondWithError(res, http.StatusInternalServerError, "error on revoke token")
		return
	}

	helpers.RespondWithJson(res, http.StatusNoContent, "")
}

func returnAuthError(res http.ResponseWriter) {
	helpers.RespondWithError(res, http.StatusUnauthorized, "Incorrect email or password")
}
