package main

import (
	"bytes"
	"image"
	"log"

	desktop "github.com/axet/desktop/go"
	"github.com/mangalaman93/giggle/gigglesync"
	"github.com/mangalaman93/giggle/giggleui"
)

const (
	cAppName           = "giggle"
	cIconFile          = "icons/giggle.png"
	cSettingsMenuEntry = "settings"
	cSettingsIconFile  = "icons/settings.png"
	cExitMenuEntry     = "exit"
	cExitIconFile      = "icons/exit.png"
)

func onClickListener(mn *desktop.Menu) {
	switch mn.Name {
	case cSettingsMenuEntry:
		giggleui.ShowSettingsDialogue()
	case cExitMenuEntry:
		gigglesync.Exit()
	default:
		log.Println("[WARNING] unknown menu entry: ", mn.Name)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("#################### BEGIN OF LOG ##########################")

	iconImages := make(map[string]image.Image)
	iconFiles := []string{cIconFile, cSettingsIconFile, cExitIconFile}
	for _, iconFile := range iconFiles {
		iconData, err := Asset(iconFile)
		if err != nil {
			log.Fatalln("[ERROR] unable to open icon image file")
		}

		iconReader := bytes.NewReader(iconData)
		icon, _, err := image.Decode(iconReader)
		if err != nil {
			panic(err)
		}

		iconImages[iconFile] = icon
	}

	menu := []desktop.Menu{
		{
			Type:    desktop.MenuItem,
			Enabled: true,
			Name:    cSettingsMenuEntry,
			Action:  onClickListener,
			Icon:    iconImages[cSettingsIconFile],
		},
		{
			Type:    desktop.MenuItem,
			Enabled: true,
			Name:    cExitMenuEntry,
			Action:  onClickListener,
			Icon:    iconImages[cExitIconFile],
		},
	}

	deskApp := desktop.DesktopSysTrayNew()
	deskApp.SetIcon(iconImages[cIconFile])
	deskApp.SetTitle(cAppName)
	deskApp.SetMenu(menu)
	deskApp.Show()

	desktop.Main()
}
