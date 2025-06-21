package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/mbrunoon/bootdev-chirpy/internal/application"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("error database connection")
	}

	mux := application.NewRouter(db)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}
