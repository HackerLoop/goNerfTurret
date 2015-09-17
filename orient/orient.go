package orient

import (
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/pca9685"

	log "github.com/Sirupsen/logrus"

	_ "github.com/kidoman/embd/host/rpi"
)

const i2cAddr = 0x40

const servoXChannel = 0
const servoYChannel = 1

const servoPWMFreq = 50

const servoMin = 250
const servoMax = 430
const servoCenter = 340

var turret *Turret

func init() {
	log.Info("Init Orientation")
	// TODO check i2c bus number for raspi 1 or 2
	bus := embd.NewI2CBus(1)
	turret = &Turret{pca9685.New(bus, i2cAddr)}
	turret.Freq = servoPWMFreq

	turret.setPosition(0.5, 0.5)
}

type Event struct {

	X float64 `json:"x"`
	Y float64 `json:"y"`

}

type Turret struct {
	*pca9685.PCA9685
}

func (t *Turret) setX(value int) {
	t.SetPwm(servoXChannel, 0, value)
}

func (t *Turret) setY(value int) {
	t.SetPwm(servoYChannel, 0, value)
}

func (t *Turret) setPosition(xf float64, yf float64) {
	log.Infof("Settings position %f %f", xf, yf)
	var x int
	var y int
	if xf < 0.5 {
		x = int(servoCenter - (0.5 - xf) * (servoCenter - servoMin))
	} else if xf > 0.5 {
		x = int(servoCenter + (xf - 0.5) * (servoMax - servoCenter))
	} else {
		x = servoCenter
	}

	if yf < 0.5 {
		y = int(servoCenter - (0.5 - yf) * (servoCenter - servoMin))
	} else if yf > 0.5 {
		y = int(servoCenter + (yf - 0.5) * (servoMax - servoCenter))
	} else {
		y = servoCenter
	}

	turret.setX(x)
	turret.setY(y)
}

func Start() chan Event {
	c := make(chan Event)

	go func() {
		for {
			event := <-c
			turret.setPosition(event.X, event.Y)
		}
	}()

	return c
}
