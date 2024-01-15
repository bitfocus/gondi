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
			videoFrame := gondi.NewVideoFrameV2()
			videoFrame.FourCC = gondi.FourCCTypeRGBA
			videoFrame.FrameFormatType = gondi.FrameFormatProgressive
			videoFrame.Xres = int32(video.Width())
			videoFrame.Yres = int32(video.Height())
			videoFrame.LineStride = 0 // 2 bytes per pixel
			videoFrame.FrameRateN = 30000
			videoFrame.FrameRateD = 1001
			videoFrame.Data = &video.FrameBuffer()[0]
			sender.SendVideoFrame(videoFrame)
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
