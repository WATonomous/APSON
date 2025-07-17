package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/WATonomous/APSON/internal/config"
	"github.com/WATonomous/APSON/internal/plantops"
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

	fmt.Println("APSON service starting...")
	fmt.Printf("Monitoring buildings: %v\n", cfg.Buildings)

	interval := time.Duration(cfg.PollingIntervalMinutes) * time.Minute
	for {
		fmt.Println("Polling PlantOps...")
		announcements, err := plantops.FetchAndParse(cfg.Buildings)
		if err != nil {
			log.Printf("Error fetching PlantOps: %v", err)
		} else {
			// TODO: Filter out already-notified announcements using state
			// TODO: Send notifications for new announcements
			fmt.Printf("Found %d relevant announcements:\n", len(announcements))
			for _, a := range announcements {
				fmt.Printf("- %s (%s)\n", a.Title, a.Link)
			}
		}
		fmt.Printf("Sleeping for %v...\n", interval)
		time.Sleep(interval)
	}
}
