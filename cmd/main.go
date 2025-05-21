package main

import (
	"log"
	"net/http"

	"github.com/mbrunoon/bootdev-chirpy/internal/application"
)

func main() {
	mux := application.NewRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}
