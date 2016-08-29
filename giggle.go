package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/axet/desktop/go"
	"github.com/kardianos/service"
	"github.com/mangalaman93/giggle/giggleservice"
	"github.com/mangalaman93/giggle/giggleui"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	cAppName          = "giggle"
	cAppDescription   = "syncs Overleaf repositories with Github"
	cAppFolder        = "giggle"
	cLogFolder        = "log"
	cLogFile          = "giggle.log"
	cFilePerm         = 640
	cLogFileMaxSize   = 50 // MB
	cLogMaxNumBackups = 5
	cLogMaxAge        = 30 // days
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	homePath := desktop.GetAppDataFolder()
	logFolderPath := filepath.Join(homePath, cAppFolder, cLogFolder)
	_, err := os.Stat(logFolderPath)

	if err != nil {
		if strings.Contains(err.Error(), "cannot find the path") {
			errDir := os.MkdirAll(logFolderPath, os.FileMode(cFilePerm))
			if errDir != nil {
				message := fmt.Sprintf("[ERROR] unable to create app folder :: %v", errDir)
				giggleui.ShowDialog(message)
				log.Fatalln(message)
			} else {
				log.Println("[INFO] created log directory")
			}
		} else {
			message := fmt.Sprintf("[ERROR] unable to get app folder stats :: %v", err)
			giggleui.ShowDialog(message)
			log.Fatalln(message)
		}
	} else {
		log.Println("[INFO] log directory already exists")
	}
}

func main() {
	homePath := desktop.GetAppDataFolder()
	appPath := filepath.Join(homePath, cAppFolder, cLogFolder)
	logFilePath := filepath.Join(appPath, cLogFile)
	log.SetOutput(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    cLogFileMaxSize,
		MaxBackups: cLogMaxNumBackups,
		MaxAge:     cLogMaxAge,
		LocalTime:  true,
	})
	log.Println("#################### BEGIN OF LOG ##########################")

	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        cAppName,
		DisplayName: cAppName,
		Description: cAppDescription,
	}

	giggleService := giggleservice.NewGiggleService(logFilePath, cAppName)
	gService, errSvc := service.New(giggleService, svcConfig)
	if errSvc != nil {
		message := fmt.Sprintf("[ERROR] unable to create giggle service :: %v", errSvc)
		giggleui.ShowDialog(message)
		log.Fatalln(message)
	}
	log.Println("[INFO] ready to run giggle service")

	if len(*svcFlag) != 0 {
		errCtrl := service.Control(gService, *svcFlag)
		if errCtrl != nil {
			message := fmt.Sprintf("[ERROR] valid actions for giggle service %q :: %v",
				service.ControlAction, errCtrl)
			giggleui.ShowDialog(message)
			log.Fatalln(message)
		}

		return
	}

	errRun := gService.Run()
	if errRun != nil {
		message := fmt.Sprintf("[ERROR] unable to run giggle service :: %v", errRun)
		giggleui.ShowDialog(message)
		log.Fatalln(message)
	}

	log.Println("[INFO] exiting giggle service")
}
