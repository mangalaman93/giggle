package giggleui

import (
	"github.com/sqweek/dialog"
)

func Error(message string) {
	dialog.Message(message).Error()
}
