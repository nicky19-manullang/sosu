package fetcher

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"sosu-backend/internal/model"
	"sosu-backend/internal/repository"
)

const indonesiaBBox = "95,-11,141,6"

type FIRMSFetcher struct {
	mapKey string
	repo   *repository.HotspotRepository
}

func NewFIRMSFetcher(mapKey string, repo *repository.HotspotRepository) *FIRMSFetcher {
	return &FIRMSFetcher{mapKey: mapKey, repo: repo}
}

func (f *FIRMSFetcher) Run() error {
	url := fmt.Sprintf(
		"https://firms.modaps.eosdis.nasa.gov/api/area/csv/%s/VIIRS_SNPP_NRT/%s/3",
		f.mapKey, indonesiaBBox,
	)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	header, err := reader.Read()
	if err != nil {
		return err
	}
	col := columnIndex(header)

	saved := 0
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		lat, _ := strconv.ParseFloat(row[col["latitude"]], 64)
		lng, _ := strconv.ParseFloat(row[col["longitude"]], 64)
		frp, _ := strconv.ParseFloat(row[col["frp"]], 64)
		confidence := row[col["confidence"]]
		acqDate := row[col["acq_date"]]
		acqTime := row[col["acq_time"]]

		acquiredAt, err := time.Parse("2006-01-02 1504", acqDate+" "+acqTime)
		if err != nil {
			continue
		}
		if !isLikelyLandHotspot(lat, lng) {
			continue
		}

		h := &model.Hotspot{
			Latitude:   lat,
			Longitude:  lng,
			Confidence: mapConfidence(confidence),
			FRP:        frp,
			Source:     "FIRMS-VIIRS",
			AcquiredAt: acquiredAt,
		}

		if err := f.repo.Create(h); err == nil {
			saved++
		}
	}

	fmt.Printf("[FIRMS] saved %d hotspots\n", saved)
	return nil
}

func columnIndex(header []string) map[string]int {
	idx := make(map[string]int)
	for i, name := range header {
		idx[name] = i
	}
	return idx
}

func mapConfidence(raw string) string {
	switch raw {
	case "l":
		return "low"
	case "n":
		return "nominal"
	case "h":
		return "high"
	default:
		return raw
	}
}

func isLikelyLandHotspot(lat, lng float64) bool {
	if lat < -12 || lat > 8 || lng < 94 || lng > 141 {
		return false
	}

	blocks := []struct {
		minLat, maxLat, minLng, maxLng float64
	}{
		{minLat: -4.5, maxLat: 6.2, minLng: 95.0, maxLng: 106.5},
		{minLat: -10.0, maxLat: -5.0, minLng: 104.0, maxLng: 116.0},
		{minLat: -4.8, maxLat: 4.8, minLng: 108.0, maxLng: 119.5},
		{minLat: -8.5, maxLat: 2.2, minLng: 118.2, maxLng: 126.5},
		{minLat: -10.5, maxLat: 3.0, minLng: 126.0, maxLng: 141.0},
		{minLat: -9.5, maxLat: -0.5, minLng: 130.0, maxLng: 141.0},
	}

	for _, block := range blocks {
		if lat >= block.minLat && lat <= block.maxLat && lng >= block.minLng && lng <= block.maxLng {
			return true
		}
	}

	return false
}

const worldArea = "world"

func (f *FIRMSFetcher) RunGlobal() error {
	url := fmt.Sprintf(
		"https://firms.modaps.eosdis.nasa.gov/api/area/csv/%s/VIIRS_SNPP_NRT/%s/1",
		f.mapKey, worldArea,
	)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	header, err := reader.Read()
	if err != nil {
		return err
	}
	col := columnIndex(header)

	saved := 0
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		lat, _ := strconv.ParseFloat(row[col["latitude"]], 64)
		lng, _ := strconv.ParseFloat(row[col["longitude"]], 64)
		frp, _ := strconv.ParseFloat(row[col["frp"]], 64)
		confidence := row[col["confidence"]]
		acqDate := row[col["acq_date"]]
		acqTime := row[col["acq_time"]]

		acquiredAt, err := time.Parse("2006-01-02 1504", acqDate+" "+acqTime)
		if err != nil {
			continue
		}
		if !isLikelyLandHotspot(lat, lng) {
			continue
		}

		h := &model.Hotspot{
			Latitude:   lat,
			Longitude:  lng,
			Confidence: mapConfidence(confidence),
			FRP:        frp,
			Source:     "FIRMS-VIIRS-GLOBAL",
			AcquiredAt: acquiredAt,
		}

		if err := f.repo.Create(h); err == nil {
			saved++
		}
	}

	fmt.Printf("[FIRMS-GLOBAL] saved %d hotspots\n", saved)
	return nil
}
