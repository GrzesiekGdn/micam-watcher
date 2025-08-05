//go:build windows

package platformservices

import (
	"fmt"
	"github.com/GrzesiekGdn/micam-watcher/internal/common"
	"golang.org/x/sys/windows/registry"
)

const (
	CameraKeyPath     = `SOFTWARE\Microsoft\Windows\CurrentVersion\CapabilityAccessManager\ConsentStore\webcam`
	MicrophoneKeyPath = `SOFTWARE\Microsoft\Windows\CurrentVersion\CapabilityAccessManager\ConsentStore\microphone`
	lastUsedTimeStop  = "LastUsedTimeStop"
)

func IsCameraInUse(config common.Config) (bool, error) {
	userSid := config.UserSid
	activeApps, err := _getActiveApplications(userSid, CameraKeyPath)
	if err != nil {
		return false, err
	}

	if len(activeApps) > 0 {
		fmt.Println("Camera is in use by the following apps:")
		for _, app := range activeApps {
			fmt.Printf("  %s\n", app)
		}
	} else {
		fmt.Println("Camera is not in use.")
	}

	return len(activeApps) > 0, err
}

func IsMicrophoneInUse(config common.Config) (bool, error) {
	userSid := config.UserSid
	activeApps, err := _getActiveApplications(userSid, MicrophoneKeyPath)
	if err != nil {
		return false, err
	}

	if len(activeApps) > 0 {
		fmt.Println("Microphone is in use by the following apps:")
		for _, app := range activeApps {
			fmt.Printf("  %s\n", app)
		}
	} else {
		fmt.Println("Microphone is not in use.")
	}

	return len(activeApps) > 0, err
}

func _getActiveApplications(userSid string, keyPath string) ([]string, error) {
	var activeApps []string

	// Open base key
	key, err := getRegistryKey(registry.USERS, userSid, keyPath)
	if err != nil {
		return activeApps, err
	}
	defer key.Close()

	// Get subkeys
	subKeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return activeApps, err
	}

	for _, subKeyName := range subKeys {
		fullSubKeyPath := keyPath + `\` + subKeyName

		if subKeyName == "NonPackaged" {
			npKey, err := getRegistryKey(registry.USERS, userSid, fullSubKeyPath)
			if err != nil {
				continue
			}
			defer npKey.Close()

			npSubKeys, err := npKey.ReadSubKeyNames(-1)
			if err != nil {
				continue
			}

			activeApps = _getApplicationsFromSubKeys(userSid, npSubKeys, fullSubKeyPath, activeApps)
		} else {
			subKey, err := getRegistryKey(registry.USERS, userSid, fullSubKeyPath)
			if err != nil {
				continue
			}
			defer subKey.Close()

			timestamp, err := _getSubKeyTimestamp(subKey, lastUsedTimeStop)
			if err == nil && timestamp == 0 {
				activeApps = append(activeApps, subKeyName)
			}
		}
	}

	return activeApps, nil
}

func _getApplicationsFromSubKeys(userSid string, npSubKeys []string, fullSubKeyPath string, activeApps []string) []string {
	for _, npSubKey := range npSubKeys {
		fullNpSubKeyPath := fullSubKeyPath + `\` + npSubKey
		subKey, err := getRegistryKey(registry.USERS, userSid, fullNpSubKeyPath)
		if err != nil {
			continue
		}
		defer subKey.Close()

		timestamp, err := _getSubKeyTimestamp(subKey, lastUsedTimeStop)
		if err == nil && timestamp == 0 {
			activeApps = append(activeApps, npSubKey)
		}
	}
	return activeApps
}

func _getSubKeyTimestamp(k registry.Key, name string) (int64, error) {
	val, _, err := k.GetIntegerValue(name)
	if err != nil {
		return -1, err
	}
	return int64(val), nil
}

func getRegistryKey(baseKey registry.Key, sid, subPath string) (registry.Key, error) {
	fullPath := sid + `\` + subPath
	return registry.OpenKey(baseKey, fullPath, registry.READ)
}
