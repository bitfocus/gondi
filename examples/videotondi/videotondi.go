package main

import (
	"fmt"
	"log"
	"time"

	vidio "github.com/AlexEidt/Vidio"
	"github.com/bitfocus/gondi"
)

var (
	videoFrames int64 = 0
)

func videoToNDI(sender *gondi.SendInstance) {
	for {
		video, err := vidio.NewVideo("koala.mp4")
		if err != nil {
			log.Println("videoToNDI: failed to read video", err)
			continue
		}

		for video.Read() {
			startFrameTime := time.Now()
			videoFrame := gondi.NewVideoFrameV2()
			videoFrame.FourCC = gondi.FourCCTypeRGBA
			videoFrame.FrameFormatType = gondi.FrameFormatProgressive
			videoFrame.Xres = int32(video.Width())
			videoFrame.Yres = int32(video.Height())
			videoFrame.LineStride = 0 // 2 bytes per pixel
			videoFrame.FrameRateN = 30000
			videoFrame.FrameRateD = 1001
			videoFrame.Data = &video.FrameBuffer()[0]
			sender.SendVideoFrameAsync(videoFrame)

			// force frame lock to the video FPS
			// https://stackoverflow.com/a/61878644
			elapsed := time.Since(startFrameTime)
			fps := time.Duration((1.0 / video.FPS()) * float64(time.Second))
			if elapsed.Nanoseconds() < fps.Nanoseconds() {
				syncTime := int(fps.Nanoseconds() - elapsed.Nanoseconds())
				time.Sleep(time.Duration(syncTime) * time.Nanosecond)
			}
		}
	}
}

func main() {
	gondi.InitLibrary("")

	// Set up sender, block on both audio and video as we are using separate threads for audio and video
	sender, err := gondi.NewSendInstance("Output 1", "", true, true)
	if err != nil {
		panic(err)
	}
	defer sender.Destroy()

	// Set up video goroutine
	go videoToNDI(sender)

	// Show info
	for {
		fmt.Printf("Generated %d video frames \n", videoFrames)
		time.Sleep(1 * time.Second)
	}
}
