//go:build linux

package platformservices

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/GrzesiekGdn/micam-watcher/internal/common"
)

func IsCameraInUse(config common.Config) (bool, error) {
	matches, err := filepath.Glob(config.Device)
	if err != nil {
		return false, err
	}

	for _, device := range matches {
		cmd := exec.Command("lsof", device)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err == nil {
			if strings.TrimSpace(out.String()) != "" {
				return true, nil
			}
		}
	}

	return false, nil
}

func IsMicrophoneInUse(config common.Config) (bool, error) {
	cmd := exec.Command("pactl", "list", "short", "source-outputs")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return len(output) > 0, nil
}
