package widgets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type WeatherTickMsg time.Time

type WeatherWidget struct {
	location location
	weather  weather
	hasError bool
}

type location struct {
	City      string
	Latitude  float64
	Longitude float64
}

type weather struct {
	Temperature   float64
	Precipitation float64
	WeatherCode   int
	WindSpeed     float64
	IsDay         int
}

type locationResponse struct {
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type weatherResponse struct {
	Current struct {
		Temperature   float64 `json:"temperature_2m"`
		Precipitation float64 `json:"precipitation"`
		WeatherCode   int     `json:"weathercode"`
		WindSpeed     float64 `json:"wind_speed_10m"`
		IsDay         int     `json:"is_day"`
	} `json:"current"`
}

var weatherCodeDescriptions map[string]weatherCodeDesc

type weatherCodeDesc struct {
	Day struct {
		Description string `json:"description"`
		Ascii       string `json:"ascii"`
	} `json:"day"`
	Night struct {
		Description string `json:"description"`
		Ascii       string `json:"ascii"`
	} `json:"night"`
}

func weatherIcon(code int, isDay int) string {
	key := strconv.Itoa(code)
	desc, ok := weatherCodeDescriptions[key]
	if !ok {
		return "[?]"
	}

	if isDay == 1 {
		return desc.Day.Ascii
	}
	return desc.Night.Ascii
}

func init() {
	weatherCodeDescriptions = make(map[string]weatherCodeDesc)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Printf("unable to determine widget source file path")
		return
	}

	jsonPath := filepath.Join(filepath.Dir(filename), "..", "weather-codes.json")
	b, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Printf("failed to read weather code file %s: %v", jsonPath, err)
		return
	}

	if err := json.Unmarshal(b, &weatherCodeDescriptions); err != nil {
		log.Printf("failed to decode weather code descriptions: %v", err)
	}
}

func weatherDescription(code int, isDay int) string {
	key := strconv.Itoa(code)
	if desc, ok := weatherCodeDescriptions[key]; ok {
		if isDay == 1 {
			return desc.Day.Description
		}
		return desc.Night.Description
	}
	return "Unknown conditions"
}

func getLocationFromIpApi() (locationResponse, error) {
	resp, err := http.Get("https://ipapi.co/json/")
	if err != nil {
		return locationResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return locationResponse{}, fmt.Errorf("ipapi.co returned status code %d", resp.StatusCode)
	}

	var location locationResponse
	err = json.NewDecoder(resp.Body).Decode(&location)
	return location, err
}

func getLocationFromIpInfo() (locationResponse, error) {
	resp, err := http.Get("https://ipinfo.io/json")
	if err != nil {
		return locationResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return locationResponse{}, fmt.Errorf("ipinfo.io returned status code %d", resp.StatusCode)
	}

	var ipInfoResponse struct {
		City      string `json:"city"`
		Loc       string `json:"loc"`
		Latitude  float64
		Longitude float64
	}
	err = json.NewDecoder(resp.Body).Decode(&ipInfoResponse)
	if err != nil {
		return locationResponse{}, err
	}

	// Parse coordinates from "latitude,longitude" format
	parts := strings.Split(ipInfoResponse.Loc, ",")
	if len(parts) == 2 {
		lat, errLat := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		lon, errLon := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if errLat == nil && errLon == nil {
			ipInfoResponse.Latitude = lat
			ipInfoResponse.Longitude = lon
		}
	}

	return locationResponse{
		City:      ipInfoResponse.City,
		Latitude:  ipInfoResponse.Latitude,
		Longitude: ipInfoResponse.Longitude,
	}, nil
}

func getLocation() (locationResponse, error) {
	// Try ipap.co first
	location, err := getLocationFromIpApi()
	if err == nil {
		return location, nil
	}

	// Retry with ipinfo.io
	return getLocationFromIpInfo()
}

func getWeather(location location) (weatherResponse, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,is_day,precipitation,weather_code,wind_speed_10m&wind_speed_unit=mph", location.Latitude, location.Longitude)
	resp, err := http.Get(url)
	if err != nil {
		return weatherResponse{}, err
	}
	defer resp.Body.Close()

	var weather weatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weather)
	return weather, err
}

func marshalLocationResponse(r locationResponse) location {
	return location{
		City:      r.City,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
	}
}

func marshalWeatherResponse(r weatherResponse) weather {
	return weather{
		Temperature:   r.Current.Temperature,
		Precipitation: r.Current.Precipitation,
		WeatherCode:   r.Current.WeatherCode,
		WindSpeed:     r.Current.WindSpeed,
		IsDay:         r.Current.IsDay,
	}
}

func NewWeatherWidget() WeatherWidget {
	locationResponse, err := getLocation()
	if err != nil {
		return WeatherWidget{
			hasError: true,
		}
	}

	locationToStore := marshalLocationResponse(locationResponse)

	weatherData, err := getWeather(locationToStore)
	if err != nil {
		return WeatherWidget{
			hasError: true,
		}
	}

	return WeatherWidget{
		location: locationToStore,
		weather:  marshalWeatherResponse(weatherData),
	}
}

func (w WeatherWidget) Init() tea.Cmd {
	return weatherTick(time.Minute)
}

func (w WeatherWidget) Update(msg tea.Msg) (WeatherWidget, tea.Cmd) {
	switch msg.(type) {
	case WeatherTickMsg:
		weatherData, err := getWeather(w.location)
		if err != nil {
			w.hasError = true
			return w, weatherTick(time.Minute)
		}

		w.weather = marshalWeatherResponse(weatherData)
		w.hasError = false

		return w, weatherTick(time.Minute)
	}
	return w, nil
}

func (w WeatherWidget) View() string {
	if w.hasError {
		errorStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.BrightRed)

		return errorStyle.Render("Error fetching weather data")
	}

	cityStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.BrightCyan).
		Width(35).
		PaddingBottom(1)

	tempStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.BrightMagenta)

	conditionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.BrightWhite).
		Align(lipgloss.Center).
		Width(12)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.BrightWhite)

	locationText := w.location.City
	if locationText == "" {
		locationText = "Unknown location"
	}

	locationStr := cityStyle.Render(locationText)
	tempStr := tempStyle.Render(fmt.Sprintf("%.1f°C", w.weather.Temperature))
	conditionStr := conditionStyle.Render(weatherDescription(w.weather.WeatherCode, w.weather.IsDay))
	iconStr := labelStyle.Render(weatherIcon(w.weather.WeatherCode, w.weather.IsDay))
	precipStr := labelStyle.Render(fmt.Sprintf("Precip: %.1f mm", w.weather.Precipitation))
	windStr := labelStyle.Render(fmt.Sprintf("Wind: %.1f mp/h", w.weather.WindSpeed))

	leftColumn := lipgloss.JoinVertical(
		lipgloss.Left,
		locationStr,
		tempStr,
	)

	rightColumn := lipgloss.JoinVertical(
		lipgloss.Left,
		iconStr,
		conditionStr,
	)

	topSection := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftColumn,
		"  ",
		rightColumn,
	)

	bottomSection := lipgloss.JoinHorizontal(lipgloss.Top, windStr, " | ", precipStr)

	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		topSection,
		bottomSection,
	)

	weatherWidgetStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	return weatherWidgetStyle.Render(fullContent)
}

func weatherTick(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return WeatherTickMsg(t)
	})
}
