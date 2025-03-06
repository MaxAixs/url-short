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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mock_handlers.NewMockGetURL(ctrl)
			tt.mockBehavior(m)

			// Настройка роутера Chi
			r := chi.NewRouter()
			r.Get("/{alias}", Redirect(log, m))

			req := httptest.NewRequest(http.MethodGet, tt.urlPath, nil)
			w := httptest.NewRecorder()

			// Выполнение запроса
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			// Проверка статус кода
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			// Проверка Location для успешного редиректа
			if tt.expectedURL != "" {
				location := resp.Header.Get("Location")
				assert.Equal(t, tt.expectedURL, location)
			}
		})
	}
}
