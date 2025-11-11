package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/mbrunoon/bootdev-chirpy/helpers"
	"github.com/mbrunoon/bootdev-chirpy/internal/auth"
	"github.com/mbrunoon/bootdev-chirpy/internal/database"
)

func (cfg *apiConfig) PolkaUserEventController(res http.ResponseWriter, req *http.Request) {

	apiToken, err := auth.GetAPIKey(req.Header)
	if err != nil {
		helpers.RespondWithError(res, http.StatusUnauthorized, fmt.Sprintf("Error: %v", err))
		return
	}

	if apiToken != cfg.PolkaKey {
		helpers.RespondWithError(res, http.StatusUnauthorized, "invalid api token")
		return
	}

	var params struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		helpers.RespondWithError(res, http.StatusInternalServerError, "")
		return
	}

	if params.Event != "user.upgraded" {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		helpers.RespondWithError(res, http.StatusBadRequest, "error parsing uuid")
		return
	}

	_, err = cfg.DB.FindUserByID(req.Context(), userID)
	if err == sql.ErrNoRows {
		helpers.RespondWithError(res, http.StatusNotFound, "")
		return
	}

	_, err = cfg.DB.UpdateUserToRed(req.Context(), database.UpdateUserToRedParams{
		ID:          userID,
		IsChirpyRed: sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		helpers.RespondWithError(res, http.StatusInternalServerError, "error updating user")
		return
	}

	res.WriteHeader(http.StatusNoContent)
}
