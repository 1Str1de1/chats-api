package handler_test

import (
	"chats-api/internal/handler"
	"chats-api/internal/model"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockMessagesService struct {
	mock.Mock
}

func (m *MockMessagesService) ValidateMessageCreate(text string) error {
	args := m.Called(text)
	return args.Error(0)
}

func (m *MockMessagesService) CreateMessage(ctx context.Context, text string, chatId int) (*model.Message, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessagesService) GetAllMessagesFromChat(id int, limit int) ([]*model.Message, error) {
	//TODO implement me
	panic("implement me")
}

func TestHandler_HandleMessagesCreate(t *testing.T) {
	tests := []struct {
		name           string
		chatID         string
		requestBody    string
		setupMocks     func(*MockChatsService, *MockMessagesService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "invalid chat id - not a number",
			chatID:         "abc",
			requestBody:    `{"text":"Hello"}`,
			setupMocks:     func(c *MockChatsService, m *MockMessagesService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid chat_id"}`,
		},
		{
			name:           "invalid chat id - zero",
			chatID:         "0",
			requestBody:    `{"text":"Hello"}`,
			setupMocks:     func(c *MockChatsService, m *MockMessagesService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid chat_id"}`,
		},
		{
			name:        "invalid json body",
			chatID:      "1",
			requestBody: `{bad json`,
			setupMocks: func(c *MockChatsService, m *MockMessagesService) {

			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `invalid json body:`,
		},
		{
			name:        "empty message text",
			chatID:      "1",
			requestBody: `{"text":""}`,
			setupMocks: func(c *MockChatsService, m *MockMessagesService) {
				m.On("ValidateMessageCreate", "").
					Return(errors.New("message text is required"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `message text is required`,
		},
		{
			name:        "message text too long",
			chatID:      "1",
			requestBody: `{"text":"` + strings.Repeat("a", 5001) + `"}`,
			setupMocks: func(c *MockChatsService, m *MockMessagesService) {
				m.On("ValidateMessageCreate", strings.Repeat("a", 5001)).
					Return(errors.New("message text is too long"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"message text is too long`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockChats := new(MockChatsService)
			mockMessages := new(MockMessagesService)

			test.setupMocks(mockChats, mockMessages)

			h := handler.NewHandler(mockChats, mockMessages, slog.Default())

			url := fmt.Sprintf("%s/%s/messages", apiPrefix, test.chatID)
			req := httptest.NewRequest(http.MethodPost, url,
				strings.NewReader(test.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			mux := http.NewServeMux()
			mux.HandleFunc("POST "+apiPrefix+"/{id}/messages", h.HandleMessagesCreate())

			mux.ServeHTTP(w, req)
			require.Equal(t, test.expectedStatus, w.Code,
				"Expected status %d, got %d. Response: %s",
				test.expectedStatus, w.Code, w.Body.String())

			if test.expectedBody != "" {
				require.Contains(t, w.Body.String(), test.expectedBody,
					"Expected body to contain: %s\nGot: %s",
					test.expectedBody, w.Body.String())
			}

			mockChats.AssertExpectations(t)
			mockMessages.AssertExpectations(t)
		})
	}
}
