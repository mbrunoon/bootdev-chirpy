package application

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	mux := mux.NewRouter()

	mux.Handle("/app", http.StripPrefix("/app", http.FileServer(http.Dir("./web/"))))
	mux.PathPrefix("/app/assets").Handler(
		http.StripPrefix("/app/assets", http.FileServer(http.Dir("./assets/"))),
	)

	mux.HandleFunc("/healthz", HealthzHandler)

	return mux
}
