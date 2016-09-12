package giggletray

import (
	"bytes"
	"fmt"
	"image"
	"log"

	"github.com/axet/desktop/go"
	"github.com/mangalaman93/giggle/giggleconf"
	"github.com/mangalaman93/giggle/gigglecontent"
	"github.com/mangalaman93/giggle/giggleui"
	"github.com/skratchdot/open-golang/open"
)

const (
	cSettingsMenuEntry = "Settings"
	cLogFileMenuEntry  = "Open Logs"
)

var (
	vIconFile         = giggleconf.GetIconFile()
	vSettingsIconFile = giggleconf.GetSettingsIconFile()
	vLogIconFile      = giggleconf.GetLogIconFile()
)

type GiggleTray struct {
	tray *desktop.DesktopSysTray
}

func NewGiggleTray() *GiggleTray {
	return &GiggleTray{
		tray: desktop.DesktopSysTrayNew(),
	}
}

func (gt *GiggleTray) Start() error {
	log.Println("[INFO] starting giggle tray")

	go gt.run()
	return nil
}

func (gt *GiggleTray) Stop() error {
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

func (gt *GiggleTray) onClickListener(mn *desktop.Menu) {
	switch mn.Name {
	case cSettingsMenuEntry:
		log.Println("[INFO] settings menu option selected")
		urlToUI := giggleconf.GetURLToUI()
		err := open.Run(urlToUI)
		if err != nil {
			message := fmt.Sprintf("[ERROR] unable to open url %s :: %v", urlToUI, err)
			giggleui.Error(message)
			log.Println(message)
		}
	case cLogFileMenuEntry:
		log.Println("[INFO] log file menu option selected")
		logFolder := giggleconf.GetLogFolder()
		err := open.Start(logFolder)
		if err != nil {
			message := fmt.Sprintf("[ERROR] unable to open log folder %s :: %v",
				logFolder, err)
			giggleui.Error(message)
			log.Println(message)
		}
	default:
		log.Println("[WARNING] unknown menu entry:", mn.Name)
	}
}

func (gt *GiggleTray) run() {
	iconImages := make(map[string]image.Image)
	iconFiles := []string{vIconFile, vSettingsIconFile, vLogIconFile}

	for _, iconFile := range iconFiles {
		iconData, err := gigglecontent.Asset(iconFile)
		if err != nil {
			message := "[ERROR] unable to open icon image file :: internal error!"
			log.Println(message)
			giggleui.Error(message)
			panic(err)
		}

		iconReader := bytes.NewReader(iconData)
		icon, _, err := image.Decode(iconReader)
		if err != nil {
			message := "[ERROR] unable to decode image :: internal error!"
			log.Println(message)
			giggleui.Error(message)
			panic(err)
		}

		iconImages[iconFile] = icon
	}

	menu := []desktop.Menu{
		{
			Type:    desktop.MenuItem,
			Enabled: true,
			Name:    cSettingsMenuEntry,
			Action:  gt.onClickListener,
			Icon:    iconImages[vSettingsIconFile],
		},
		{
			Type:    desktop.MenuItem,
			Enabled: true,
			Name:    cLogFileMenuEntry,
			Action:  gt.onClickListener,
			Icon:    iconImages[vLogIconFile],
		},
	}
	log.Println("[INFO] constructed system tray menu")

	gt.tray.SetIcon(iconImages[vIconFile])
	gt.tray.SetTitle(giggleconf.GetAppName())
	gt.tray.SetMenu(menu)
	gt.tray.Show()
	log.Println("[INFO] done setting up system tray")

	desktop.Main()
}
