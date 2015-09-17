package orient

import (
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/pca9685"
)

const i2cAddr = 0x40

const servoXChannel = 0
const servoYChannel = 0

const servoPWMFreq = 50

const servoMin = 250
const servoMax = 430
const servoCenter = 340

var turret *Turret

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

func init() {
	bus := embd.NewI2CBus(1)
	turret = &Turret{pca9685.New(bus, i2cAddr)}
	turret.Freq = servoPWMFreq
}

func orient(e Event) {
	var x int
	var y int
	if e.X < 0.5 {
		x = int(servoCenter - (0.5 - e.X) * (servoCenter - servoMin))
	} else {
		x = int(servoCenter + (e.X - 0.5) * (servoMax - servoCenter))
	}

	if e.Y < 0.5 {
		y = int(servoCenter - (0.5 - e.Y) * (servoCenter - servoMin))
	} else {
		y = int(servoCenter + (e.Y - 0.5) * (servoMax - servoCenter))
	}

	turret.setX(x)
	turret.setY(y)
}

func Start() chan Event {
	c := make(chan Event)

	go func() {
		for {
			event := <-c
			orient(event)
		}
	}()

	return c
}
