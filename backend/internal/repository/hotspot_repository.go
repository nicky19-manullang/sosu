package repository

import (
	"time"

	"gorm.io/gorm"

	"sosu-backend/internal/model"
)

type HotspotRepository struct {
	db *gorm.DB
}

func NewHotspotRepository(db *gorm.DB) *HotspotRepository {
	return &HotspotRepository{db: db}
}

func (r *HotspotRepository) Create(h *model.Hotspot) error {
	return r.db.Create(h).Error
}

func (r *HotspotRepository) FindWithinRadius(lat, lng, radiusMeters float64, since time.Time) ([]model.Hotspot, error) {
	var hotspots []model.Hotspot
	err := r.db.Raw(`
		SELECT id, latitude, longitude, confidence, frp, source, acquired_at FROM hotspots
		WHERE ST_DWithin(
			geom::geography,
			ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
			?
		)
		AND acquired_at >= ?
		ORDER BY acquired_at DESC
		LIMIT 3000
	`, lng, lat, radiusMeters, since).Scan(&hotspots).Error

	return hotspots, err
}

func (r *HotspotRepository) CountByConfidenceWithinRadius(lat, lng, radiusMeters float64, confidence string, since time.Time) (int64, error) {
	var count int64
	err := r.db.Raw(`
		SELECT COUNT(*) FROM hotspots
		WHERE ST_DWithin(
			geom::geography,
			ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
			?
		)
		AND confidence = ?
		AND acquired_at >= ?
	`, lng, lat, radiusMeters, confidence, since).Scan(&count).Error

	return count, err
}
func (r *HotspotRepository) FindWithinBBox(minLat, minLng, maxLat, maxLng float64, since time.Time) ([]model.Hotspot, error) {
	var hotspots []model.Hotspot
	err := r.db.Raw(`
		SELECT id, latitude, longitude, confidence, frp, source, acquired_at FROM hotspots
		WHERE geom && ST_MakeEnvelope(?, ?, ?, ?, 4326)
		AND acquired_at >= ?
		ORDER BY acquired_at DESC
		LIMIT 2000
	`, minLng, minLat, maxLng, maxLat, since).Scan(&hotspots).Error

	return hotspots, err
}

func (r *HotspotRepository) DeleteOlderThan(t time.Time) error {
	return r.db.Where("acquired_at < ?", t).Delete(&model.Hotspot{}).Error
}
func (r *HotspotRepository) CountWithinRadiusRange(lat, lng, radiusMeters float64, from, to time.Time) (int64, error) {
	var count int64
	err := r.db.Raw(`
		SELECT COUNT(*) FROM hotspots
		WHERE ST_DWithin(
			geom::geography,
			ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
			?
		)
		AND acquired_at >= ? AND acquired_at < ?
	`, lng, lat, radiusMeters, from, to).Scan(&count).Error

	return count, err
}

type DailyHotspotCount struct {
	Date  string `gorm:"column:day_date"`
	Count int64
}

func (r *HotspotRepository) DailyCountsWithinRadius(lat, lng, radiusMeters float64, since time.Time) ([]DailyHotspotCount, error) {
	var results []DailyHotspotCount
	err := r.db.Raw(`
		SELECT TO_CHAR(acquired_at, 'YYYY-MM-DD') as day_date, COUNT(*) as count
		FROM hotspots
		WHERE ST_DWithin(
			geom::geography,
			ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
			?
		)
		AND acquired_at >= ?
		GROUP BY day_date
		ORDER BY day_date
	`, lng, lat, radiusMeters, since).Scan(&results).Error

	return results, err
}
