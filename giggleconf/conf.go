package giggleconf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mangalaman93/giggle/gigglesync"
	"github.com/mangalaman93/giggle/giggleutils"
)

const (
	cAppName          = "giggle"
	cUIPort           = 4444
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

	cIconFile         = "images/giggle.png"
	cSettingsIconFile = "images/settings.png"
	clogIconFile      = "images/log.png"
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

func GetAppName() string {
	return cAppName
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

func GetIconFile() string {
	return cIconFile
}

func GetSettingsIconFile() string {
	return cSettingsIconFile
}

func GetLogIconFile() string {
	return clogIconFile
}

func LoadConfigFile() (*GiggleConfig, error) {
	baseFolder := GetBaseFolder()
	settingsFile := filepath.Join(baseFolder, cSettingsFile)

	_, err := os.Stat(settingsFile)
	if os.IsNotExist(err) {
		emptyConfig, err := json.MarshalIndent(NewGiggleConfig(), "", "\t")
		if err != nil {
			return nil, err
		}

		ioutil.WriteFile(settingsFile, emptyConfig, cSecureFilePerm)
	} else if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(settingsFile)
	if err != nil {
		return nil, err
	}

	var giggleConfig GiggleConfig
	err = json.Unmarshal(data, &giggleConfig)
	if err != nil {
		return nil, err
	}

	return &giggleConfig, nil
}

func GetURLToUI() string {
	return fmt.Sprintf("http://localhost:%s/", cUIPort)
}
