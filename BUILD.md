# 🔧 Build & Install Guide

This guide explains how to build **micam-watcher** from source, cross-compile for other platforms, and install it as a background service.

---

## 📋 Prerequisites

- **Go**: version 1.20+ (check with `go version`)
- **Git**: to clone the repository
- **Platform-specific tools**:
  - **Linux**: `lsof`, `pactl` (PulseAudio)
  - **Windows**: PowerShell (for `whoami /user`)
- A valid `config.json` file (see README for format)

---

## 🖥️ Build for Current OS

```bash
go build -o ./bin/micam-watcher ./cmd/micam-watcher
```
The binary will be placed in the `./bin/` directory.

## 🪟 Build for Windows from Linux/macOS

```bash
GOOS=windows GOARCH=amd64 go build -o ./bin/windows/micam-watcher.exe ./cmd/micam-watcher
```

## 🐧 Build for Linux from Windows/macOS

```bash
GOOS=linux GOARCH=amd64 go build -o ./bin/linux/micam-watcher ./cmd/micam-watcher
```

## 🚀 Install as a Service

### Windows

```powershell
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/windows/micam-watcher.exe ./cmd/micam-watcher
mkdir "C:\Program Files\micam-watcher"
# sc.exe stop Micam-watcher   # optional, if already installed
copy ./bin/windows/micam-watcher.exe "C:\Program Files\micam-watcher\"
copy ./bin/config.json "C:\Program Files\micam-watcher\"
sc.exe create Micam-watcher binPath="C:\Program Files\micam-watcher\micam-watcher.exe"
sc.exe start Micam-watcher
```
Tip: Use `whoami /user` in PowerShell to get your user SID for the configuration file.

### Linux (systemd --user)

```bash
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/linux/micam-watcher ./cmd/micam-watcher
mkdir -p /usr/local/bin/micam-watcher/
# systemctl --user stop --now micam-watcher   # optional, if already installed
sudo cp ./bin/linux/micam-watcher /usr/local/bin/micam-watcher/
sudo cp ./bin/config.json /usr/local/bin/micam-watcher/
mkdir -p ~/.config/systemd/user/
sudo cp ./micam-watcher.service ~/.config/systemd/user/
systemctl --user enable --now micam-watcher
```
#### Example `micam-watcher.service` file:

```ini
[Unit]
Description=Micam Watcher
After=network.target

[Service]
ExecStart=/usr/local/bin/micam-watcher/micam-watcher
Restart=on-failure

[Install]
WantedBy=default.target
```

## 🧪 Testing Without Service
You can run the binary directly from a terminal to verify it works before installing as a service.

### Linux
```bash
./bin/linux/micam-watcher
```

### Windows
```powershell
.\bin\windows\micam-watcher.exe
```

## 📦 Cross-Compilation Notes

- Cross-compiling from Linux to Windows works without Wine, but you cannot test service installation until you run the binary on a real or virtual Windows machine.

- When cross-compiling from Windows to Linux, you’ll need to copy the binary to a Linux environment for testing.
