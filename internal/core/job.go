package core

import (
	"log"
	"github.com/GrzesiekGdn/micam-watcher/internal/platformservices"
)

var lastNumberOfCameras = -1
var lastNumberOfMicrophones = -1

func RunMainJob(config Config) {
	log.Println("Checking for active camera and microphone applications...")
	numberOfCameras, err := platformservices.GetActiveCameraApplicationsCount(config.UserSid)
	if err != nil {
		log.Printf("Error checking active camera applications: %v", err)
		return
	}

	numberOfMicrophones, err := platformservices.GetActiveMicrophoneApplicationsCount(config.UserSid)
	if err != nil {
		log.Printf("Error checking active microphone applications: %v", err)
		return
	}

	log.Println("Done checking active camera and microphone applications.")

	if numberOfCameras != lastNumberOfCameras || numberOfMicrophones != lastNumberOfMicrophones {
		cameraOn := numberOfCameras > 0
		microphoneOn := numberOfMicrophones > 0

		err := SendPost(config.Url, cameraOn, microphoneOn)
		if err != nil {
			log.Printf("Error sending POST request: %v", err)
			return
		}

		log.Println("Sent POST request to Home Assistant.")

		lastNumberOfCameras = numberOfCameras
		lastNumberOfMicrophones = numberOfMicrophones
	}
}
