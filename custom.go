package main

import (
	"encoding/json"
	"log"
	"os"
)

type customButton struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type customDevice struct {
	ID      string         `json:"id"`
	Buttons []customButton `json:"buttons"`
}

var customDevices = map[string]customDevice{}

func init() {
	s, ok := os.LookupEnv("CUSTOM_DEVICES")
	if !ok {
		return
	}
	var devs []customDevice
	if err := json.Unmarshal([]byte(s), &devs); err != nil {
		log.Printf("Failed to parse CUSTOM_DEVICES: %v", err)
		return
	}
	for _, d := range devs {
		for i := range d.Buttons {
			if d.Buttons[i].Icon == "" {
				d.Buttons[i].Icon = "circle"
			}
		}
		customDevices[d.ID] = d
	}
}
