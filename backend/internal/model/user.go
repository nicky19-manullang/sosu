package model

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"size:150;uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    time.Time
}

func (User) TableName() string {
	return "users"
}