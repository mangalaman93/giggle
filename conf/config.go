package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ConfigHolder holds the configuration and provides safe access to it.
type ConfigHolder struct {
	mu     sync.RWMutex
	config Config
}

// Config stores all the configuration.
type Config struct {
	Sync   []SyncConfig           `json:"sync"`
	Auth   map[string]*AuthMethod `json:"auth"`
	Period duration               `json:"period"`
}

// SyncConfig stores configuration for completing sync.
type SyncConfig struct {
	Name   string `json:"name"`
	From   Repo   `json:"from"`
	ToList []Repo `json:"to"`
}

// Repo is one of the repos (github or overleaf).
type Repo struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	URLToRepo string `json:"url"`
	AuthToUse string `json:"auth"`
}

// New reads the configuration file from disk.
func New() (*ConfigHolder, error) {
	configFile := filepath.Join(baseFolder(), cConfigFile)
	return newFrom(configFile)
}

func newFrom(configFile string) (*ConfigHolder, error) {
	// set secure permission for the file
	if err := os.Chmod(configFile, cSecureFilePerm); err != nil {
		return nil, fmt.Errorf("error modifying perm for conf file :: %w", err)
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading conf file :: %w", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling conf file :: %w", err)
	}

	return &ConfigHolder{config: config}, nil
}

// Get returns a copy of the config
func (h *ConfigHolder) Get() Config {
	h.mu.Lock()
	defer h.mu.Unlock()

	authMap := make(map[string]*AuthMethod, len(h.config.Auth))
	for authName, authMethod := range h.config.Auth {
		var sb strings.Builder
		_, _ = sb.WriteString(authName)
		authMap[sb.String()] = authMethod.copy()
	}

	return Config{
		Period: h.config.Period,
		Sync:   h.config.Sync,
		Auth:   authMap,
	}
}

// GetPeriod returns the periodicity of sync.
func (h *ConfigHolder) GetPeriod() time.Duration {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.config.Period.Duration
}
