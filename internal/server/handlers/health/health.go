package health

// health check response and handler

import (
	"log/slog"
	"net/http"

	resp "github.com/Agero19/AnnotateX-api/internal/lib/api/response"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// HealthCheckResponse represents the response structure for the health check.
type HealthCheckResponse struct {
	// default response - status + error
	Response resp.Response `json:"response"`
	Version  string        `json:"version"`
}

// HealthCheckHandler handles the health check requests.
func HealthCheckHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.health.HealthCheckHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		response := HealthCheckResponse{
			Response: resp.OK(),
			Version:  "1.0.0", // TODO: Update version dynamically
		}

		render.JSON(w, r, response)
	}
}
