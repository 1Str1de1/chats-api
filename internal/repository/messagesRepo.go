package repository

import (
	"chats-api/internal/model"
	"context"
)

type MessagesRepository interface {
	Create(ctx context.Context, message *model.Message) error
	GetAll(chatId int, limit int) ([]*model.Message, error)
}
