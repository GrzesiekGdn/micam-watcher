package core

import (
	"github.com/GrzesiekGdn/micam-watcher/internal/common"
	"github.com/GrzesiekGdn/micam-watcher/internal/platformservices"
	"log"
)

var lastIsCameraInUse = false
var lastIsMicrophoneInUse = false

func RunMainJob(config common.Config) {
	log.Println("Checking for active camera and microphone applications...")
	isCameraInUse, err := platformservices.IsCameraInUse(config)
	if err != nil {
		log.Printf("Error checking active camera applications: %v", err)
		return
	}

	isMicrophoneInUse, err := platformservices.IsMicrophoneInUse(config)
	if err != nil {
		log.Printf("Error checking active microphone applications: %v", err)
		return
	}

	log.Println("Done checking active camera and microphone applications.")

	if isCameraInUse != lastIsCameraInUse || isMicrophoneInUse != lastIsMicrophoneInUse {
		err := SendPost(config.Url, isCameraInUse, isMicrophoneInUse)
		if err != nil {
			log.Printf("Error sending POST request: %v", err)
			return
		}

		log.Println("Sent POST request to Home Assistant.")

		lastIsCameraInUse = isCameraInUse
		lastIsMicrophoneInUse = isMicrophoneInUse
	}
}
