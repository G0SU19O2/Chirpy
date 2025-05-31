package handlers

import (
	"fmt"
	"net/http"

	"github.com/G0SU19O2/Chirpy/internal/config"
)

func HandleMetrics(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf(`<html>
  <body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.FileserverHits.Load())))
	}
}
