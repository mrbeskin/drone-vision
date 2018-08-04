package main

import (
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"time"

	"github.com/mrbeskin/drone-hack/control"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func main() {
	drone := tello.NewDriver("8888")

	mIn, err := GetMPlayerInput()
	if err != nil {
		fmt.Println(err)
	}

	droneVideoOutput := GetCamStream(drone)

	go WriteCameraOutputToMplayer(droneVideoOutput, mIn)

	work := func() {
		gobot.After(5*time.Second, func() {
			drone.TakeOff()
		})

		go control.InitControl(drone)
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}

func GetMPlayerInput() (io.WriteCloser, error) {
	var mPlayer *exec.Cmd
	if runtime.GOOS == "darwin" {
		fmt.Println("Mac OS detected")
		mPlayer = exec.Command("mplayer", "-vo", "-fps", "30", "-")
	} else {
		mPlayer = exec.Command("mplayer", "-vo", "x11", "-fps", "30", "-")
	}
	defer mPlayer.Start()
	return mPlayer.StdinPipe()
}

func WriteCameraOutputToMplayer(droneVideoOutput chan []byte, mPlayerIn io.WriteCloser) {
	for frame := range droneVideoOutput {
		if _, err := mPlayerIn.Write(frame); err != nil {
			fmt.Printf("failed to write frame to movie player: %v\n", err)
		}
	}
}
