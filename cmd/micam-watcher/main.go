package main

import "github.com/GrzesiekGdn/micam-watcher/internal/platform"

import (
    "flag"
    "fmt"
)

var version = "dev" // will be replaced by ldflags

func main() {
    showVersion := flag.Bool("version", false, "Print version and exit")
    flag.Parse()

    if *showVersion {
        fmt.Println("micam-watcher version", version)
        return
    }
	
	platform.StartService()
}
