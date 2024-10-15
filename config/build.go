package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type BuildInfo struct {
	// Version is the application version
	Version string `json:"version"`
	// BuildTime is the time when the application was built
	BuildTime string `json:"buildTime"`
}

const (
	buildFile         = "build.properties"
	versionProperty   = "version"
	buildTimeProperty = "buildTime"
)

var buildInfo BuildInfo

// initializeBuildData logs the application version and build time
func initializeBuildData() {
	file, err := os.Open(buildFile)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, "=")
		if len(split) == 2 {
			key := split[0]
			value := split[1]
			switch key {
			case versionProperty:
				buildInfo.Version = value
			case buildTimeProperty:
				buildInfo.BuildTime = value
			}
		}
	}
	log.Printf("Build version: %s", buildInfo.Version)
	log.Printf("Build time: %s", buildInfo.BuildTime)
}

// GetBuildInfo returns the application version and build time
func GetBuildInfo() BuildInfo {
	return buildInfo
}
