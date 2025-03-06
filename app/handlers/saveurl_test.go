package handlers

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	mock_handlers "url-shortener/app/handlers/mocks"
)

func TestHandler_SaveURL(t *testing.T) {
	type mockBehavior func(m *mock_handlers.MockURLSaver)

	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name         string
		requestBody  string
		mockBehavior mockBehavior
		expectedBody string
	}{
		{
			name:        "success save url with provided alias",
			requestBody: `{"url":"https://example.com","alias":"example123"}`,
			mockBehavior: func(m *mock_handlers.MockURLSaver) {
				m.EXPECT().SaveURL("https://example.com", "example123").Return(1, nil)
			},
			expectedBody: `{"status":"OK","alias":"example123"}`,
		},
		{
			name:        "success save url with generated alias",
			requestBody: `{"url":"https://example.com"}`,
			mockBehavior: func(m *mock_handlers.MockURLSaver) {
				m.EXPECT().SaveURL("https://example.com", gomock.Any()).Return(1, nil)
			},
			expectedBody: `{"status":"OK","alias":""}`,
		},
		{
			name:         "empty request body",
			requestBody:  ``,
			mockBehavior: func(m *mock_handlers.MockURLSaver) {},
			expectedBody: `{"status":"ERROR","error":"Empty request body"}`,
		},
		{
			name:        "save url error",
			requestBody: `{"url":"https://example.com","alias":"example123"}`,
			mockBehavior: func(m *mock_handlers.MockURLSaver) {
				m.EXPECT().SaveURL("https://example.com", "example123").Return(0, errors.New("DB error"))
			},
			expectedBody: `{"status":"ERROR","error":"Error saving URL"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSaveURL := mock_handlers.NewMockURLSaver(ctrl)
			tt.mockBehavior(mockSaveURL)

			req := httptest.NewRequest(http.MethodPost, "/save", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			handler := SaveURL(log, mockSaveURL)
			handler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			if tt.name == "success save url with generated alias" {
				var respBody map[string]interface{}
				err := json.Unmarshal(body, &respBody)
				assert.NoError(t, err)
				assert.Equal(t, "OK", respBody["status"])
				assert.NotEmpty(t, respBody["alias"])
			} else {
				assert.JSONEq(t, tt.expectedBody, string(body))
			}
		})
	}
}
