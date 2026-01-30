package services

import (
	"chats-api/internal/model"
	"context"
)

type ChatService struct {
	repo TaskRepository
}

type TaskRepository interface {
	Create(ctx context.Context, chat *model.Chat) error
	Get(id int) (*model.Chat, error)
}
