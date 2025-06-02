package router

import (
	"net/http"

	"github.com/G0SU19O2/Chirpy/internal/config"
	"github.com/G0SU19O2/Chirpy/internal/handlers"
	"github.com/G0SU19O2/Chirpy/internal/middleware"
)

func SetupRoutes(cfg *config.Config) *http.ServeMux {
	mux := http.NewServeMux()

	// Static file serving with metrics middleware
	mux.Handle("/app/", middleware.MetricsInc(cfg)(http.StripPrefix("/app/", http.FileServer(http.Dir("web/static")))))

	// Admin routes
	mux.HandleFunc("GET /admin/metrics", handlers.HandleMetrics(cfg))
	mux.HandleFunc("POST /admin/reset", handlers.HandleReset(cfg))

	// API routes
	mux.HandleFunc("GET /api/healthz", handlers.HandleReadiness)
	mux.HandleFunc("POST /api/chirps", handlers.HandleCreateChip(cfg))
	mux.HandleFunc("POST /api/users", handlers.HandleCreateUser(cfg))

	return mux
}
