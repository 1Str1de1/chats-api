package services

import (
	"chats-api/internal/model"
	"chats-api/internal/repository"
	"context"
	"errors"
	"strings"
)

type ChatsService interface {
	ValidateChatCreate(title string) (string, error)
	CreateChat(ctx context.Context, title string) (*model.Chat, error)
	GetChat(id int) (*model.Chat, error)
	DeleteChat(id int) error
}
type chatsService struct {
	repo repository.ChatsRepository
}

func NewChatsRepository(repo repository.ChatsRepository) ChatsService {
	return &chatsService{repo: repo}
}

func (s *chatsService) ValidateChatCreate(title string) (string, error) {
	str := strings.TrimSpace(title)

	if len(str) == 0 {
		return "", errors.New("chat title is required")
	}

	if len(str) > 200 {
		return "", errors.New("chat title cannot be too long")
	}

	return str, nil
}

func (s *chatsService) CreateChat(ctx context.Context, title string) (*model.Chat, error) {
	chat := &model.Chat{Title: title}

	if err := s.repo.Create(ctx, chat); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *chatsService) GetChat(id int) (*model.Chat, error) {
	chat, err := s.repo.Get(id)

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *chatsService) DeleteChat(id int) error {
	return s.repo.Delete(id)
}
