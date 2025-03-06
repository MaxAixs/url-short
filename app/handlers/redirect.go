package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"url-shortener/app/repository"
	"url-shortener/pkg/logger"
)

type GetURL interface {
	GetURL(shortUrl string) (string, error)
}

func Redirect(log *slog.Logger, getURL GetURL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")
		if aliasIsEmpty(alias) {

			log.Error("Request Alias is empty")

			render.JSON(w, r, ErrorResp("Bad request"))

			return
		}

		url, err := getURL.GetURL(alias)
		if errors.Is(err, repository.ErrNotFound) {

			log.Info("URL not found", "alias", alias)

			render.JSON(w, r, ErrorResp("URL not found"))

			return
		}
		if err != nil {

			log.Info("cant get URL ", logger.Err(err))

			render.JSON(w, r, ErrorResp("failed to get URL"))

			return
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}
