# micam-watcher

**micam-watcher** is a cross-platform Go service that monitors camera and microphone activity and sends a POST request to a configured endpoint — for example, to [Home Assistant](https://www.home-assistant.io/) for automations such as turning on or off a lamp.

---

## ✨ Features

- Cross-platform (Linux & Windows)
- Detects camera and microphone activity
- Sends POST requests to any HTTP endpoint
- Configurable polling interval
- Runs as a background service

---

## ⚙️ How It Works

### Common Code
- Reads configuration
- Main monitoring loop
- Logging

### Linux
- **Camera detection:** via `lsof` on `/dev/video*`
- **Microphone detection:** via `pactl` for active sources  
  > Note: `pactl` only lists microphones for the current user.  
  The app must run as a *systemd user service*.

### Windows
- **Detection:** via Windows Registry keys for camera & microphone usage
- Runs as a standard Windows Service
- Requires the **user SID** of the monitored account in `config.json`:
  ```powershell
  whoami /user
  ```

## 📊 Performance

- Minimal CPU usage for most configurations  
- On Linux, CPU usage depends on polling frequency (`timestamp` in config)  
- Lower `timestamp` values mean faster detection but slightly higher CPU usage

---

## 📝 Configuration

Example `config.json`:
```json
{
  "endpoint": "http://homeassistant.local/api/webhook/device-status",
  "timestamp": 1000,
  "userSID": "S-1-5-21-1234567890-123456789-1234567890-1001",
  "device":  "/dev/video*"
}
```

Fields:
- endpoint - URL to send POST requests to
- timestamp - interval in milliseconds between checks
- userSID - (Windows only) SID of the monitored user
- device - (Linux only) device to monitor camera

## 🚀 Quick Start

### Linux
```bash
# Copy binary & config to ~/.local/bin or /usr/local/bin
cp micam-watcher ~/.local/bin/
cp config.json ~/.local/bin/

# Install and start systemd user service
systemctl --user enable --now micam-watcher

```

### Windows (PowerShell)
```powershell
# Copy binary & config to Program Files
mkdir "C:\Program Files\micam-watcher"
copy micam-watcher.exe "C:\Program Files\micam-watcher\"
copy config.json "C:\Program Files\micam-watcher\"

# Install and start service
sc.exe create Micam-watcher binPath= "C:\Program Files\micam-watcher\micam-watcher.exe"
sc.exe start Micam-watcher
```

## 📦 Architecture Overview

```
+-----------------+       Camera/Mic activity       +-------------------+
| Camera / Mic    | ------------------------------> | micam-watcher     |
| (hardware)      |                                 | (Go service)      |
+-----------------+                                 +---------+---------+
                                                              |
                                                              v
                                                   HTTP POST to endpoint
                                                              |
                                                              v
                                            +-----------------+-----------------+
                                            | Automation platform (e.g. HA)     |
                                            | Turns on/off lamp or other device |
                                            +-----------------------------------+

```

## 🔧 Build & Installation

See BUILD.md for:
- Building from source
- Cross-compilation (Linux ↔ Windows)
- Full service installation commands
