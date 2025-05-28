package main

import (
	"database/sql"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/G0SU19O2/Chirpy/internal/database"
	_ "github.com/go-sql-driver/mysql"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_ = database.New(db)

	cfg := apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)

	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	server := http.Server{Handler: mux, Addr: ":8080"}
	server.ListenAndServe()
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
