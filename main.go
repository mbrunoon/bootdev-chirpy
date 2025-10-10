package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	mux := http.NewServeMux()
	apiCfg := apiConfig{}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.Handle("/app/assets/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))

	mux.HandleFunc("GET /healthz", healthzController)
	mux.HandleFunc("GET /metrics", apiCfg.metricsController)
	mux.HandleFunc("POST /reset", apiCfg.resetController)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(res, req)
	})
}

func (cfg *apiConfig) metricsController(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetController(res http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
}

func healthzController(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(http.StatusText(http.StatusOK)))
}
