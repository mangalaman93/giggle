package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	desktop "github.com/axet/desktop/go"
	"github.com/mangalaman93/giggle/gigglesync"
	"github.com/mangalaman93/giggle/giggleui"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	cAppName           = "giggle"
	cSettingsMenuEntry = "settings"
	cExitMenuEntry     = "exit"
	cExitIconFile      = "icons/exit.png"
	cIconFile          = "icons/giggle.png"
	cSettingsIconFile  = "icons/settings.png"
	cAppFolder         = ".giggle"
	cLogFolder         = "log"
	cLogFile           = "giggle.log"
	cFilePerm          = 660
	cLogFileMaxSize    = 50 // MB
	cLogMaxNumBackups  = 5
	cLogMaxAge         = 30 // days
)

var deskApp *desktop.DesktopSysTray

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	homePath := desktop.GetHomeFolder()
	logFolder := filepath.Join(homePath, cAppFolder, cLogFolder)
	_, err := os.Stat(logFolder)

	if err != nil {
		if strings.Contains(err.Error(), "cannot find the path") {
			errDir := os.MkdirAll(logFolder, os.FileMode(cFilePerm))
			if errDir != nil {
				message := fmt.Sprintf("[ERROR] unable to create app folder :: %v", errDir)
				giggleui.ShowDialog(message)
				log.Fatalln(message)
			}
		} else {
			message := fmt.Sprintf("[ERROR] unable to get app folder stats :: %v", err)
			giggleui.ShowDialog(message)
			log.Fatalln(message)
		}
	}
}

func onClickListener(mn *desktop.Menu) {
	switch mn.Name {
	case cSettingsMenuEntry:
		giggleui.ShowSettingsDialog()
	case cExitMenuEntry:
		deskApp.Close()
		log.Println("exiting system tray!")
		gigglesync.Exit()
	default:
		log.Println("[WARNING] unknown menu entry:", mn.Name)
	}
}

func main() {
	homePath := desktop.GetHomeFolder()
	logFilePath := filepath.Join(homePath, cAppFolder, cLogFolder, cLogFile)
	log.SetOutput(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    cLogFileMaxSize,
		MaxBackups: cLogMaxNumBackups,
		MaxAge:     cLogMaxAge,
		LocalTime:  true,
	})
	log.Println("#################### BEGIN OF LOG ##########################")

	iconImages := make(map[string]image.Image)
	iconFiles := []string{cIconFile, cSettingsIconFile, cExitIconFile}
	for _, iconFile := range iconFiles {
		iconData, err := Asset(iconFile)
		if err != nil {
			message := "[ERROR] unable to open icon image file :: internal error!"
			giggleui.ShowDialog(message)
			log.Fatalln(message)
		}

		iconReader := bytes.NewReader(iconData)
		icon, _, err := image.Decode(iconReader)
		if err != nil {
			message := "[ERROR] unable to decode image :: internal error!"
			giggleui.ShowDialog(message)
			log.Fatalln(message)
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

	deskApp = desktop.DesktopSysTrayNew()
	deskApp.SetIcon(iconImages[cIconFile])
	deskApp.SetTitle(cAppName)
	deskApp.SetMenu(menu)
	deskApp.Show()
	log.Println("done setting up system tray")

	desktop.Main()
}
