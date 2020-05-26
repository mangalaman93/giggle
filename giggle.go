package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mangalaman93/giggle/conf"
	"github.com/mangalaman93/giggle/svc"
	"github.com/mangalaman93/giggle/tray"
	"github.com/sqweek/dialog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func dialogAndPanic(message string, err error) {
	log.Println(message)
	dialog.Message(message).Error()
	panic(err)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logFolder := conf.LogFolder()
	if _, err := os.Stat(logFolder); err != nil {
		if os.IsNotExist(err) {
			errDir := os.MkdirAll(logFolder, conf.DirPerm())
			if errDir != nil {
				message := fmt.Sprintf("[ERROR] unable to create log folder :: %v", errDir)
				dialogAndPanic(message, errDir)
			} else {
				log.Println("[INFO] created log directory")
			}
		} else {
			message := fmt.Sprintf("[ERROR] unable to get log folder stats :: %v", err)
			dialogAndPanic(message, err)
		}
	} else {
		log.Println("[INFO] log directory already exists")
	}

	logFilePath := conf.LogFilePath()
	log.SetOutput(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    conf.LogFileMaxSize(),
		MaxBackups: conf.LogMaxNumBackups(),
		MaxAge:     conf.LogFileMaxAge(),
		LocalTime:  true,
	})

	log.Println("#################### BEGIN OF LOG ##########################")

	// register ctrl+c
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Println("[INFO] adding signal handler for SIGTERM")

	// read configuration file.
	config, err := conf.New()
	if err != nil {
		message := fmt.Sprintf("[ERROR] error reading config file :: %v", err)
		dialogAndPanic(message, err)
	}
	log.Printf("[INFO] read config file: %+v\n", config)

	// giggle system tray
	stray := tray.Start()
	defer func() {
		if err := stray.Stop(); err != nil {
			log.Println("[WARN] unable to stop giggle tray ::", err)
		}
	}()

	// giggle service
	gsvc := svc.Start(config)
	defer func() {
		if err := gsvc.Stop(); err != nil {
			log.Println("[WARN] unable to stop giggle service ::", err)
		}
	}()

	// wait for ctrl+c
	log.Println("[INFO] waiting for ctrl+c signal")
	<-sigs
	log.Println("[INFO] exiting giggle")
}
