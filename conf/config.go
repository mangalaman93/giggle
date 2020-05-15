package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// ConfigHolder holds the configuration and provides safe access to it.
type ConfigHolder struct {
	mu     sync.RWMutex
	config Config
}

// Config stores all the configuration.
type Config struct {
	Sync []SyncConfig              `json:"sync"`
	Auth map[string]http.BasicAuth `json:"auth"`
}

// SyncConfig stores configuration for completing sync.
type SyncConfig struct {
	Name string `json:"name"`
	From Repo   `json:"from"`
	To   []Repo `json:"to"`
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

	// set secure permission for the file
	if err := os.Chmod(configFile, cSecureFilePerm); err != nil {
		return nil, fmt.Errorf("error in modifying perm for conf file :: %w", err)
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("error in reading conf file :: %w", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error in unmarshalling config from conf file :: %w", err)
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid config :: %w", err)
	}

	return &ConfigHolder{config: config}, nil
}

func (c *Config) validate() error {
	for _, s := range c.Sync {
		if _, ok := c.Auth[s.From.AuthToUse]; !ok {
			return fmt.Errorf("auth [%v] not found", s.From.AuthToUse)
		}

		for _, t := range s.To {
			if _, ok := c.Auth[t.AuthToUse]; !ok {
				return fmt.Errorf("auth [%v] not found", t.AuthToUse)
			}
		}
	}

	return nil
}

// Get returns a copy of the config
func (h *ConfigHolder) Get() Config {
	h.mu.Lock()
	defer h.mu.Unlock()

	authMap := make(map[string]http.BasicAuth, len(h.config.Auth))
	for k, v := range h.config.Auth {
		authMap[k] = v
	}

	return Config{
		Sync: h.config.Sync,
		Auth: authMap,
	}
}
