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
func Start(ch *conf.ConfigHolder) *Service {
	log.Println("[INFO] starting giggle service")

	gs := &Service{
		quit: make(chan struct{}),
		done: make(chan struct{}),
	}
	go gs.run(ch)
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

func (gs *Service) run(ch *conf.ConfigHolder) {
	defer close(gs.done)

	tick := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-gs.quit:
			log.Println("[INFO] exiting service loop")
			return
		case <-tick.C:
			if err := performSync(context.Background(), ch); err != nil {
				log.Printf("[ERROR] error in sync :: %v\n", err)
			}
		}
	}
}
