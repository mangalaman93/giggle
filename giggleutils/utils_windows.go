package giggleutils

import (
  "github.com/luisiturrios/gowin"
)

func GetAppFolder() string {
  folders := gowin.ShellFolders{gowin.ALL}
  return folders.AppData()
}
