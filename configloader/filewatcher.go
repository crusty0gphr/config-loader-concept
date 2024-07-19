package configloader

import (
	"fmt"
	"log"
	"os"
	"time"
)

const checkInterval = 5 * time.Second

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
			for i, file := range filesToWatch {
				currentModTime, err := getFileModTime(file.path)
				if err != nil {
					log.Printf("Error getting file modification time for %s: %v", file.path, err)
					continue
				}

				if currentModTime.After(file.modTime) {
					fmt.Printf("File %s has been modified at %s\n", file.path, currentModTime)
					filesToWatch[i].modTime = currentModTime

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
