package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSONResponse(res http.ResponseWriter, code int, payload interface{}) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	resData, err := json.Marshal(payload)

	if err != nil {
		log.Println("error: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(code)
	res.Write(resData)
}

type errorResponse struct {
	Error string `error:"json"`
}

func JSONResponseError(res http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code >= 500 {
		log.Println("error 5XX: ", err)
	}

	JSONResponse(res, code, errorResponse{
		Error: msg,
	})
}
