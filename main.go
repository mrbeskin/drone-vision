package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
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
	if pipeIn == nil {
		panic("pipe nil")
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
	fmt.Println("done")
	if err != nil {
		return nil, err
	}
	return os.OpenFile(fifoPath, os.O_RDWR, os.ModeNamedPipe)
}

func runGuidedDrone() io.ReadCloser {
	cmd := exec.Command("docker", "run", "-it", "-a", "STDOUT", "-a", "STDIN", "--name", "guided_flying", "--rm", "--mount", "type=bind,bind,source=\"/tmp/\",target=/apptarget,readonly", "guided_flying")
	defer func() {
		err := cmd.Start()
		if err != nil {
			panic(err)
		}
	}()
	out, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	return out
}

func parseValues(xydistMsg string) (x, y, dist float64, err error) {
	// ex. [1.0,-1.0,10.0]
	s := strings.Trim(xydistMsg, " \n[]")
	sValArray := strings.Split(s, ",")
	if len(sValArray) > 3 {
		return float64(0.0), float64(0.0), float64(0.0), fmt.Errorf("string value array parsed does not three values, abort! value: %v", xydistMsg)
	}
	// TODO: errors
	x, err = strconv.ParseFloat(sValArray[0], 64)
	checkFloatErr(err)
	y, err = strconv.ParseFloat(sValArray[1], 64)
	checkFloatErr(err)
	dist, err = strconv.ParseFloat(sValArray[2], 64)
	checkFloatErr(err)
	return x, y, dist, nil
}

func checkFloatErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("couldn't parse float from model: %v", err))
	}
}
