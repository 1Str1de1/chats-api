package model

import "time"

type Chat struct {
	Id        int `gorm:"primary key"`
	Title     string
	CreatedAt time.Time
}
