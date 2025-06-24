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
	mux.HandleFunc("POST /api/users", middleware.JSONContentType(handlers.HandleCreateUser(cfg)))
	mux.HandleFunc("POST /api/chirps", middleware.JSONContentType(handlers.HandleCreateChirp(cfg)))
	mux.HandleFunc("GET /api/chirps", middleware.JSONContentType(handlers.HandleGetAllChirps(cfg)))
	mux.HandleFunc("GET /api/chirps/{chirpID}", middleware.JSONContentType(handlers.HandleGetChirpById(cfg)))
	mux.HandleFunc("POST /api/login", middleware.JSONContentType(handlers.HandleLoginUser(cfg)))
	mux.HandleFunc("POST /api/refresh", middleware.JSONContentType(handlers.HandleRefreshToken(cfg)))
	mux.HandleFunc("POST /api/revoke", middleware.JSONContentType(handlers.HandleRevokeToken(cfg)))
	return mux
}
