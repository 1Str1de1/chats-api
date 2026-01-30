package handler_test

import (
	"chats-api/internal/handler"
	"chats-api/internal/model"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var apiPrefix = "/api/v1/chats"

type MockChatsService struct {
	mock.Mock
}

func (m *MockChatsService) ValidateChatCreate(title string) (string, error) {
	args := m.Called(title)
	return args.String(0), args.Error(1)
}

func (m *MockChatsService) CreateChat(ctx context.Context, title string) (*model.Chat, error) {
	args := m.Called(title)
	return args.Get(0).(*model.Chat), args.Error(1)
}

func (m *MockChatsService) GetChat(id int) (*model.Chat, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Chat), args.Error(1)
}

func (m *MockChatsService) DeleteChat(id int) error {
	//TODO implement me
	panic("implement me")
}

func TestHandler_HandleChatsCreate(t *testing.T) {

	tests := []struct {
		name          string
		requestedBody string
		setupMock     func(*MockChatsService)
		expectedCode  int
		expectedBody  string
	}{
		{
			name:          "invalid json body",
			requestedBody: "{badJSON)",
			setupMock:     func(m *MockChatsService) {},
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `invalid json body:`,
		},
		{
			name:          "empty title",
			requestedBody: `{"title":""}`,
			setupMock: func(m *MockChatsService) {
				m.On("ValidateChatCreate", "").
					Return("", errors.New("chat title is required"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `chat title is required`,
		},
		{
			name:          "too long title",
			requestedBody: `{"title":"` + strings.Repeat("a", 300) + `"}`,
			setupMock: func(m *MockChatsService) {
				m.On("ValidateChatCreate", strings.Repeat("a", 300)).
					Return("", errors.New("chat title is too long"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `chat title is too long`,
		},
		{
			name:          "successful create chat",
			requestedBody: `{"title":"Family Chat"}`,
			setupMock: func(m *MockChatsService) {
				m.On("ValidateChatCreate", "Family Chat").
					Return("Family Chat", nil)
				m.On("CreateChat", "Family Chat").
					Return(&model.Chat{
						Id:        1,
						Title:     "Family Chat",
						CreatedAt: time.Now(),
					}, nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: `"Id":1,"Title":"Family Chat"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockChats := new(MockChatsService)
			mockMessages := new(MockMessagesService)

			test.setupMock(mockChats)
			h := handler.NewHandler(mockChats, mockMessages, slog.Default())

			req := httptest.NewRequest(http.MethodPost, apiPrefix, strings.NewReader(test.requestedBody))
			w := httptest.NewRecorder()

			h.HandleChatsCreate()(w, req)

			require.Equal(t, test.expectedCode, w.Code)

			if test.expectedBody != "" {
				require.Contains(t, w.Body.String(), test.expectedBody)
			}
			mockChats.AssertExpectations(t)
		})
	}

}
