package model

import "time"

type Message struct {
	Id        int `gorm:"primary key"`
	ChatId    int
	Text      string
	CreatedAt time.Time
}
