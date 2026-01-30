package repository

import (
	"chats-api/internal/model"
	"context"
)

type ChatsRepository interface {
	Create(ctx context.Context, chat *model.Chat) error
	Get(id int) (*model.Chat, error)
	Delete(id int) error
}
