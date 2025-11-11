package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mbrunoon/bootdev-chirpy/helpers"
	"github.com/mbrunoon/bootdev-chirpy/internal/auth"
	"github.com/mbrunoon/bootdev-chirpy/internal/database"
)

type chirpMap struct {
	ID     uuid.UUID `json:"id"`
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type chirpResponse struct {
	chirpMap
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) CreateChirpsController(res http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		helpers.RespondWithError(res, http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.SecretToken)
	if err != nil {
		helpers.RespondWithError(res, http.StatusUnauthorized, err.Error())
		return
	}

	var params struct {
		Body string `json:"body"`
	}

	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		helpers.RespondWithError(res, http.StatusBadRequest, fmt.Sprintf("Error decoding body: %v", err))
		return
	}

	ok, clean := validateAndCleanChirp(params.Body)
	if !ok {
		helpers.RespondWithError(res, http.StatusUnprocessableEntity, "Chirp must be less than 140 chars")
		return
	}

	newChirp, err := cfg.DB.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   clean,
		UserID: userID,
	})
	if err != nil {
		helpers.RespondWithError(res, http.StatusUnprocessableEntity, fmt.Sprintf("Error creating chirp: %v", err))
		return
	}

	helpers.RespondWithJson(res, http.StatusCreated, mapChirp(newChirp))
}

func (cfg *apiConfig) IndexChirpsController(res http.ResponseWriter, req *http.Request) {
	allChirps, _ := cfg.DB.AllChirps(req.Context())
	allChirpsResponse := make([]chirpResponse, len(allChirps))

	for i := range allChirps {
		allChirpsResponse[i] = mapChirp(allChirps[i])
	}

	helpers.RespondWithJson(res, http.StatusOK, allChirpsResponse)
}

func (cfg *apiConfig) ShowChirpController(res http.ResponseWriter, req *http.Request) {
	idParam, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		helpers.RespondWithError(res, http.StatusBadRequest, "Invalid Chirp ID format")
		return
	}

	chirp, err := cfg.DB.FindChirp(req.Context(), idParam)

	if errors.Is(err, sql.ErrNoRows) {
		helpers.RespondWithJson(res, http.StatusNotFound, "Chirp not found")
		return
	}

	if err != nil {
		helpers.RespondWithError(res, http.StatusInternalServerError, fmt.Sprintf("[cfg.DB.FindChirp] %v", err))
		return
	}

	helpers.RespondWithJson(res, http.StatusOK, mapChirp(chirp))
}

func (cfg *apiConfig) DeleteChirpController(res http.ResponseWriter, req *http.Request) {
	idParam, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		helpers.RespondWithError(res, http.StatusBadRequest, "Invalid Chirp Format")
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		helpers.RespondWithError(res, http.StatusUnauthorized, "Invalid Token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.SecretToken)
	if err != nil {
		helpers.RespondWithError(res, http.StatusUnauthorized, "Invalid Token")
		return
	}

	chirp, err := cfg.DB.FindChirp(req.Context(), idParam)

	if err == sql.ErrNoRows {
		helpers.RespondWithError(res, http.StatusNotFound, "chirp not found")
		return
	}

	if userID != chirp.UserID {
		helpers.RespondWithError(res, http.StatusForbidden, "You can only delete your chirps")
		return
	}

	err = cfg.DB.DeleteChirp(req.Context(), idParam)
	if err != nil {
		helpers.RespondWithError(res, http.StatusInternalServerError, "Chirp not deleted")
		return
	}

	helpers.RespondWithJson(res, http.StatusNoContent, "")
}

/* private */

func mapChirp(dbChirp database.Chirp) chirpResponse {
	newChirpMapped := chirpResponse{}
	newChirpMapped.ID = dbChirp.ID
	newChirpMapped.Body = dbChirp.Body
	newChirpMapped.UserID = dbChirp.UserID
	newChirpMapped.CreatedAt = dbChirp.CreatedAt.Time
	newChirpMapped.UpdatedAt = dbChirp.UpdatedAt.Time

	return newChirpMapped
}
