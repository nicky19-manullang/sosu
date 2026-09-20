package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sosu-backend/internal/model"
	"sosu-backend/internal/repository"
	"sosu-backend/internal/service"
)

type openAQLocationsResponse struct {
	Results []struct {
		Name        string `json:"name"`
		Coordinates struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"coordinates"`
		Sensors []struct {
			ID        int `json:"id"`
			Parameter struct {
				Name string `json:"name"`
			} `json:"parameter"`
		} `json:"sensors"`
	} `json:"results"`
}

type openAQMeasurementsResponse struct {
	Results []struct {
		Value  float64 `json:"value"`
		Period struct {
			DatetimeFrom struct {
				UTC string `json:"utc"`
			} `json:"datetimeFrom"`
		} `json:"period"`
	} `json:"results"`
}

type OpenAQFetcher struct {
	apiKey string
	repo   *repository.AQIRepository
}

func NewOpenAQFetcher(apiKey string, repo *repository.AQIRepository) *OpenAQFetcher {
	return &OpenAQFetcher{apiKey: apiKey, repo: repo}
}

func (f *OpenAQFetcher) fetchLatestPM25(sensorID int) (float64, time.Time, bool) {
	url := fmt.Sprintf("https://api.openaq.org/v3/sensors/%d/measurements?limit=1", sensorID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-API-Key", f.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, time.Time{}, false
	}
	defer resp.Body.Close()

	var parsed openAQMeasurementsResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return 0, time.Time{}, false
	}
	if len(parsed.Results) == 0 {
		return 0, time.Time{}, false
	}

	measuredAt, _ := time.Parse(time.RFC3339, parsed.Results[0].Period.DatetimeFrom.UTC)
	return parsed.Results[0].Value, measuredAt, true
}

func (f *OpenAQFetcher) Run() error {
	url := "https://api.openaq.org/v3/locations?bbox=95,-11,141,6&iso=ID&limit=1000"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-API-Key", f.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var parsed openAQLocationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return err
	}

	saved := 0
	for _, loc := range parsed.Results {
		var pm25SensorID int
		for _, s := range loc.Sensors {
			if s.Parameter.Name == "pm25" {
				pm25SensorID = s.ID
				break
			}
		}
		if pm25SensorID == 0 {
			continue
		}

		value, measuredAt, ok := f.fetchLatestPM25(pm25SensorID)
		if !ok || value <= 0 {
			continue
		}

		ispu := service.PM25ToISPU(value)

		reading := &model.AQIReading{
			StationName: loc.Name,
			Latitude:    loc.Coordinates.Latitude,
			Longitude:   loc.Coordinates.Longitude,
			ISPUValue:   ispu,
			Category:    service.ISPUCategory(ispu),
			Source:      "OpenAQ",
			MeasuredAt:  measuredAt,
		}
		if err := f.repo.Create(reading); err == nil {
			saved++
		}
	}

	fmt.Printf("[OpenAQ] saved %d readings with ISPU values\n", saved)
	return nil
}