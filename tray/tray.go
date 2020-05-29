package tray

import (
	"log"
	"sync"

	"github.com/getlantern/systray"
	"github.com/mangalaman93/giggle/conf"
	"github.com/mangalaman93/giggle/content"
	"github.com/skratchdot/open-golang/open"
)

const (
	cSettingsMenuEntry = "Settings"
	cLogFileMenuEntry  = "Open Logs"
	cExitMenuEntry     = "Exit"
)

// GTray is a giggle system tray.
type GTray struct {
	mainQuit chan struct{}
	quit     chan struct{}
	wg       sync.WaitGroup
}

// Start starts the system tray.
func Start(mainQuit chan struct{}) *GTray {
	log.Println("[INFO] starting giggle tray")
	return &GTray{mainQuit: mainQuit, quit: make(chan struct{})}
}

// Stop stops the system tray.
func (gt *GTray) Stop() error {
	log.Println("[INFO] stopping giggle tray")
	gt.quit <- struct{}{}
	gt.wg.Wait()
	return nil
}

func (gt *GTray) onSettingsMenuClick() {
	log.Println("[INFO] settings menu option selected")
	settingsPath := conf.SettingsFilePath()
	if err := open.Start(settingsPath); err != nil {
		log.Printf("[ERROR] unable to open %s :: %v\n", settingsPath, err)
	}
}

func (gt *GTray) onLogFileMenuClick() {
	log.Println("[INFO] log file menu option selected")
	logFolder := conf.LogFolder()
	if err := open.Start(logFolder); err != nil {
		log.Printf("[ERROR] unable to open %s :: %v\n", logFolder, err)
	}
}

func (gt *GTray) onExitMenuClick() {
	log.Println("[INFO] exit menu option selected")
	gt.mainQuit <- struct{}{}
}

// OnReady is the function that is passed to systray.Run.
func (gt *GTray) OnReady() {
	systray.SetIcon(getIcon(conf.IconFile()))
	systray.SetTooltip(conf.AppName())

	settingsMenu := systray.AddMenuItem(cSettingsMenuEntry, "")
	settingsMenu.SetIcon(getIcon(conf.SettingsIconFile()))

	logFileMenu := systray.AddMenuItem(cLogFileMenuEntry, "")
	logFileMenu.SetIcon(getIcon(conf.LogIconFile()))

	exitMenu := systray.AddMenuItem(cExitMenuEntry, "")
	exitMenu.SetIcon(getIcon(conf.ExitIconFile()))

	gt.wg.Add(1)
	go gt.handleClicks(settingsMenu, logFileMenu, exitMenu)
	log.Println("[INFO] constructed system tray menu")
}

func (gt *GTray) handleClicks(settingsMenu, logFileMenu, exitMenu *systray.MenuItem) {
	defer gt.wg.Done()

	for {
		select {
		case <-gt.quit:
			return
		case <-settingsMenu.ClickedCh:
			gt.onSettingsMenuClick()
		case <-logFileMenu.ClickedCh:
			gt.onLogFileMenuClick()
		case <-exitMenu.ClickedCh:
			gt.onExitMenuClick()
		}
	}
}

func getIcon(iconFile string) []byte {
	iconData, err := content.Asset(iconFile)
	if err != nil {
		log.Printf("[ERROR] unable to open icon for %v\n", iconFile)
		panic(err)
	}

	return iconData
}
