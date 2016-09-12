package giggleconf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mangalaman93/giggle/gigglesync"
	"github.com/mangalaman93/giggle/giggleutils"
)

const (
	cPasswordLength   = 24
	cAppFolder        = ".giggle"
	cLogFolder        = "log"
	cSettingsFile     = ".gsettings"
	cLogFile          = "giggle.log"
	cSecureFilePerm   = 600
	cFilePerm         = 640
	cLogFileMaxSize   = 50 // MB
	cLogMaxNumBackups = 5
	cLogFileMaxAge    = 30 // days
)

type GiggleConfig struct {
	OnlySyncData bool                         `json:"only_sync_data"`
	RConfSecret  string                       `json:"rconf_secret"`
	Repositories []*gigglesync.GiggleSyncRepo `json:"repos"`
}

func NewGiggleConfig() *GiggleConfig {
	return &GiggleConfig{
		OnlySyncData: false,
		RConfSecret:  giggleutils.RandomString(cPasswordLength),
	}
}

func GetBaseFolder() string {
	rootFolder := giggleutils.GetAppFolder()
	return filepath.Join(rootFolder, cAppFolder)
}

func GetLogFolder() string {
	baseFolder := GetBaseFolder()
	return filepath.Join(baseFolder, cLogFolder)
}

func GetLogFilePath() string {
	logFolder := GetLogFolder()
	return filepath.Join(logFolder, cLogFile)
}

func GetFilePerm() int {
	return cFilePerm
}

func GetLogFileMaxSize() int {
	return cLogFileMaxSize
}

func GetLogMaxNumBackups() int {
	return cLogMaxNumBackups
}

func GetLogFileMaxAge() int {
	return cLogFileMaxAge
}

func LoadConfigFile() (*GiggleConfig, error) {
	baseFolder := GetBaseFolder()
	settingsFile := filepath.Join(baseFolder, cSettingsFile)

	_, err := os.Stat(settingsFile)
	if os.IsNotExist(err) {
		emptyConfig, err := json.Marshal(NewGiggleConfig())
		if err != nil {
			return nil, err
		}

		ioutil.WriteFile(settingsFile, emptyConfig, cSecureFilePerm)
	}

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(settingsFile)
	if err != nil {
		return nil, err
	}

	var giggleConfig GiggleConfig
	err = json.Unmarshal(data, &giggleConfig)
	if err != nil {
		return nil, err
	}

	return giggleConfig, nil
}
