package giggleservice

import (
	"log"
)

type GiggleService struct {
}

func NewGiggleService() *GiggleService {
	return &GiggleService{
	}
}

func (gs *GiggleService) Start() {
	log.Println("[INFO] starting giggle service")
	go gs.run()
}

func (gs *GiggleService) Stop() error {
	log.Println("[INFO] stopping giggle service")
	return nil
}

func (gs *GiggleService) run() {
}
