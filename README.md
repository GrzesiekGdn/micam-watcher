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
  "device":  "/dev/video*",
  "log-path": "/var/log/micam-watcher.log",
}
```

Fields:
- endpoint - URL to send POST requests to
- timestamp - interval in milliseconds between checks
- userSID - (Windows only) SID of the monitored user
- device - (Linux only) device to monitor camera
- log-path - optional path to log file.

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

## 🔩 Alternative: hardware LED sensor

`micam-watcher` has to run **on the machine you are monitoring**. On a corporate or otherwise
locked-down computer that is often impossible — no admin rights, no installing services, no
touching registry keys or systemd units.

For that case there is a companion project:
**[brio105-led-sensor](https://github.com/GrzesiekGdn/brio105-led-sensor)** — a 3D-printable
clip that holds a phototransistor over a Logitech Brio 105's activity LED, so an Arduino can
tell whether the camera is streaming without any software on the host at all. Slower to set
up, camera-specific, and it cannot see microphone activity — but it needs no access to the
monitored machine whatsoever.

Use `micam-watcher` when you can install software; use the hardware sensor when you cannot.

---

## 🔧 Build & Installation

See BUILD.md for:
- Building from source
- Cross-compilation (Linux ↔ Windows)
- Full service installation commands

## ⚠️ Antivirus Notes

`micam-watcher` is a small monitoring tool written in Go.  
Since it accesses **system resources** (camera, microphone, registry on Windows, `lsof/pactl` on Linux), Windows Defender or other antivirus tools may occasionally flag it or require you to add an exception.

This is a **false positive** — the program does **not** contain malicious code, does **not** collect or send data anywhere except the configured endpoint (e.g., Home Assistant), and is fully open source, so you can audit the code yourself.

If you encounter a warning:
- Verify you are using an official release from [this repository](https://github.com/GrzesiekGdn/micam-watcher).
- Optionally add the binary path to your Defender exclusions.