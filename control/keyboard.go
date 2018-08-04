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
				fmt.Println("up Pressed")
			case keyboard.A:
				controller.Left()
				fmt.Println("left pressed")
			case keyboard.S:
				controller.Backward()
				fmt.Println("down pressed")
			case keyboard.D:
				controller.Right()
				fmt.Println("right pressed")
			case keyboard.Q:
				controller.CounterClockwise()
				fmt.Println("q pressed")
			case keyboard.E:
				controller.Clockwise()
				fmt.Println("e pressed")
			case keyboard.R:
				controller.Up()
				fmt.Println("r pressed")
			case keyboard.F:
				controller.Down()
				fmt.Println("f pressed")
			case keyboard.L:
				driver.Land()
			case keyboard.T:
				driver.TakeOff()
				fmt.Println("t pressed")
			case keyboard.M:
				controller.ThrottleUp()
				fmt.Println("m pressed")
			case keyboard.N:
				controller.ThrottleDown()
				fmt.Println("n pressed")
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
