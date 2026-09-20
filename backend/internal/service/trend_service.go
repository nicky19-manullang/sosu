package service

import (
	"time"

	"sosu-backend/internal/repository"
)

type TrendPoint struct {
	Date         string  `json:"date"`
	HotspotCount int64   `json:"hotspot_count"`
	AQIValue     float64 `json:"aqi_value"`
	AQIAvailable bool    `json:"aqi_available"`
}

type TrendService struct {
	hotspotRepo *repository.HotspotRepository
	aqiRepo     *repository.AQIRepository
}

func NewTrendService(hotspotRepo *repository.HotspotRepository, aqiRepo *repository.AQIRepository) *TrendService {
	return &TrendService{hotspotRepo: hotspotRepo, aqiRepo: aqiRepo}
}

func (s *TrendService) GetTrend(lat, lng float64, days int) ([]TrendPoint, error) {
	since := time.Now().AddDate(0, 0, -days)

	hotspotDaily, err := s.hotspotRepo.DailyCountsWithinRadius(lat, lng, 10000, since)
	if err != nil {
		return nil, err
	}
	aqiDaily, err := s.aqiRepo.DailyAverageWithinRadius(lat, lng, 75000, since)
	if err != nil {
		return nil, err
	}

	hotspotMap := make(map[string]int64)
	for _, h := range hotspotDaily {
		hotspotMap[h.Date] = h.Count
	}
	aqiMap := make(map[string]float64)
	for _, a := range aqiDaily {
		aqiMap[a.Date] = a.AvgISPU
	}

	points := make([]TrendPoint, 0, days)
	now := time.Now()
	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		aqiVal, ok := aqiMap[date]
		points = append(points, TrendPoint{
			Date:         date,
			HotspotCount: hotspotMap[date],
			AQIValue:     aqiVal,
			AQIAvailable: ok,
		})
	}

	return points, nil
}