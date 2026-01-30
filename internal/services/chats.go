package services

import (
	"chats-api/internal/model"
	"chats-api/internal/repository"
	"context"
	"errors"
	"strings"
)

type ChatsService struct {
	repo repository.ChatsRepository
}

func NewChatsRepository(repo repository.ChatsRepository) *ChatsService {
	return &ChatsService{repo: repo}
}

func (s *ChatsService) ValidateChatCreate(title string) (string, error) {
	str := strings.TrimSpace(title)

	if len(str) == 0 {
		return "", errors.New("chat title is required")
	}

	if len(str) > 200 {
		return "", errors.New("chat title cannot be too long")
	}

	return str, nil
}

func (s *ChatsService) CreateChat(ctx context.Context, title string) (*model.Chat, error) {
	chat := &model.Chat{Title: title}

	if err := s.repo.Create(ctx, chat); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatsService) GetChat(id int) (*model.Chat, error) {
	chat, err := s.repo.Get(id)

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatsService) DeleteChat(id int) error {
	return s.repo.Delete(id)
}
