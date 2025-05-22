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

	mux.HandleFunc("/healthz", HealthzHandler)
	mux.HandleFunc("/metrics", cfg.MetricsHandler)
	mux.HandleFunc("/reset", cfg.ResetMetricsHandler)

	return mux
}
