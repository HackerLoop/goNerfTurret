package shoot

import (
	"time"
	"sync"

	"github.com/kidoman/embd"
	log "github.com/Sirupsen/logrus"

	_ "github.com/kidoman/embd/host/rpi"
)

const trigger1 = 14
const trigger2 = 15

const triggerReturn = 18

const slowSwitch = 27
const fastSwitch = 17

var pinTrigger1 embd.DigitalPin
var pinTrigger2 embd.DigitalPin
var pinTriggerReturn embd.DigitalPin

var pinSlowSwitch embd.DigitalPin
var pinFastSwitch embd.DigitalPin

func inputPin(n int) embd.DigitalPin {
	pin, err := embd.NewDigitalPin(n)
	if err != nil {
		panic(err)
	}
	pin.SetDirection(embd.In)

	return pin
}

func outputPin(n int, value int) embd.DigitalPin {
	pin, err := embd.NewDigitalPin(n)
	if err != nil {
		panic(err)
	}
	pin.SetDirection(embd.Out)
	pin.Write(value)

	return pin
}

func init() {
	log.Info("Init Shooter")
	embd.InitGPIO()

	pinTrigger1 = outputPin(trigger1, embd.High)
	pinTrigger2 = outputPin(trigger2, embd.High)
	pinTriggerReturn = inputPin(trigger2)
	pinSlowSwitch = outputPin(slowSwitch, embd.High)
	pinFastSwitch = outputPin(fastSwitch, embd.Low)

	startShooter()
}

type Event struct {
	N int `json:"n"`
}

var toShootMutex sync.Mutex
var toShoot = 0

func addToShoot(n int) {
	toShootMutex.Lock()
	defer toShootMutex.Unlock()

	log.Infof("Shooting %d darts", n)
	toShoot += n
}

func getToShoot() (n int) {
	toShootMutex.Lock()
	defer toShootMutex.Unlock()

	n = toShoot
	return
}

func awaitTriggerReturnValue(value int) {
	for {
		if val, err := pinTriggerReturn.Read(); val != value {
			time.Sleep(50 * time.Millisecond)
		} else if err != nil {
			panic(err)
		} else {
			break
		}

	}
}

func startShooter() {
	go func () {
		for {
			if (getToShoot() <= 0) {
				time.Sleep(100 * time.Millisecond)
				return
			}

			pinTrigger1.Write(embd.High)
			pinTrigger2.Write(embd.Low)

			pinFastSwitch.Write(embd.High)
			for {
				awaitTriggerReturnValue(embd.Low)
				awaitTriggerReturnValue(embd.High)

				addToShoot(-1)
				if getToShoot() <= 0 {
					pinFastSwitch.Write(embd.Low)
					break
				}
			}

			pinTrigger1.Write(embd.Low)
			pinTrigger2.Write(embd.High)
			time.Sleep(20 * time.Millisecond)

			pinTrigger1.Write(embd.High)
			pinTrigger2.Write(embd.High)
			time.Sleep(300 * time.Millisecond)

			if value, err := pinTriggerReturn.Read(); value != embd.Low {
				pinTrigger1.Write(embd.Low)
				pinTrigger2.Write(embd.High)

				awaitTriggerReturnValue(embd.High)

				pinTrigger1.Write(embd.High)
				pinTrigger2.Write(embd.High)

			} else if err != nil {
				panic(err)
			}
		}
	}()
}

func Start() chan Event {
	c := make(chan Event)

	go func() {
		for {
			event := <-c

			addToShoot(event.N)
		}
	}()

	return c
}
