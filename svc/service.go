package svc

import (
	"context"
	"log"
	"time"

	"github.com/mangalaman93/giggle/conf"
)

// Service runs a periodic sync between various repos.
type Service struct {
	quit chan struct{}
	done chan struct{}
}

// Start starts the service.
func Start() *Service {
	log.Println("[INFO] starting giggle service")

	gs := &Service{
		quit: make(chan struct{}),
		done: make(chan struct{}),
	}
	go gs.run()
	return gs
}

// Stop stops the service.
func (gs *Service) Stop() error {
	log.Println("[INFO] stopping giggle service")
	close(gs.quit)
	<-gs.done
	log.Println("[INFO] stopped giggle service")
	return nil
}

func (gs *Service) run() {
	defer close(gs.done)

	for {
		cf, err := conf.ReadConfig(conf.SettingsFilePath())
		if err != nil {
			log.Printf("[ERROR] error in reading config file :: %v\n", err)
			log.Println("sleeping for a minute...")
			time.Sleep(time.Minute)
			continue
		}

		select {
		case <-gs.quit:
			log.Println("[INFO] exiting service loop")
			return
		case <-time.After(cf.Period.Duration):
			if err := performSync(context.Background(), cf); err != nil {
				log.Printf("[ERROR] error syncing :: %v\n", err)
			}
		}
	}
}
