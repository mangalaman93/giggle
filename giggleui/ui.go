package giggleui

import (
	"github.com/sqweek/dialog"
)

type uiConfig struct {
	exposeExternal bool
}

func Error(message string) {
	dialog.Message(message).Error()
}
