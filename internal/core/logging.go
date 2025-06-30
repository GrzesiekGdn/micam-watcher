package core

import (
    "log"
    "os"
)

func SetupLogging(logPath string) {
    if logPath == "" {
        return
    }

    file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Printf("Failed to open log file: %v", err)
        return
    }

    log.SetOutput(file)
}
