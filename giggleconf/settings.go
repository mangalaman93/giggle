package giggleconf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mangalaman93/giggle/giggleutils"
  "github.com/mangalaman93/giggle/gigglesync"
)

const (
	cPasswordLength = 24
	cSettingsFile   = ".gsettings"
	cSecureFilePerm = 600
)

type GiggleSettings struct {
	OnlySyncData bool   `json:"only_sync_data"`
	RConfSecret  string `json:"rconf_secret"`
  Repos []*gigglesync.Repo `json:"repos"`
}

func NewGiggleSettings() *GiggleSettings {
	return &GiggleSettings{
		OnlySyncData: false,
		RConfSecret:  giggleutils.RandomString(cPasswordLength),
	}
}

func ReadGiggleSettings() (*GiggleSettings, error) {
	baseFolder := giggleutils.GetAppFolder()
	settingsFile := filepath.Join(baseFolder, cSettingsFile)

	_, err := os.Stat(settingsFile)
	if os.IsNotExist(err) {
		emptyConfig, err := json.Marshal(NewGiggleSettings())
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

	var giggleSettings GiggleSettings
	err = json.Unmarshal(data, &giggleSettings)
	if err != nil {
		return nil, err
	}

	return giggleSettings, nil
}
