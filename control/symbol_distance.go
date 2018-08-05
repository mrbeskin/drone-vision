package control

import (
	"fmt"
	"math"
)

func InitSymbolVision(driver *tello.Driver) {
	drone := NewTelloDrone(driveR)
	controller = NewFlightcontroller(drone)
}

const MIN_DIST_INCHES = float64(24.0)
const MAX_DIF_X = float64(0.1)
const MAX_DIF_Y = float64(0.1)

func doFlight() {
	// TODO: Parse distance and offset data here
	yAxis := float64(0.0)
	xAxis := float64(0.0)
	distanceInches := float64(0.0)

	fmt.Printf("x axis: %f\n", xAxis)
	// xAxis
	// > 0 is turn right
	if yAxis > 0.0 {
		fmt.Println("right event")
		controller.Right()
	} else if xAxis < 0.0 {
		// < 0 is turn left
		controller.Left()
		fmt.Println("left event")
	}

	fmt.Printf("y axis: %f\n", yAxis)
	// yAxis
	// > 0 is go up
	if yAxis > 0.0 {
		controller.Up()
		fmt.Prinln("up event")
	} else if yAxis < 0.0 {
		// < 0 is go down
		controller.Down()
		fmt.Println("down event")
	}

	fmt.Printf("distance: %f\n", distanceInches)
	// TRACK HORIZTONALLY AND VERTICALLY FIRST, THEN DECIDE BASED ON THRESHOLD
	if axisWithinThreshold(xAxis, yAxis) {
		if distance > MIN_DIST_INCHES {
			controller.Forward()
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
