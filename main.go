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
	"gocv.io/x/gocv"
)

const (
	frameSize = 960 * 720 * 3
)

func main() {
	drone := tello.NewDriver("8888")

	mIn, err := GetMPlayerInput()
	if err != nil {
		fmt.Println(err)
	}

	droneVideoOutput := GetCamStream(drone)

	ffmpegIn, ffmpegOut := InitFfmpeg()

	go WriteCameraOutputToMplayerAndFfmpeg(droneVideoOutput, ffmpegIn, mIn)

	work := func() {
		gobot.After(5*time.Second, func() {
			drone.TakeOff()
		})
		go writeFramesToTmp(ffmpegOut)
		go control.InitControl(drone)
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}

func writeFramesToTmp(ffmpegOut io.ReadCloser) {
	for {
		buf := make([]byte, frameSize)
		if _, err := io.ReadFull(ffmpegOut, buf); err != nil {
			fmt.Println(err)
			continue
		}

		img, err := gocv.NewMatFromBytes(720, 960, gocv.MatTypeCV8UC3, buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		gocv.IMWrite("/tmp/drone-vision/capture.png", img)
	}
}

func InitFfmpeg() (io.WriteCloser, io.ReadCloser) {
	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-pix_fmt", "bgr24", "-vcodec", "rawvideo",
		"-an", "-sn", "-s", "960x720", "-f", "rawvideo", "pipe:1")
	ffmpegIn, _ := ffmpeg.StdinPipe()
	ffmpegOut, _ := ffmpeg.StdoutPipe()
	if err := ffmpeg.Start(); err != nil {
		panic(err)
	}
	return ffmpegIn, ffmpegOut
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

func WriteCameraOutputToMplayerAndFfmpeg(droneVideoOutput chan []byte, mPlayerIn io.WriteCloser, ffmpegIn io.WriteCloser) {
	frameCount := 0
	for frame := range droneVideoOutput {
		if _, err := mPlayerIn.Write(frame); err != nil {
			fmt.Printf("failed to write frame to movie player: %v\n", err)
		}
		if (frameCount > 100) && ((frameCount % 10) == 0) {
			if _, err := ffmpegIn.Write(frame); err != nil {
				fmt.Printf("failed to write frame to ffmpeg: %v\n", err)
			}
		}
		frameCount++
	}
}
