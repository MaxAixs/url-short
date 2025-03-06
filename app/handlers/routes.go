package handlers

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"url-shortener/app/repository"
)

func InitRoutes(log *slog.Logger, storage *repository.URLRepository) http.Handler {
	r := chi.NewRouter()

	r.Post("/url", SaveURL(log, storage))
	r.Get("/{alias}", Redirect(log, storage))
	r.Delete("/{alias}", DeleteURL(log, storage))

	return r
}
