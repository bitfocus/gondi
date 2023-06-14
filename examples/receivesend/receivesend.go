package main

import (
	"fmt"
	"time"

	"github.com/bitfocus/gondi"
)

func captureAndSendVideo(receiver *gondi.RecvInstance, sender *gondi.SendInstance) {
	for {
		videoInput := gondi.NewVideoFrameV2()
		frametype := receiver.CaptureV2(videoInput, nil, nil, 1000)
		if frametype == gondi.FrameTypeVideo {
			// ... Manipulate video frame here ...
			sender.SendVideoFrame(videoInput)
			receiver.FreeVideoV2(videoInput)
		}
	}
}

func captureAndSendAudio(receiver *gondi.RecvInstance, sender *gondi.SendInstance) {
	for {
		audioInput := gondi.NewAudioFrameV2()
		frametype := receiver.CaptureV2(nil, audioInput, nil, 1000)
		if frametype == gondi.FrameTypeAudio {
			// ... Manipulate audio frame here ...
			sender.SendAudioFrame(audioInput)
			receiver.FreeAudioV2(audioInput)
		}
	}
}

func main() {
	fmt.Println("Initializing NDI")
	gondi.InitLibrary("")

	version := gondi.GetVersion()
	fmt.Printf("NDI version: %s\n", version)

	findInstance, err := gondi.NewFindInstance(true, "", "10.20.10.42")
	if err != nil {
		panic(err)
	}
	defer findInstance.Destroy()

	// Wait for sources to appear
	fmt.Println("Looking for sources...")
	for {
		more := findInstance.WaitForSources(5000)
		if !more {
			break
		}
	}

	// Fetch the sources
	sources := findInstance.GetCurrentSources()

	if len(sources) == 0 {
		fmt.Println("No sources found, cannot continue")
		return
	}

	selectedSource := sources[0]
	fmt.Printf("Source selected: %s\n", selectedSource.Name())

	// Set up receiver
	receiver, err := gondi.NewRecvInstance(&gondi.NewRecvInstanceSettings{
		SourceToConnectTo: selectedSource,
		ColorFormat:       gondi.RecvColorFormatBGRXBGRA,
		Bandwidth:         gondi.RecvBandwidthHighest,
		AllowVideoFields:  true,
		Name:              "Receive 1",
	})
	if err != nil {
		panic(err)
	}
	defer receiver.Destroy()

	// Set up sender, block on both audio and video as we are using separate threads for audio and video
	sender, err := gondi.NewSendInstance("Output 1", "", true, true)
	if err != nil {
		panic(err)
	}
	defer sender.Destroy()

	// Set up threads
	go captureAndSendVideo(receiver, sender)
	go captureAndSendAudio(receiver, sender)

	// Show info
	for {
		totals, dropped := receiver.GetPerformance()
		fmt.Printf("Total video frames received: %d, total dropped: %d\n", totals.VideoFrames, dropped.VideoFrames)
		fmt.Printf("Total audio frames received: %d, total dropped: %d\n", totals.AudioFrames, dropped.AudioFrames)
		time.Sleep(1 * time.Second)
	}
}
