package service

import (
	"time"

	"sosu-backend/internal/repository"
)

type Status string

const (
	StatusAman    Status = "Aman"
	StatusWaspada Status = "Waspada"
	StatusBahaya  Status = "Tidak Aman"
)

const maxAQIDistanceMeters = 75000

type StatusResult struct {
	FinalStatus   Status     `json:"final_status"`
	AQIStatus     Status     `json:"aqi_status"`
	AQIValue      int        `json:"aqi_value"`
	AQICategory   string     `json:"aqi_category"`
	AQIAvailable  bool       `json:"aqi_available"`
	AQIDistanceKm float64    `json:"aqi_distance_km,omitempty"`
	HotspotStatus Status     `json:"hotspot_status"`
	HotspotCount  int64      `json:"hotspot_count"`
	HotspotTrend  string     `json:"hotspot_trend"`
	DataUpdatedAt *time.Time `json:"data_updated_at,omitempty"`
	Disclaimer    string     `json:"disclaimer"`
}

type StatusService struct {
	hotspotRepo *repository.HotspotRepository
	aqiRepo     *repository.AQIRepository
}

func NewStatusService(hotspotRepo *repository.HotspotRepository, aqiRepo *repository.AQIRepository) *StatusService {
	return &StatusService{hotspotRepo: hotspotRepo, aqiRepo: aqiRepo}
}

func hotspotStatus(highCount, nominalCount int64) Status {
	if highCount > 0 && highCount+nominalCount > 5 {
		return StatusBahaya
	}
	if highCount > 0 || nominalCount > 0 {
		return StatusWaspada
	}
	return StatusAman
}

func worstStatus(a, b Status) Status {
	rank := map[Status]int{StatusAman: 0, StatusWaspada: 1, StatusBahaya: 2}
	if rank[a] >= rank[b] {
		return a
	}
	return b
}

func (s *StatusService) GetStatus(lat, lng float64) (*StatusResult, error) {
	since := time.Now().Add(-24 * time.Hour)

	highCount, err := s.hotspotRepo.CountByConfidenceWithinRadius(lat, lng, 5000, "high", since)
	if err != nil {
		return nil, err
	}
	nominalCount, err := s.hotspotRepo.CountByConfidenceWithinRadius(lat, lng, 10000, "nominal", since)
	if err != nil {
		return nil, err
	}
	totalHotspots, err := s.hotspotRepo.FindWithinRadius(lat, lng, 10000, since)
	if err != nil {
		return nil, err
	}

	hStatus := hotspotStatus(highCount, nominalCount)

	nearest, distanceMeters, err := s.aqiRepo.FindNearest(lat, lng, maxAQIDistanceMeters)
	if err != nil {
		return nil, err
	}

	aqiStatus := StatusAman
	aqiAvailable := nearest != nil
	var ispuValue int
	var category string
	var distanceKm float64

	if aqiAvailable {
		ispuValue = nearest.ISPUValue
		category = nearest.Category
		distanceKm = distanceMeters / 1000
		if !IsAQISafe(ispuValue) {
			aqiStatus = StatusBahaya
		}
	}

	finalStatus := hStatus
	if aqiAvailable {
		finalStatus = worstStatus(aqiStatus, hStatus)
	}

	var dataUpdatedAt *time.Time
	updateDataTimestamp := func(candidate time.Time) {
		if candidate.IsZero() || (dataUpdatedAt != nil && !candidate.After(*dataUpdatedAt)) {
			return
		}
		value := candidate
		dataUpdatedAt = &value
	}
	if aqiAvailable {
		updateDataTimestamp(nearest.MeasuredAt)
	}
	for _, hotspot := range totalHotspots {
		updateDataTimestamp(hotspot.AcquiredAt)
	}

	now := time.Now()
	last24h, _ := s.hotspotRepo.CountWithinRadiusRange(lat, lng, 10000, now.Add(-24*time.Hour), now)
	prev24h, _ := s.hotspotRepo.CountWithinRadiusRange(lat, lng, 10000, now.Add(-48*time.Hour), now.Add(-24*time.Hour))

	trend := "stabil"
	if last24h > prev24h {
		trend = "meningkat"
	} else if last24h < prev24h {
		trend = "menurun"
	}

	return &StatusResult{
		FinalStatus:   finalStatus,
		AQIStatus:     aqiStatus,
		AQIValue:      ispuValue,
		AQICategory:   category,
		AQIAvailable:  aqiAvailable,
		AQIDistanceKm: distanceKm,
		HotspotStatus: hStatus,
		HotspotCount:  int64(len(totalHotspots)),
		HotspotTrend:  trend,
		DataUpdatedAt: dataUpdatedAt,
		Disclaimer:    "Data hotspot merupakan indikasi titik panas dari citra satelit, bukan kepastian kebakaran aktif.",
	}, nil
}
