package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mangalaman93/giggle/giggleconf"
	"github.com/mangalaman93/giggle/giggleservice"
	"github.com/mangalaman93/giggle/giggletray"
	"github.com/mangalaman93/giggle/giggleui"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logFolder := giggleconf.GetLogFolder()
	_, err := os.Stat(logFolder)
	if err != nil {
		if strings.Contains(err.Error(), "cannot find the path") {
			errDir := os.MkdirAll(logFolder, os.FileMode(giggleconf.GetFilePerm()))
			if errDir != nil {
				message := fmt.Sprintf("[ERROR] unable to create app folder :: %v", errDir)
				log.Println(message)
				giggleui.Error(message)
				panic(errDir)
			} else {
				log.Println("[INFO] created log directory")
			}
		} else {
			message := fmt.Sprintf("[ERROR] unable to get app folder stats :: %v", err)
			log.Println(message)
			giggleui.Error(message)
			panic(err)
		}
	} else {
		log.Println("[INFO] log directory already exists")
	}

	logFilePath := giggleconf.GetLogFilePath()
	log.SetOutput(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    giggleconf.GetLogFileMaxSize(),
		MaxBackups: giggleconf.GetLogMaxNumBackups(),
		MaxAge:     giggleconf.GetLogFileMaxAge(),
		LocalTime:  true,
	})

	log.Println("#################### BEGIN OF LOG ##########################")

	// register ctrl+c
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Println("[INFO] adding signal handler for SIGTERM")

	giggleConfig, err := giggleconf.LoadConfigFile()
	if err != nil {
		message := fmt.Sprintf("[ERROR] error in reading config file :: %v", err)
		log.Println(message)
		giggleui.Error(message)
		panic(err)
	}
	log.Println("[INFO] read config file:", giggleConfig)

	// giggle system tray
	giggleTray := giggletray.NewGiggleTray()
	err = giggleTray.Start()
	if err != nil {
		message := fmt.Sprintf("[ERROR] error in starting giggle tray :: %v", err)
		log.Println(message)
		giggleui.Error(message)
		panic(err)
	}
	defer func() {
		err = giggleTray.Stop()
		if err != nil {
			log.Println("[WARN] unable to stop giggle tray ::", err)
		}
	}()

	// giggle UI server
	giggleUI := giggleui.NewGiggleUI()
	err = giggleUI.Start()
	if err != nil {
		message := fmt.Sprintf("[ERROR] error in starting giggle UI Server :: %v", err)
		log.Println(message)
		giggleui.Error(message)
		panic(err)
	}
	defer func() {
		err = giggleUI.Stop()
		if err != nil {
			log.Println("[WARN] unable to stop giggle UI server ::", err)
		}
	}()

	// giggle service
	giggleSvc := giggleservice.NewGiggleService()
	err = giggleSvc.Start()
	if err != nil {
		message := fmt.Sprintf("[ERROR] error in starting giggle service :: %v", err)
		log.Println(message)
		giggleui.Error(message)
		panic(err)
	}
	defer func() {
		err = giggleSvc.Stop()
		if err != nil {
			log.Println("[WARN] unable to stop giggle service ::", err)
		}
	}()

	// wait for ctrl+c
	log.Println("[INFO] waiting for ctrl+c signal")
	<-sigs

	log.Println("[INFO] exiting giggle")
}
