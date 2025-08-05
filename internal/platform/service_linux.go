//go:build linux

package platform

import (
	"github.com/GrzesiekGdn/micam-watcher/internal/common"
	"github.com/GrzesiekGdn/micam-watcher/internal/core"
	"log"
	"time"
)

func StartService() {
	config, err := common.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	for {
		core.RunMainJob(config)
		time.Sleep(time.Duration(config.Timespan) * time.Millisecond)
	}
}
