package repository

import (
	"gorm.io/gorm"

	"sosu-backend/internal/model"
	"time"
)

type AQIRepository struct {
	db *gorm.DB
}

func NewAQIRepository(db *gorm.DB) *AQIRepository {
	return &AQIRepository{db: db}
}

func (r *AQIRepository) Create(a *model.AQIReading) error {
	return r.db.Create(a).Error
}

type nearestResult struct {
	model.AQIReading
	DistanceMeters float64
}

func (r *AQIRepository) FindNearest(lat, lng, maxDistanceMeters float64) (*model.AQIReading, float64, error) {
	var result nearestResult
	err := r.db.Raw(`
		SELECT *,
			ST_Distance(
				geom::geography,
				ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography
			) AS distance_meters
		FROM aqi_readings
		WHERE ST_DWithin(
			geom::geography,
			ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
			?
		)
		ORDER BY geom <-> ST_SetSRID(ST_MakePoint(?, ?), 4326)
		LIMIT 1
	`, lng, lat, lng, lat, maxDistanceMeters, lng, lat).Scan(&result).Error

	if err != nil {
		return nil, 0, err
	}
	if result.ID == 0 {
		return nil, 0, nil
	}
	return &result.AQIReading, result.DistanceMeters, nil
}

func (r *AQIRepository) FindWithinRadius(lat, lng, radiusMeters float64) ([]model.AQIReading, error) {
	var readings []model.AQIReading
	err := r.db.Raw(`
		SELECT * FROM aqi_readings
		WHERE ST_DWithin(
			geom::geography,
			ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
			?
		)
		ORDER BY measured_at DESC
	`, lng, lat, radiusMeters).Scan(&readings).Error

	return readings, err
}
type DailyAQIAverage struct {
	Date    string `gorm:"column:day_date"`
	AvgISPU float64
}

func (r *AQIRepository) DailyAverageWithinRadius(lat, lng, radiusMeters float64, since time.Time) ([]DailyAQIAverage, error) {
	var results []DailyAQIAverage
	err := r.db.Raw(`
		SELECT TO_CHAR(measured_at, 'YYYY-MM-DD') as day_date, AVG(ispu_value) as avg_ispu
		FROM aqi_readings
		WHERE ST_DWithin(
			geom::geography,
			ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
			?
		)
		AND measured_at >= ?
		GROUP BY day_date
		ORDER BY day_date
	`, lng, lat, radiusMeters, since).Scan(&results).Error

	return results, err
}