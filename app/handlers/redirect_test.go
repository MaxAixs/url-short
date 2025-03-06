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

func TestHandler_Redirect(t *testing.T) {
	type mockBehavior func(m *mock_handlers.MockGetURL)

	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name         string
		urlPath      string
		mockBehavior mockBehavior
		expectedCode int
		expectedURL  string
	}{
		{
			name:    "Successfully Redirect GET",
			urlPath: "/example",
			mockBehavior: func(m *mock_handlers.MockGetURL) {
				m.EXPECT().GetURL("example").Return("https://original.com", nil)
			},
			expectedCode: http.StatusFound,
			expectedURL:  "https://original.com",
		},
		{
			name:         "Empty Alias",
			urlPath:      "/",
			mockBehavior: func(m *mock_handlers.MockGetURL) {},
			expectedCode: http.StatusNotFound,
		},
		{
			name:    "URL Not Found",
			urlPath: "/non-existent",
			mockBehavior: func(m *mock_handlers.MockGetURL) {
				m.EXPECT().GetURL("non-existent").Return("", repository.ErrNotFound)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockGetURL(ctrl)
			tt.mockBehavior(m)

			r := chi.NewRouter()
			r.Get("/{alias}", Redirect(log, m))

			req := httptest.NewRequest(http.MethodGet, tt.urlPath, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			if tt.expectedURL != "" {
				location := resp.Header.Get("Location")
				assert.Equal(t, tt.expectedURL, location)
			}
		})
	}
}
