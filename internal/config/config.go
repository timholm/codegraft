package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds all runtime configuration for codegraft.
type Config struct {
	DataDir            string  `json:"data_dir"`
	OllamaURL          string  `json:"ollama_url"`
	BaseModel          string  `json:"base_model"`
	ServerPort         int     `json:"server_port"`
	LogLevel           string  `json:"log_level"`
	MaxExamples        int     `json:"max_examples"`
	EvolutionDepth     int     `json:"evolution_depth"`
	DiversityThreshold float64 `json:"diversity_threshold"`
}

// Default returns a Config populated with default values.
func Default() *Config {
	return &Config{
		DataDir:            defaultDataDir(),
		OllamaURL:          DefaultOllamaURL,
		BaseModel:          DefaultBaseModel,
		ServerPort:         DefaultServerPort,
		LogLevel:           DefaultLogLevel,
		MaxExamples:        DefaultMaxExamples,
		EvolutionDepth:     DefaultEvolutionDepth,
		DiversityThreshold: DefaultDiversityThreshold,
	}
}

// Load reads a config file at path. Missing fields fall back to defaults.
func Load(path string) (*Config, error) {
	cfg := Default()
	if path == "" {
		path = filepath.Join(cfg.DataDir, "config.json")
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	return cfg, nil
}

// Save writes the config to path, creating directories as needed.
func (c *Config) Save(path string) error {
	if path == "" {
		path = filepath.Join(c.DataDir, "config.json")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// LibraryDir returns the directory used for a named library's data.
func (c *Config) LibraryDir(library string) string {
	return filepath.Join(c.DataDir, "libraries", library)
}

// ExamplesFile returns the path to the training examples file for a library.
func (c *Config) ExamplesFile(library string) string {
	return filepath.Join(c.LibraryDir(library), "examples.jsonl")
}

// ModelfileDir returns the directory for generated Modelfiles.
func (c *Config) ModelfileDir(library string) string {
	return filepath.Join(c.LibraryDir(library), "modelfiles")
}
