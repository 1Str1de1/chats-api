package repository

import (
	"chats-api/internal/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type messagesRepo struct {
	db *gorm.DB
}

var ErrChatNotFound = errors.New("chat not found")

func NewMessagesRepo(db *gorm.DB) MessagesRepository {
	return &messagesRepo{db: db}
}

func (r *messagesRepo) Create(ctx context.Context, message *model.Message) error {
	if err := r.db.Where("chat_id = ?", message.ChatId).Error; err != nil {
		return ErrChatNotFound
	}

	return r.db.Create(message).WithContext(ctx).Error
}

func (r *messagesRepo) GetAll(chatId int, limit int) ([]*model.Message, error) {
	var messages []*model.Message

	result := r.db.Where("chat_id = ?", chatId).
		Limit(limit).
		Order("created_at desc").
		Find(&messages)

	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}
