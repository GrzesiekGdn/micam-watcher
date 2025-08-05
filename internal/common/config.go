package common

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Url      string `json:"url"`
	Timespan int    `json:"timespan"`
	UserSid  string `json:"userSid,omitempty"`
	Device   string `json:"device,omitempty"`
	LogPath  string `json:"logPath,omitempty"`
}

func LoadConfig() (Config, error) {
	var config Config
	exePath, err := os.Executable()
	if err != nil {
		return config, err
	}
	configPath := filepath.Join(filepath.Dir(exePath), "config.json")

	fileData, err := os.ReadFile(configPath)

	if err != nil {
		return config, err
	}

	return config, json.Unmarshal(fileData, &config)
}
