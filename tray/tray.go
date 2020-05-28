package tray

import (
	"bytes"
	"image"
	"log"

	"github.com/mangalaman93/giggle/conf"
	"github.com/mangalaman93/giggle/content"
	"github.com/skratchdot/open-golang/open"
	desktop "gitlab.com/axet/desktop/go"
)

const (
	cSettingsMenuEntry = "Settings"
	cLogFileMenuEntry  = "Open Logs"
	cExitMenuEntry     = "Exit"
)

var (
	vIconFile         = conf.IconFile()
	vSettingsIconFile = conf.SettingsIconFile()
	vLogIconFile      = conf.LogIconFile()
	vExitIconFile     = conf.ExitIconFile()
)

// GTray is a giggle system tray.
type GTray struct {
	tray *desktop.DesktopSysTray
	quit chan struct{}
}

// Start starts the system tray.
func Start(quit chan struct{}) *GTray {
	log.Println("[INFO] starting giggle tray")

	gt := &GTray{
		tray: desktop.DesktopSysTrayNew(),
		quit: quit,
	}
	go gt.run()
	return gt
}

// Stop stops the system tray.
func (gt *GTray) Stop() error {
	log.Println("[INFO] stopping giggle tray")

	// There is most likely a race condition when interrupt
	// signal is sent. We avoid this by recovering from panic.
	defer func() {
		if r := recover(); r != nil {
			log.Println("[WARNING] systray stop panic ::", r)
		}
	}()

	if gt.tray != nil {
		gt.tray.Close()
	}

	return nil
}

func (gt *GTray) onSettingsMenuClick(mn *desktop.Menu) {
	log.Println("[INFO] settings menu option selected")
	settingsPath := conf.SettingsFilePath()
	if err := open.Start(settingsPath); err != nil {
		log.Printf("[ERROR] unable to open %s :: %v\n", settingsPath, err)
	}
}

func (gt *GTray) onLogFileMenuClick(mn *desktop.Menu) {
	log.Println("[INFO] log file menu option selected")
	logFolder := conf.LogFolder()
	if err := open.Start(logFolder); err != nil {
		log.Printf("[ERROR] unable to open %s :: %v\n", logFolder, err)
	}
}

func (gt *GTray) onExitMenuClick(mn *desktop.Menu) {
	log.Println("[INFO] exit menu option selected")
	gt.quit <- struct{}{}
}

func (gt *GTray) run() {
	iconImages := make(map[string]image.Image)
	iconFiles := []string{vIconFile, vSettingsIconFile, vLogIconFile}

	for _, iconFile := range iconFiles {
		iconData, err := content.Asset(iconFile)
		if err != nil {
			log.Println("[ERROR] unable to open icon image file :: internal error!")
			panic(err)
		}

		iconReader := bytes.NewReader(iconData)
		icon, _, err := image.Decode(iconReader)
		if err != nil {
			log.Println("[ERROR] unable to decode image :: internal error!")
			panic(err)
		}

		iconImages[iconFile] = icon
	}

	menu := []desktop.Menu{
		{
			Type:    desktop.MenuItem,
			Enabled: true,
			Name:    cSettingsMenuEntry,
			Action:  gt.onSettingsMenuClick,
			Icon:    iconImages[vSettingsIconFile],
		},
		{
			Type:    desktop.MenuItem,
			Enabled: true,
			Name:    cLogFileMenuEntry,
			Action:  gt.onLogFileMenuClick,
			Icon:    iconImages[vLogIconFile],
		},
		{
			Type:    desktop.MenuItem,
			Enabled: true,
			Name:    cExitMenuEntry,
			Action:  gt.onExitMenuClick,
			Icon:    iconImages[vExitIconFile],
		},
	}
	log.Println("[INFO] constructed system tray menu")

	gt.tray.SetIcon(iconImages[vIconFile])
	gt.tray.SetTitle(conf.AppName())
	gt.tray.SetMenu(menu)
	gt.tray.Show()
	log.Println("[INFO] done setting up system tray")

	desktop.Main()
}
