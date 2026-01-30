package services

import (
	"chats-api/internal/model"
	"chats-api/internal/repository"
	"context"
	"errors"
)

type messagesService struct {
	repo repository.MessagesRepository
}

type MessagesService interface {
	ValidateMessageCreate(text string) error
	CreateMessage(ctx context.Context, text string, chatId int) (*model.Message, error)
	GetAllMessagesFromChat(id int, limit int) ([]*model.Message, error)
}

func NewMessagesRepository(repo repository.MessagesRepository) MessagesService {
	return &messagesService{repo: repo}
}

func (s *messagesService) ValidateMessageCreate(text string) error {
	if len(text) == 0 {
		return errors.New("message is empty")
	}

	if len(text) > 5000 {
		return errors.New("message is too long")
	}

	return nil
}

func (s *messagesService) CreateMessage(ctx context.Context, text string, chatId int) (*model.Message, error) {
	message := &model.Message{Text: text, ChatId: chatId}

	if err := s.repo.Create(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *messagesService) GetAllMessagesFromChat(id int, limit int) ([]*model.Message, error) {
	messages, err := s.repo.GetAll(id, limit)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
