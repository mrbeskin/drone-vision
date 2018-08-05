package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func GetCamStream(driver *tello.Driver) chan []byte {
	videoStream := make(chan []byte)

	// register camera connection
	driver.On(tello.ConnectedEvent, func(data interface{}) {
		fmt.Println("Drone: camera connected")
		driver.StartVideo()
		driver.SetVideoEncoderRate(2)
		gobot.Every(2000*time.Millisecond, func() {
			driver.StartVideo()
		})
	})

	// register camera feed
	driver.On(tello.VideoFrameEvent, func(data interface{}) {
		pkt := data.([]byte)
		printFrameForDebug(pkt)
		videoStream <- pkt
	})
	return videoStream
}

var i = 0

func printFrameForDebug(b []byte) {
	if (i % 67) == 0 {
		fmt.Printf("length of video buffer: %d", len(b))
	}
	if (i % 199) == 0 {
		writeValueToTmpFile(b)
		fmt.Println(b)
	}
	i++

}

func writeValueToTmpFile(b []byte) {
	f, err := ioutil.TempFile("./", "tempy")
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(f.Name(), b, 0777)
}
