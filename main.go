package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/mrbeskin/drone-vision/control"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

const (
	frameSize = 960 * 720 * 3
	mkfifoCmd = "mkfifo"
	fifoPath  = "/tmp/drone_vision_pipe"
)

func main() {
	drone := tello.NewDriver("8888")

	mIn, err := GetMPlayerInput()
	if err != nil {
		fmt.Println(err)
	}

	droneVideoOutput := GetCamStream(drone)

	pipeIn, err := mkfifo()
	if err != nil {
		fmt.Println(err)
	}

	go WriteCameraOutputToMplayerAndPipe(droneVideoOutput, mIn, pipeIn)

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

func WriteCameraOutputToMplayerAndPipe(droneVideoOutput chan []byte, mPlayerIn io.WriteCloser, pipeIn io.WriteCloser) {
	for frame := range droneVideoOutput {
		if _, err := mPlayerIn.Write(frame); err != nil {
			fmt.Printf("failed to write frame to movie player: %v\n", err)
		}
		WriteVideoFeedToNamedPipe(frame, pipeIn)
	}
}

func WriteVideoFeedToNamedPipe(droneFrame []byte, pipeOut io.WriteCloser) {
	_, err := pipeOut.Write(droneFrame)
	if err != nil {
		fmt.Println(err)
	}
}

func mkfifo() (io.WriteCloser, error) {
	cmd := exec.Command(mkfifoCmd, fifoPath)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return os.OpenFile(fifoPath, os.O_RDWR, 0755)
}
