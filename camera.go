package main

import (
	"fmt"
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
		videoStream <- pkt
	})
	return videoStream
}
