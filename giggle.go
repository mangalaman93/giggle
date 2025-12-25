package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
	"github.com/mangalaman93/giggle/conf"
	"github.com/mangalaman93/giggle/svc"
	"github.com/mangalaman93/giggle/tray"
	"github.com/sevlyar/go-daemon"
	"github.com/sqweek/dialog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func dialogAndPanic(message string, err error) {
	log.Println(message)
	dialog.Message(message).Error()
	panic(err)
}

func main() {
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

	dctx := &daemon.Context{
		PidFileName: conf.PidFilePath(),
		PidFilePerm: 0644,
	}
	child, err := dctx.Reborn()
	if err != nil {
		message := fmt.Sprintf("[ERROR] unable to daemonize :: %v", err)
		dialogAndPanic(message, err)
	}

	if child != nil {
		log.Println("running the service as a daemon")
	} else {
		defer func() {
			if err := dctx.Release(); err != nil {
				log.Printf("error releasing daemon context: %v", err)
			}
		}()
		runChild()
	}
}

func runChild() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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

	// giggle system tray
	quit := make(chan struct{})
	gt := tray.Start(quit)
	defer func() {
		if err := gt.Stop(); err != nil {
			log.Println("[WARN] unable to stop giggle tray ::", err)
		}
	}()

	// giggle service
	gsvc := svc.Start()
	defer func() {
		if err := gsvc.Stop(); err != nil {
			log.Println("[WARN] unable to stop giggle service ::", err)
		}
	}()

	go func() {
		// wait for ctrl+c
		log.Println("[INFO] waiting for ctrl+c signal")
		select {
		case <-quit:
		case <-sigs:
		}

		systray.Quit()
	}()

	// This has to be called here in the main thread, fails on mac otherwise.
	systray.Run(gt.OnReady, nil)
	log.Println("[INFO] exiting giggle")
}
