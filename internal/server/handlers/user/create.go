package user

import (
	"log/slog"
	"net/http"

	resp "github.com/Agero19/AnnotateX-api/internal/lib/api/response"
	"github.com/Agero19/AnnotateX-api/internal/lib/hash"
	"github.com/Agero19/AnnotateX-api/internal/repository"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=6,max=18"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

type CreateUserResponse struct {
	Response  resp.Response
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
}

func CreateUserHandler(repo repository.Users, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.CreateUserHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req CreateUserRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("Failed to decode request", "error", err)
			render.JSON(w, r, resp.Error("Failed to decode request"))
			return
		}
		log.Info("Request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("Invalid request payload", "error", err)
			render.JSON(w, r, resp.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		hashedPwd, err := hash.HashPassword(req.Password)
		if err != nil {
			log.Error("Failed to hash password", "error", err)
			render.JSON(w, r, resp.Error("Failed to hash password"))
			return
		}

		user := &repository.User{
			ID:        "",
			Username:  req.Username,
			Email:     req.Email,
			Password:  hashedPwd,
			CreatedAt: "",
		}

		if err := repo.Create(user); err != nil {
			log.Error("Failed to create user", "error", err)
			render.JSON(w, r, resp.Error("Failed to create user"))
			return
		}

		log.Info(
			"User created successfully",
			slog.String("user_id", user.ID),
			slog.String("created_at", user.CreatedAt),
		)

		response := CreateUserResponse{
			Response:  resp.OK(),
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
		}

		render.JSON(w, r, response)
	}
}
