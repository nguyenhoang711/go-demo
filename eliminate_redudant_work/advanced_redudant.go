package eliminateredudantwork

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type WeatherService struct {
	requestGroup singleflight.Group
	cache        sync.Map
}


func (w *WeatherService) fetchWeatherData(city string) (string, error) {
	resp, err := http.Get("http://example.com/weather/" + city)
	if err != nil {
		return "error", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (w *WeatherService) GetWeather(city string) (string, error) {
	// check data in cache or not
	if data, ok := w.cache.Load(city); ok {
		return data.(string), nil
	}
	// if not, using singleflight so only one request which same city can fetch
	data, err, _ := w.requestGroup.Do(city, func() (interface{}, error) {
		result, err := w.fetchWeatherData(city)
		if err == nil {
			w.cache.Store(city, result)
		}
		return result, err
	})
	if err != nil {
		return "error when fetching data", err
	}
	return data.(string), nil
}

func AdvancedFetchWeather() {
	service := &WeatherService{}

	for i := 0; i< 10; i++ {
		go func(i int) {
			weather, err := service.GetWeather("NewYork")
			if err == nil {
				fmt.Printf("Goroutine %d got weather data: %v\n", i, weather)
			} else {
				fmt.Printf("Goroutine %d encoutered an error :: %s\n", i, err)
			}
		}(i)
	}
	time.Sleep(5 * time.Second)
}