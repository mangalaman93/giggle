package giggleutils

import (
	"github.com/axet/desktop/go"
)

func GetAppFolder() string {
	return desktop.GetAppDataFolder()
}
