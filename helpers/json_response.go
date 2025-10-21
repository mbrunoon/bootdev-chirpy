package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJson(res http.ResponseWriter, code int, payload interface{}) {
	res.Header().Add("Content-Type", "application/json")

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("[error] json.Marshal(payload): %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(code)
	res.Write(dat)
}

type jsonError struct {
	Error string `json:"error"`
}

func RespondWithError(res http.ResponseWriter, code int, msg string) {
	RespondWithJson(res, code, jsonError{Error: msg})
}
