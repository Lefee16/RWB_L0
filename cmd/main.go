package main

import (
	"log"
	"os"

	"RWB_L0/internal/app"
)

func main() {
	// Создаём приложение
	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
		os.Exit(1)
	}

	// Запускаем приложение
	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
		os.Exit(1)
	}
}
