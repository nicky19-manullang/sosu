package model

import "time"

type Subscription struct {
	ID                 uint      `gorm:"primaryKey"`
	UserID             uint      `gorm:"not null;index"`
	Latitude           float64   `gorm:"not null"`
	Longitude          float64   `gorm:"not null"`
	LocationLabel      string    `gorm:"size:150"`
	LastNotifiedStatus string    `gorm:"size:20"`
	LastNotifiedAt     *time.Time
	CreatedAt          time.Time
}

func (Subscription) TableName() string {
	return "subscriptions"
}