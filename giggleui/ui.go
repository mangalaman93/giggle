package giggleui

import (
	"log"
)

type GiggleUI struct {
}

func NewGiggleUI() *GiggleUI {
	return &GiggleUI{}
}

func (gu *GiggleUI) Start() error {
	log.Println("[INFO] starting giggle UI server")

	go gu.run()
	return nil
}

func (gu *GiggleUI) Stop() error {
	log.Println("[INFO] stopping giggle UI server")
	return nil
}

func (gu *GiggleUI) run() {
}
