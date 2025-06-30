package core

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type StatusPayload struct {
    CameraOn     bool `json:"cameraOn"`
    MicrophoneOn bool `json:"microphoneOn"`
}

func SendPost(url string, cameraOn, microphoneOn bool) error {
    payload := StatusPayload{CameraOn: cameraOn, MicrophoneOn: microphoneOn}
    body, _ := json.Marshal(payload)

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    _, err := http.DefaultClient.Do(req)
    return err
}
