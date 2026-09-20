package model

import "time"

type AQIReading struct {
	ID          uint      `gorm:"primaryKey"`
	StationName string    `gorm:"size:150"`
	Latitude    float64   `gorm:"not null"`
	Longitude   float64   `gorm:"not null"`
	ISPUValue   int       `gorm:"not null"`
	Category    string    `gorm:"size:30"`
	Source      string    `gorm:"size:50;not null"`
	MeasuredAt  time.Time `gorm:"not null"`
	CreatedAt   time.Time
}

func (AQIReading) TableName() string {
	return "aqi_readings"
}