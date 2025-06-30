package core

import (
    "encoding/json"
    "os"
)

type Config struct {
    Url      string `json:"url"`
    Timespan int    `json:"timespan"`
    UserSid  string `json:"userSid,omitempty"`
    LogPath  string `json:"logPath,omitempty"`
}

func LoadConfig(path string) (Config, error) {
    var config Config
    data, err := os.ReadFile(path)
    if err != nil {
        return config, err
    }
    err = json.Unmarshal(data, &config)
    return config, err
}
