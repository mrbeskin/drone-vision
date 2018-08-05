package control

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gobot.io/x/gobot/platforms/keyboard"
)

const speed = 40

func InitControl(driver *tello.Driver) {

	drone := NewTelloDrone(driver)
	controller := NewFlightController(drone)

	keyb := keyboard.NewDriver()

	fmt.Println("Drone: control intitialized")

	work := func() {

		keyb.On(keyboard.Key, func(data interface{}) {

			k := data.(keyboard.KeyEvent)

			switch k.Key {
			case keyboard.W:
				controller.Forward()
			case keyboard.A:
				controller.Left()
			case keyboard.S:
				controller.Backward()
			case keyboard.D:
				controller.Right()
			case keyboard.Q:
				controller.CounterClockwise()
			case keyboard.E:
				controller.Clockwise()
			case keyboard.R:
				controller.Up()
			case keyboard.F:
				controller.Down()
			case keyboard.L:
				driver.Land()
			case keyboard.T:
				driver.TakeOff()
			case keyboard.M:
				controller.ThrottleUp()
			case keyboard.N:
				controller.ThrottleDown()
			}

		})

	}

	robot := gobot.NewRobot("keyboardbot",
		[]gobot.Connection{},
		[]gobot.Device{keyb},
		work,
	)

	fmt.Println("keyboard initialized")
	controller.StartControl()
	robot.Start()
}
