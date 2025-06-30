# 🔧 Build & Cross-Compile Instructions

## 🖥️ Build for Current OS

```bash
go build -o micam-watcher ./cmd/micam-watcher
```

## 🪟 Build for Windows from Linux/macOS

```bash
GOOS=windows GOARCH=amd64 go build -o ./bin/windows/micam-watcher.exe ./cmd/micam-watcher
```

## 🐧 Build for Linux from Windows/macOS

```bash
GOOS=linux GOARCH=amd64 go build -o ./bin/linux/micam-watcher ./cmd/micam-watcher
```

## 🧪 Test Run Locally

Make sure your `config.json` is in the same directory as the binary.

```bash
./micam-watcher
```

or

```cmd
micam-watcher.exe
```
