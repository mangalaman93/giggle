package giggletray

import (
	"bytes"
	"fmt"
	"image"
	"log"

	"github.com/axet/desktop/go"
	"github.com/kardianos/service"
	"github.com/mangalaman93/giggle/giggleui"
	"github.com/skratchdot/open-golang/open"
)

const (
	cSettingsMenuEntry = "Settings"
	cLogFileMenuEntry  = "Open Logs"
	cExitMenuEntry     = "Exit"
	cIconFile          = "icons/giggle.png"
	cSettingsIconFile  = "icons/settings.png"
	clogFileMenuFile   = "icons/log.png"
	cExitIconFile      = "icons/exit.png"
)

type GiggleService struct {
	appName     string
	deskApp     *desktop.DesktopSysTray
	logFilePath string
}

func NewGiggleService(logFilePath, appName string) *GiggleService {
	return &GiggleService{
		appName:     appName,
		deskApp:     nil,
		logFilePath: logFilePath,
	}
}

func (gs *GiggleService) Start(s service.Service) error {
	log.Println("[INFO] starting giggle service")

	go gs.run()
	return nil
}

func (gs *GiggleService) Stop(s service.Service) error {
	log.Println("[INFO] stopping giggle service")

	// There is most likely a race condition when interrupt
	// signal is sent. We avoid this by recovering from panic.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[WARNING] systray stop panic ::", r)
		}
	}()

	if gs.deskApp != nil {
		gs.deskApp.Close()
	}

	return nil
}

func (gs *GiggleService) onClickListener(mn *desktop.Menu) {
	switch mn.Name {
	case cSettingsMenuEntry:
		log.Println("[INFO] settings menu option selected")
		giggleui.ShowSettingsDialog()
	case cLogFileMenuEntry:
		log.Println("[INFO] log file menu option selected")
		errOpen := open.Run(gs.logFilePath)
		if errOpen != nil {
			message := fmt.Sprintf("[ERROR] unable to open log file %s :: %v",
				gs.logFilePath, errOpen)
			giggleui.ShowDialog(message)
			log.Fatalln(message)
		}
	case cExitMenuEntry:
		log.Println("[INFO] exit menu option selected, closing systray!")
		gs.deskApp.Close()
	default:
		log.Println("[WARNING] unknown menu entry:", mn.Name)
	}
}

func (gs *GiggleService) run() {
	iconImages := make(map[string]image.Image)
	iconFiles := []string{cIconFile, cSettingsIconFile, clogFileMenuFile, cExitIconFile}
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
			Action:  gs.onClickListener,
			Icon:    iconImages[cSettingsIconFile],
		},
		{
			Type:    desktop.MenuItem,
			Enabled: true,
			Name:    cLogFileMenuEntry,
			Action:  gs.onClickListener,
			Icon:    iconImages[clogFileMenuFile],
		},
		{
			Type:    desktop.MenuItem,
			Enabled: true,
			Name:    cExitMenuEntry,
			Action:  gs.onClickListener,
			Icon:    iconImages[cExitIconFile],
		},
	}
	log.Println("[INFO] constructed system tray menu")

	gs.deskApp = desktop.DesktopSysTrayNew()
	gs.deskApp.SetIcon(iconImages[cIconFile])
	gs.deskApp.SetTitle(gs.appName)
	gs.deskApp.SetMenu(menu)
	gs.deskApp.Show()
	log.Println("[INFO] done setting up system tray")

	desktop.Main()
}
