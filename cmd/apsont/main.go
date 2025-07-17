package main

import (
	"fmt"
	"log"
	"os"

	"github.com/WATonomous/APSON/internal/config"
)

// APSON main entrypoint
func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs/config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Loaded config: %+v\n", cfg)
	// TODO: Wire up PlantOps polling, state, and notifier
}
