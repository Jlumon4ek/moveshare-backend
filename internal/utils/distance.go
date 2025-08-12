package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"moveshare/internal/config"
	"net/http"
	"net/url"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type DistanceElement struct {
	Distance struct {
		Text  string `json:"text"`
		Value int    `json:"value"`
	} `json:"distance"`
	Duration struct {
		Text  string `json:"text"`
		Value int    `json:"value"`
	} `json:"duration"`
	Status string `json:"status"`
}

type DistanceMatrixResponse struct {
	Rows []struct {
		Elements []DistanceElement `json:"elements"`
	} `json:"rows"`
	Status string `json:"status"`
}

type DistanceResult struct {
	Distance      string `json:"distance"`
	DistanceValue int    `json:"distance_value"`
	Duration      string `json:"duration"`
	DurationValue int    `json:"duration_value"`
}

// formatDistance converts meters to kilometers and rounds to 1 decimal place
func formatDistance(meters int) string {
	km := float64(meters) / 1000.0
	rounded := math.Round(km*10) / 10
	return fmt.Sprintf("%.1f km", rounded)
}

func GetDistance(pointA, pointB Point, cfg *config.GoogleMapsConfig) (*DistanceResult, error) {
	baseURL := "https://maps.googleapis.com/maps/api/distancematrix/json"

	params := url.Values{}
	params.Add("origins", fmt.Sprintf("%f,%f", pointA.Lat, pointA.Lng))
	params.Add("destinations", fmt.Sprintf("%f,%f", pointB.Lat, pointB.Lng))
	params.Add("units", "metric")
	params.Add("key", cfg.APIKey)

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var result DistanceMatrixResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	if result.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %s", result.Status)
	}

	if len(result.Rows) == 0 || len(result.Rows[0].Elements) == 0 {
		return nil, fmt.Errorf("no distance data returned")
	}

	element := result.Rows[0].Elements[0]
	if element.Status != "OK" {
		return nil, fmt.Errorf("distance calculation failed: %s", element.Status)
	}

	return &DistanceResult{
		Distance:      formatDistance(element.Distance.Value),
		DistanceValue: element.Distance.Value,
		Duration:      element.Duration.Text,
		DurationValue: element.Duration.Value,
	}, nil
}

func GetDistanceFromAddresses(pickupAddress, deliveryAddress string, cfg *config.GoogleMapsConfig) (*DistanceResult, error) {
	baseURL := "https://maps.googleapis.com/maps/api/distancematrix/json"

	params := url.Values{}
	params.Add("origins", pickupAddress)
	params.Add("destinations", deliveryAddress)
	params.Add("units", "metric")
	params.Add("key", cfg.APIKey)

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	fmt.Printf("Making Google Maps API request to: %s\n", requestURL)

	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	fmt.Printf("Google Maps API response: %s\n", string(body))

	var result DistanceMatrixResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	if result.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %s", result.Status)
	}

	if len(result.Rows) == 0 || len(result.Rows[0].Elements) == 0 {
		return nil, fmt.Errorf("no distance data returned")
	}

	element := result.Rows[0].Elements[0]
	if element.Status != "OK" {
		return nil, fmt.Errorf("distance calculation failed: %s", element.Status)
	}

	return &DistanceResult{
		Distance:      formatDistance(element.Distance.Value),
		DistanceValue: element.Distance.Value,
		Duration:      element.Duration.Text,
		DurationValue: element.Duration.Value,
	}, nil
}
