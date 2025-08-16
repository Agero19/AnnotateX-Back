package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Agero19/AnnotateX-api/internal/config"
	"github.com/Agero19/AnnotateX-api/internal/repository"
	"github.com/Agero19/AnnotateX-api/internal/server/handlers/health"
	"github.com/Agero19/AnnotateX-api/internal/server/handlers/user"
	mwLogger "github.com/Agero19/AnnotateX-api/internal/server/middleware/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

// application struct holds the configuration and repository for the server.
type application struct {
	Config config.Config
	Repo   repository.Repository
	Logger *slog.Logger
}

// NewApp creates a new application instance with the given configuration and repository.
func NewApp(cfg *config.Config, repo repository.Repository, log *slog.Logger) *application {
	return &application{
		Config: *cfg,
		Repo:   repo,
		Logger: log,
	}
}

// Mount mounts the application routes.
func (app *application) Mount() http.Handler {
	r := chi.NewRouter()

	//Middleware can be added here
	r.Use(middleware.RequestID)
	r.Use(mwLogger.New(app.Logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	// Mount routes here
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", health.HealthCheckHandler(app.Logger))
		r.Route("/users", func(r chi.Router) {
			r.Post("/", user.CreateUserHandler(app.Repo.Users, app.Logger))
		})
	})

	return r
}

func (app *application) Run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.Config.Port,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	return srv.ListenAndServe()
}
