package main

import (
	"github.com/vitaminwater/turret/control"
	"github.com/vitaminwater/turret/orient"
	"github.com/vitaminwater/turret/shoot"

	log "github.com/Sirupsen/logrus"
)

func main() {
	orientChan := orient.Start()
	shootChan := shoot.Start()
	control.Start(orientChan, shootChan)

	log.Info("All started")

	select {}
}
