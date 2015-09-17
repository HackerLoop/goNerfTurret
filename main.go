package main

import (
	"github.com/vitaminwater/turret/control"
	"github.com/vitaminwater/turret/orient"
	"github.com/vitaminwater/turret/shoot"

	_ "github.com/kidoman/embd/host/rpi"
)

func main() {
	orientChan := orient.Start()
	shootChan := shoot.Start()
	control.Start(orientChan, shootChan)

	select {}
}
