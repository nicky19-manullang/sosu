package model

import "time"

type Hotspot struct {
	ID         uint      `gorm:"primaryKey"`
	Latitude   float64   `gorm:"not null"`
	Longitude  float64   `gorm:"not null"`
	Confidence string    `gorm:"size:20"`
	FRP        float64
	Source     string    `gorm:"size:50;not null"`
	AcquiredAt time.Time `gorm:"not null"`
	CreatedAt  time.Time
}

func (Hotspot) TableName() string {
	return "hotspots"
}