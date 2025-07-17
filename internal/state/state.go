package state

import (
	"bufio"
	"os"
)

const stateFile = "notified.txt"

// LoadNotified loads the set of already-notified announcement links from the state file.
func LoadNotified() (map[string]struct{}, error) {
	f, err := os.Open(stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]struct{}{}, nil // No state file yet
		}
		return nil, err
	}
	defer f.Close()

	notified := make(map[string]struct{})
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			notified[line] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return notified, nil
}

// SaveNotified appends a new link to the state file if not already present.
func SaveNotified(link string) error {
	// Check if already present
	notified, err := LoadNotified()
	if err != nil {
		return err
	}
	if _, exists := notified[link]; exists {
		return nil // Already present
	}
	f, err := os.OpenFile(stateFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(link + "\n")
	return err
}
