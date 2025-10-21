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
	decoder := json.NewDecoder(req.Body)
	params := chirpMap{}

	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(res, http.StatusBadRequest, fmt.Sprintf("Error decoding chirpParams: %v", err))
	}

	bodyIsValid, cleanBody := validateAndCleanChirp(params.Body)
	if !bodyIsValid {
		helpers.RespondWithError(res, http.StatusUnprocessableEntity, "Chirp must be less than 140 chars")
		return
	}

	newChirp, err := cfg.DB.CreateChirp(req.Context(), database.CreateChirpParams{Body: cleanBody, UserID: uuid.MustParse(params.UserID.String())})

	if err != nil {
		helpers.RespondWithError(res, http.StatusUnprocessableEntity, fmt.Sprintf("Erro on create new chirp: %v", err))
	}

	newChirpMapped := mapChirp(newChirp)

	helpers.RespondWithJson(res, http.StatusCreated, newChirpMapped)
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
