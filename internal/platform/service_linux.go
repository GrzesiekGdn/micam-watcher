//go:build linux
package platform

import (
    "log"

    "github.com/GrzesiekGdn/micam-watcher/internal/core"
)

func StartService() {
	config, err := core.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)        
	}

    core.RunMainJob(config)
}