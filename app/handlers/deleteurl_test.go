package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	mock_handlers "url-shortener/app/handlers/mocks"
	"url-shortener/app/repository"
)

func TestHandler_deleteURL(t *testing.T) {
	type mockBehavior func(m *mock_handlers.MockURLDeleter)

	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name         string
		urlPath      string
		behavior     mockBehavior
		expectedBody string
	}{
		{
			name:    "success delete",
			urlPath: "/example123",
			behavior: func(m *mock_handlers.MockURLDeleter) {
				m.EXPECT().DeleteURL("example123").Return(nil)
			},
			expectedBody: `{"status":"OK"}`,
		}, {
			name:    "not found",
			urlPath: "/missing-alias",
			behavior: func(m *mock_handlers.MockURLDeleter) {
				m.EXPECT().DeleteURL("missing-alias").Return(repository.ErrNotFound)
			},
			expectedBody: `{"error":"URL not found", "status":"ERROR"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockURLDeleter(ctrl)
			tt.behavior(m)

			req := httptest.NewRequest(http.MethodDelete, tt.urlPath, nil)

			r := chi.NewRouter()
			w := httptest.NewRecorder()

			r.Delete("/{alias}", DeleteURL(log, m))
			r.ServeHTTP(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.JSONEq(t, tt.expectedBody, string(body))
		})
	}
}
