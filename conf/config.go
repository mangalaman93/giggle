package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

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

// ReadConfig reads the config from file on disk.
func ReadConfig(configFile string) (*Config, error) {
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

	return &config, nil
}
