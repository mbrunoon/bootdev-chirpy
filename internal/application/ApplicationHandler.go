package application

import (
	"fmt"
	"net/http"
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
