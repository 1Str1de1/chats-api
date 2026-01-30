package repository

import (
	"chats-api/internal/model"
	"context"

	"gorm.io/gorm"
)

type chatsRepo struct {
	db *gorm.DB
}

func NewChatsRepo(db *gorm.DB) ChatsRepository {
	return &chatsRepo{db: db}
}

func (r *chatsRepo) Create(ctx context.Context, chat *model.Chat) error {
	return r.db.WithContext(ctx).Create(chat).Error
}

func (r *chatsRepo) Get(id int) (*model.Chat, error) {
	var chat model.Chat

	if err := r.db.First(&chat, id).Error; err != nil {
		return nil, ErrChatNotFound
	}

	return &chat, nil
}

func (r *chatsRepo) Delete(id int) error {
	result := r.db.Delete(&model.Chat{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrChatNotFound
	}
	return nil
}
