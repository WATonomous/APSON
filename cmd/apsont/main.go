package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/WATonomous/APSON/internal/config"
	"github.com/WATonomous/APSON/internal/plantops"
	"github.com/WATonomous/APSON/internal/state"
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

	notified, err := state.LoadNotified()
	if err != nil {
		log.Fatalf("Failed to load notified state: %v", err)
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
			newCount := 0
			for _, a := range announcements {
				if _, already := notified[a.Link]; already {
					continue
				}
				// TODO: Send notification for new announcement (email, etc.)
				fmt.Printf("[NEW] Would notify: %s (%s)\n", a.Title, a.Link)
				if err := state.SaveNotified(a.Link); err != nil {
					log.Printf("Failed to save notified link: %v", err)
				} else {
					notified[a.Link] = struct{}{}
				}
				newCount++
			}
			if newCount == 0 {
				fmt.Println("No new relevant announcements.")
			}
		}
		fmt.Printf("Sleeping for %v...\n", interval)
		time.Sleep(interval)
	}
}
