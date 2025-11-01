package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/justdry/eventbus"
)

type PersonData struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func Hi(ctx context.Context, b []byte) error {
	var data PersonData
	if err := json.Unmarshal(b, &data); err != nil {
		return eventbus.NewError(fmt.Errorf("the json structure should be {\"firstName\": string, \"lastName\": string}: %w", err))
	}

	log.Printf("Hi %s %s,", data.FirstName, data.LastName)
	return nil
}

func Greeting(ctx context.Context, _ []byte) error {
	log.Println("How are you doing?")

	return nil
}
