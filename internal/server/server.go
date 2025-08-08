package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Agero19/AnnotateX-api/internal/config"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type application struct {
	Config config.Config
}

func NewApp(cfg *config.Config) *application {
	return &application{
		Config: *cfg,
	}
}

func (app *application) Mount() http.Handler {
	r := chi.NewRouter()

	//Middleware can be added here
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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

	log.Printf("Starting server on %s", app.Config.Port)

	return srv.ListenAndServe()
}
