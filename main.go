package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/mbrunoon/bootdev-chirpy/helpers"
	"github.com/mbrunoon/bootdev-chirpy/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	DB             *database.Queries
	platform       string
	SecretToken    string
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("[sql.Open()] %v", err)
		return
	}

	log.Printf("DB Connected...")
	log.Printf("Plataform: %v", os.Getenv("PLATFORM"))

	mux := http.NewServeMux()
	apiCfg := apiConfig{
		DB:          database.New(db),
		platform:    os.Getenv("PLATFORM"),
		SecretToken: os.Getenv("SECRET_TOKEN"),
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.Handle("/app/assets/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", healthzController)
	mux.HandleFunc("POST /api/users", apiCfg.CreateUserController)
	mux.HandleFunc("POST /api/chirps", apiCfg.CreateChirpsController)
	mux.HandleFunc("GET /api/chirps", apiCfg.IndexChirpsController)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.ShowChirpController)

	mux.HandleFunc("POST /api/login", apiCfg.LoginAuthController)
	mux.HandleFunc("POST /api/refresh", apiCfg.RefreshTokenController)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeController)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsController)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetController)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
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
	if cfg.platform != "dev" {
		helpers.RespondWithError(res, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cfg.DB.DeleteAllUsers(req.Context())
	cfg.DB.DeleteAllChirps(req.Context())
	cfg.fileserverHits.Store(0)

	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
}

func healthzController(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(http.StatusText(http.StatusOK)))
}

func validateAndCleanChirp(chirp string) (bool, string) {
	if len(chirp) > 140 {
		return false, chirp
	}

	return true, cleanBody(chirp)
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
