package main

import (
	"RWB_L0/internal/app"
	"log"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
