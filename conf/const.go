package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
)

const (
	cAppName          = "giggle"
	cUIPort           = 4444
	cAppFolder        = ".giggle"
	cLogFolder        = "log"
	cReposFolder      = "repos"
	cConfigFile       = "config.json"
	cLogFile          = "giggle.log"
	cSecureFilePerm   = 0600
	cDirPerm          = 0700
	cLogFileMaxSize   = 50 // MB
	cLogMaxNumBackups = 5
	cLogFileMaxAge    = 30 // days

	cIconFile         = "images/giggle.png"
	cSettingsIconFile = "images/settings.png"
	clogIconFile      = "images/log.png"
)

// AppName returns the name of the app.
func AppName() string {
	return cAppName
}

// DirPerm returns the permissions that should be used for creating a directory.
func DirPerm() os.FileMode {
	return cDirPerm
}

// baseFolder returns the directory where all the data for the app is stored.
func baseFolder() string {
	rootFolder := configdir.LocalConfig()
	return filepath.Join(rootFolder, cAppFolder)
}

// LogFolder returns the directory where logs are stored.
func LogFolder() string {
	return filepath.Join(baseFolder(), cLogFolder)
}

// LogFilePath returns the path to log file.
func LogFilePath() string {
	return filepath.Join(LogFolder(), cLogFile)
}

// reposFolder returns the path to directory where all the repos are stored.
func reposFolder() string {
	return filepath.Join(baseFolder(), cReposFolder)
}

// GetSyncTarget returns the path on disk where a given sync is cloned/stored.
func GetSyncTarget(syncName string) string {
	return filepath.Join(reposFolder(), syncName)
}

// LogFileMaxSize returns the max allowed size of a log file in MB.
func LogFileMaxSize() int {
	return cLogFileMaxSize
}

// LogMaxNumBackups returns the number of maximum log files that needs to be kept.
func LogMaxNumBackups() int {
	return cLogMaxNumBackups
}

// LogFileMaxAge returns the oldest log file that needs to be kept.
func LogFileMaxAge() int {
	return cLogFileMaxAge
}

// IconFile returns the path of the app icon file in bindata.
func IconFile() string {
	return cIconFile
}

// SettingsIconFile returns the path of the settings icon file in bindata.
func SettingsIconFile() string {
	return cSettingsIconFile
}

// LogIconFile returns the path of the log icon file in bindata.
func LogIconFile() string {
	return clogIconFile
}

// URLToUI returns a url to Giggle UI.
func URLToUI() string {
	return fmt.Sprintf("http://localhost:%d/", cUIPort)
}
