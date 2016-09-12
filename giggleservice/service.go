package giggleservice

import (
	"log"
)

type GiggleService struct {
	running bool
}

func NewGiggleService() *GiggleService {
	return &GiggleService{
		running: false,
	}
}

func (gs *GiggleService) Start() error {
	log.Println("[INFO] starting giggle service")

	go gs.run()
	return nil
}

func (gs *GiggleService) Stop() error {
	if gs.running {
		log.Println("[INFO] stopping giggle service")
	} else {
		log.Println("[INFO] giggle service is not running")
	}

	return nil
}

func (gs *GiggleService) run() {
	gs.running = true
	defer func() { gs.running = false }()
}
