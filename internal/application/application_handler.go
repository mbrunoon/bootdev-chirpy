package application

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mbrunoon/bootdev-chirpy/internal/helpers"
)

func HealthzHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"status": "OK"}`))
}

func (c *apiConfig) MetricsHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=uft-8")
	res.WriteHeader(http.StatusOK)

	res.Write([]byte(fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`,
		c.fileserverHits.Load())))
}

func (c *apiConfig) ResetMetricsHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	c.fileserverHits.Swap(0)
	res.Write([]byte(fmt.Sprintf("Metrics reseted to %d", c.fileserverHits.Load())))
}

type ChirpParams struct {
	Body string `json:"body"`
}

func ValidateChirpHandler(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := ChirpParams{}
	err := decoder.Decode(&params)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		helpers.JSONResponseError(res, http.StatusInternalServerError, "error decoding params", err)
		return
	}

	if len(params.Body) > 140 {
		helpers.JSONResponseError(res, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedBody := filterBadWords(params.Body)

	type cleanBody struct {
		Body string `json:"cleaned_body"`
	}

	res.WriteHeader(http.StatusOK)
	helpers.JSONResponse(res, http.StatusOK, cleanBody{Body: cleanedBody})
}
