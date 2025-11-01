package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/justdry/eventbus"
)

type WeatherData struct {
	Current string `json:"weather"`
}

func Weather(ctx context.Context, b []byte) error {
	var weather WeatherData
	if err := json.Unmarshal(b, &weather); err != nil {
		return eventbus.NewError(fmt.Errorf("the json structure should be {\"weather\": string}: %w", err))
	}

	log.Printf("Oh, it's a %s day, don't you like it?!", weather.Current)
	return nil
}
