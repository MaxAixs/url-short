package handlers

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"url-shortener/app/repository"
	"url-shortener/pkg/logger"
)

type ReqURL struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias"`
}

//go:generate mockgen -source=saveurl.go -destination=mocks/mockSaveUrl.go

type URLSaver interface {
	SaveURL(URL string, shortURL string) (int, error)
}

const AliasLen = 5

func SaveURL(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input ReqURL

		err := render.DecodeJSON(r.Body, &input)
		if errors.Is(err, io.EOF) {

			log.Error("Empty request body")

			render.JSON(w, r, ErrorResp("Empty request body"))

			return
		}
		if err != nil {

			log.Error("Error decoding request body", logger.Err(err))

			render.JSON(w, r, ErrorResp("Failed decoding request body"))

			return
		}

		if err := validator.New().Struct(input); err != nil {
			validationErrors := err.(validator.ValidationErrors)

			log.Error("Validation error", logger.Err(err))

			render.JSON(w, r, ValidationError(validationErrors))

			return
		}

		if aliasIsEmpty(input.Alias) {
			input.Alias = generateRandomAlias(AliasLen)
		}

		id, err := urlSaver.SaveURL(input.URL, input.Alias)
		if errors.Is(err, repository.ErrUrlExists) {
			log.Info("URL already exists", slog.String("url", input.URL))

			render.JSON(w, r, ErrorResp("URL already exists"))

			return
		}
		if err != nil {

			log.Error("Error saving URL", logger.Err(err))

			render.JSON(w, r, ErrorResp("Error saving URL"))

			return
		}

		log.Info("url added", slog.Int("id", id))

		render.JSON(w, r, OKResponse{Response: OK(), Alias: input.Alias})
	}
}
