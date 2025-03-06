package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"url-shortener/app/repository"
)

type URLDeleter interface {
	DeleteURL(shortUrl string) error
}

func DeleteURL(log *slog.Logger, URLDelete URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")
		if aliasIsEmpty(alias) {

			slog.Info("Request alias is empty")

			render.JSON(w, r, "Bad request")

			return
		}

		err := URLDelete.DeleteURL(alias)
		if errors.Is(err, repository.ErrNotFound) {

			slog.Error("URL not found")

			render.JSON(w, r, ErrorResp("URL not found"))

			return
		}

		render.JSON(w, r, OKResponse{Response: OK()})
	}
}
