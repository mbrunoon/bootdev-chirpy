package application

import (
	"net/http"
	"sync/atomic"

	"github.com/gorilla/mux"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (c *apiConfig) middlewareMetricsInt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(res, req)
	})
}

func NewRouter() *mux.Router {
	mux := mux.NewRouter()
	cfg := apiConfig{}

	mux.Handle("/app", cfg.middlewareMetricsInt(
		http.StripPrefix("/app", http.FileServer(http.Dir("./web/"))),
	))

	mux.PathPrefix("/app/assets").Handler(
		cfg.middlewareMetricsInt(
			http.StripPrefix("/app/assets", http.FileServer(http.Dir("./assets/"))),
		),
	)

	apiRouter := mux.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/healthz", HealthzHandler).Methods("GET")
	apiRouter.HandleFunc("/validate_chirp", ValidateChirpHandler).Methods("POST")

	adminRouter := mux.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/metrics", cfg.MetricsHandler).Methods("GET")
	adminRouter.HandleFunc("/reset", cfg.ResetMetricsHandler).Methods("POST")

	return mux
}
