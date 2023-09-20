package main

import (
	"log"
	"tm-backend-trainee-impl-clean-template/config"
	"tm-backend-trainee-impl-clean-template/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	// Run
	app.Run(cfg)
}
