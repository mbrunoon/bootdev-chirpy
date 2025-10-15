package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync/atomic"

	"github.com/mbrunoon/bootdev-chirpy/helpers"
	"github.com/mbrunoon/bootdev-chirpy/internal/database"

	_ "github.com/lib/pq"
)

func main() {

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("[sql.Open()] %v", err)
		return
	}

	log.Printf("DB Connected...")

	mux := http.NewServeMux()
	apiCfg := apiConfig{
		DB: database.New(db),
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.Handle("/app/assets/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", healthzController)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpController)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsController)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetController)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
	DB             *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(res, req)
	})
}

func (cfg *apiConfig) metricsController(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)

	res.Write([]byte(fmt.Sprintf(`
	<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
	</html>`, cfg.fileserverHits.Load())))
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

func validateChirpController(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "application/json")

	type parameters struct {
		Body string `json:"body"`
	}

	type responseValid struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(res, http.StatusInternalServerError, fmt.Sprintf("Error decoding %s", err))
	}

	if len(params.Body) > 140 {
		helpers.RespondWithError(res, 400, "Chirp is too long")
		return
	}

	helpers.RespondWithJson(res, http.StatusOK, responseValid{CleanedBody: cleanBody(params.Body)})
}

func cleanBody(body string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanedBody := strings.Split(body, " ")

	for i, w := range cleanedBody {
		if slices.Contains(badWords, strings.ToLower(w)) {
			cleanedBody[i] = "****"
		}
	}

	return strings.Join(cleanedBody, " ")
}
