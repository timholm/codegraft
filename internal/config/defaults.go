package config

import "os"

const (
	DefaultOllamaURL  = "http://localhost:11434"
	DefaultBaseModel  = "codellama:7b"
	DefaultServerPort = 8080
	DefaultLogLevel   = "info"
	DefaultMaxExamples = 5000
	DefaultEvolutionDepth = 3
	DefaultDiversityThreshold = 0.7
	AppName = "codegraft"
)

func defaultDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".codegraft"
	}
	return home + "/.codegraft"
}
