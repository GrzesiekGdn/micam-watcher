//go:build linux
package platformservices

import (	
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"	
)

func GetActiveCameraApplicationsCount(userSid string) (int, error) {		
	return 0, nil
}

func GetActiveMicrophoneApplicationsCount(userSid string) (int, error) {		
	return 0, nil
}

func isCameraInUse(device string) (bool, error) {
	// Check each /proc/[pid]/fd/* for a symlink pointing to /dev/video0
	procPath := "/proc"
	inUse := false

	filepath.WalkDir(procPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && strings.HasPrefix(d.Name(), "fd") {
			return filepath.SkipDir
		}
		if !d.IsDir() || !isNumeric(d.Name()) {
			return nil
		}

		fdPath := filepath.Join(procPath, d.Name(), "fd")
		entries, err := os.ReadDir(fdPath)
		if err != nil {
			return nil // may not have access
		}

		for _, entry := range entries {
			link, err := os.Readlink(filepath.Join(fdPath, entry.Name()))
			if err == nil && link == device {
				inUse = true
				return fs.SkipAll // found, no need to continue
			}
		}
		return nil
	})

	return inUse, nil
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func isMicrophoneInUse() bool {
	cmd := exec.Command("pactl", "list", "source-outputs")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "Source Output #")
}