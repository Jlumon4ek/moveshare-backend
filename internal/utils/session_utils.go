package utils

import (
	"encoding/json"
	"fmt"
	"moveshare/internal/models"
	"net"
	"net/http"
	"strings"
	"time"
)

// GetClientIP extracts the real client IP from the request
func GetClientIP(r *http.Request) string {
	// Try X-Real-IP first
	ip := r.Header.Get("X-Real-IP")
	if ip != "" && net.ParseIP(ip) != nil {
		return ip
	}

	// Try X-Forwarded-For
	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			trimmedIP := strings.TrimSpace(ips[0])
			if net.ParseIP(trimmedIP) != nil {
				return trimmedIP
			}
		}
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// ParseUserAgent extracts device information from User-Agent string
func ParseUserAgent(userAgent string) models.DeviceInfo {
	if userAgent == "" {
		return models.DeviceInfo{
			Browser:  "Unknown Browser",
			OS:       "Unknown OS",
			Device:   "Unknown Device",
			Platform: "Unknown Platform",
		}
	}

	// Log original User-Agent for debugging
	fmt.Printf("DEBUG: User-Agent: %s\n", userAgent)

	ua := strings.ToLower(userAgent)
	
	// Parse Browser (more detailed detection)
	// Order matters! Check most specific browsers first
	browser := "Unknown Browser"
	if strings.Contains(ua, "yabrowser/") || strings.Contains(ua, "yandexbrowser/") {
		browser = "Yandex Browser"
	} else if strings.Contains(ua, "atomclientelectron/") {
		browser = "Atom Browser" 
	} else if strings.Contains(ua, "sputnik/") {
		browser = "Sputnik Browser"
	} else if strings.Contains(ua, "vivaldi/") {
		browser = "Vivaldi"
	} else if strings.Contains(ua, "brave/") {
		browser = "Brave"
	} else if strings.Contains(ua, "edg/") || strings.Contains(ua, "edge/") {
		browser = "Microsoft Edge"
	} else if strings.Contains(ua, "opera/") || strings.Contains(ua, "opr/") {
		browser = "Opera"
	} else if strings.Contains(ua, "chromium/") {
		browser = "Chromium"
	} else if strings.Contains(ua, "chrome/") {
		browser = "Google Chrome"
	} else if strings.Contains(ua, "firefox/") {
		browser = "Mozilla Firefox"
	} else if strings.Contains(ua, "safari/") && !strings.Contains(ua, "chrome") {
		browser = "Safari"
	}

	// Parse OS
	os := "Unknown OS"
	if strings.Contains(ua, "windows nt") {
		os = "Windows"
	} else if strings.Contains(ua, "macintosh") || strings.Contains(ua, "mac os x") {
		os = "macOS"
	} else if strings.Contains(ua, "linux") {
		os = "Linux"
	} else if strings.Contains(ua, "android") {
		os = "Android"
	} else if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		os = "iOS"
	}

	// Parse Device Type
	device := "Desktop"
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone") {
		device = "Mobile"
	} else if strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad") {
		device = "Tablet"
	}

	// Parse Platform
	platform := "Web"
	if strings.Contains(ua, "mobile") {
		platform = "Mobile Web"
	}

	result := models.DeviceInfo{
		Browser:  browser,
		OS:       os,
		Device:   device,
		Platform: platform,
	}

	// Log parsed result for debugging
	fmt.Printf("DEBUG: Parsed browser: %s, OS: %s\n", browser, os)

	return result
}

// IPGeolocationResponse represents the response from ip-api.com
type IPGeolocationResponse struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	Query       string  `json:"query"`
}

// GetLocationInfo returns location information using ip-api.com (free tier)
func GetLocationInfo(ip string) models.LocationInfo {
	// For localhost and private IPs, get real location based on public IP
	if ip == "::1" || ip == "127.0.0.1" || strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") {
		// Get public IP for real geolocation during development
		return getLocationByPublicIP(ip)
	}

	// Use ip-api.com for geolocation (free tier: 1000 requests/month)
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		// Fallback to unknown if API fails
		return models.LocationInfo{
			Country: "Unknown",
			Region:  "Unknown", 
			City:    "Unknown",
			IP:      ip,
		}
	}
	defer resp.Body.Close()

	var geoResp IPGeolocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return models.LocationInfo{
			Country: "Unknown",
			Region:  "Unknown",
			City:    "Unknown", 
			IP:      ip,
		}
	}

	// Check if API response is successful
	if geoResp.Status != "success" {
		return models.LocationInfo{
			Country: "Unknown",
			Region:  "Unknown",
			City:    "Unknown",
			IP:      ip,
		}
	}

	return models.LocationInfo{
		Country: geoResp.Country,
		Region:  geoResp.RegionName,
		City:    geoResp.City,
		IP:      ip,
	}
}

// getLocationByPublicIP gets location based on the server's public IP for localhost requests
func getLocationByPublicIP(originalIP string) models.LocationInfo {
	// Get location without specifying IP (uses requester's public IP)
	url := "http://ip-api.com/json/"
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		// Fallback for localhost
		return models.LocationInfo{
			Country: "Local",
			Region:  "Local Network",
			City:    "Local",
			IP:      originalIP,
		}
	}
	defer resp.Body.Close()

	var geoResp IPGeolocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return models.LocationInfo{
			Country: "Local",
			Region:  "Local Network",
			City:    "Local",
			IP:      originalIP,
		}
	}

	// Check if API response is successful
	if geoResp.Status != "success" {
		return models.LocationInfo{
			Country: "Local",
			Region:  "Local Network", 
			City:    "Local",
			IP:      originalIP,
		}
	}

	return models.LocationInfo{
		Country: geoResp.Country,
		Region:  geoResp.RegionName,
		City:    geoResp.City,
		IP:      originalIP, // Keep original IP (localhost) but show real location
	}
}