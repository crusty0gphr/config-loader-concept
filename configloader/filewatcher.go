package configloader

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

const checkInterval = 3 * time.Second

type fileWatcher struct {
	service string
	path    string
	modTime time.Time
}

func RunFileWatcher(files map[string]string, updates chan<- string) {
	var filesToWatch []*fileWatcher

	for service, path := range files {
		modTime, err := getFileModTime(path)
		if err != nil {
			log.Fatalf("failed to read file: %v", err)
		}

		filesToWatch = append(filesToWatch, &fileWatcher{
			service: service,
			path:    path,
			modTime: modTime,
		})
	}

	// Create a ticker to check the files periodically
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Iterate over each file in the list of files to watch
			for i, file := range filesToWatch {
				// Get the current modification time of the file
				currentModTime, err := getFileModTime(file.path)
				if err != nil {
					log.Printf("Error getting file modification time for %s: %v", file.path, err)
					// Continue to the next file in the list
					continue
				}

				// Check if the file has been modified since the last check
				if currentModTime.After(file.modTime) {
					log.Printf("File %s has been modified at %s", file.path, currentModTime)
					// Update the modification time
					filesToWatch[i].modTime = currentModTime

					// Send the associated service to the updates channel
					updates <- file.service
				}
			}
		}
	}
}

// getFileModTime returns the modification time of the file at the given path
func getFileModTime(filePath string) (time.Time, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return fileInfo.ModTime(), nil
}

func FileChangeHandler(service, path string, kv nats.KeyValue) error {
	cfg, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read file: %w", err)
	}

	_, err = kv.Put(service, cfg)
	if err != nil {
		return fmt.Errorf("nats: unable to put config into kv: %w", err)
	}

	log.Printf("successful config update for %s", service)
	return nil
}
