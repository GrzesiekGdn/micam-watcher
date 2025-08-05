//go:build windows

package platform

import (
	"log"
	"time"

	"github.com/GrzesiekGdn/micam-watcher/internal/common"
	"github.com/GrzesiekGdn/micam-watcher/internal/core"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

func StartService() {
	isInteractive, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Fatalf("Failed to determine session type: %v", err)
	}

	config, err := common.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	core.SetupLogging(config.LogPath)

	if !isInteractive {
		// Run as a Windows Service
		runService("Camera-WS", false)
		return
	}

	// Run in interactive (dev) mode
	log.Println("Running in interactive mode (not as a service)")

	for {
		core.RunMainJob(config)
		time.Sleep(time.Duration(config.Timespan) * time.Millisecond)
	}
}

type myService struct{}

func (m *myService) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue

	config, err := common.LoadConfig()
	if err != nil {
		log.Printf("Error loading configuration: %v", err)
		log.Fatalln("Failed to load configuration. Exiting service.")
	}

	tick := time.Tick(time.Duration(config.Timespan) * time.Millisecond)

	status <- svc.Status{State: svc.StartPending}
	status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		select {
		case <-tick:
			core.RunMainJob(config)
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				status <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				log.Print("Shutting service...!")
				break loop
			case svc.Pause:
				status <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
				log.Printf("Unexpected service control request #%d", c)
			}
		}
	}

	status <- svc.Status{State: svc.StopPending}
	return false, 1
}

func runService(name string, isDebug bool) {
	if isDebug {
		err := debug.Run(name, &myService{})
		if err != nil {
			log.Fatalln("Error running service in debug mode.")
		}
	} else {
		err := svc.Run(name, &myService{})
		if err != nil {
			log.Printf("Error running service: %v", err)
			log.Fatalln("Error running service in Service Control mode.")
		}
	}
}
