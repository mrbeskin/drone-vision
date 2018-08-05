package control

import (
	"fmt"
	"math"

	"gobot.io/x/gobot/platforms/dji/tello"
)

type SymbolVisionController struct {
	controller *FlightController
}

type SymbolVision struct {
}

func InitSymbolVision(driver *tello.Driver) *SymbolVisionController {
	drone := NewTelloDrone(driver)
	controller = NewFlightController(drone)
	return &SymbolVisionController{
		controller: controller,
	}
}

const MIN_DIST_INCHES = float64(24.0)
const MAX_DIF_X = float64(0.1)
const MAX_DIF_Y = float64(0.1)

func (svc *SymbolVisionController) ListenToYourBrain() {
	// initialize brain
	// start listener to stdout
	// intermittently do flight
	// stop if no new instructions for a while
}

func (svc *SymbolVisionController) doFlight(xAxis float64, yAxis float64, distance float64) {
	fmt.Printf("x axis: %f\n", xAxis)
	// xAxis
	// > 0 is turn right
	if yAxis > 0.0 {
		fmt.Println("right event")
		svc.controller.Right()
	} else if xAxis < 0.0 {
		// < 0 is turn left
		svc.controller.Left()
		fmt.Println("left event")
	}

	fmt.Printf("y axis: %f\n", yAxis)
	// yAxis
	// > 0 is go up
	if yAxis > 0.0 {
		svc.controller.Up()
		fmt.Prinln("up event")
	} else if yAxis < 0.0 {
		// < 0 is go down
		svc.controller.Down()
		fmt.Println("down event")
	}

	fmt.Printf("distance: %f\n", distanceInches)
	// TRACK HORIZTONALLY AND VERTICALLY FIRST, THEN DECIDE BASED ON THRESHOLD
	if axisWithinThreshold(xAxis, yAxis) {
		if distance > MIN_DIST_INCHES {
			svc.controller.Forward()
			fmt.Println("forward event")
		}
	}
}

func axisWithinThreshold(xAxis float64, yAxis float64) bool {
	if (math.Abs(xAxis) < MAX_DIF_X) && (math.Abs(yAxis) < MAX_DIF_Y) {
		return true
	}
	return false
}
