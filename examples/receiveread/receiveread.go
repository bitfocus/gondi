package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"time"
	"unsafe"

	"github.com/bitfocus/gondi"
	"github.com/nsmith5/mjpeg"
)

var PreviewIMG *image.RGBA

func captureAndReadVideo(receiver *gondi.RecvInstance, sender *gondi.SendInstance) {
	for {
		videoInput := gondi.NewVideoFrameV2()
		frametype := receiver.CaptureV2(videoInput, nil, nil, 1000)
		if frametype == gondi.FrameTypeVideo {
			// ... Manipulate video frame here ...
			size := videoInput.LineStride * videoInput.Yres
			videoInputSlice := unsafe.Slice(videoInput.Data, size)
			frame := make([]byte, len(videoInputSlice))
			copy(frame, videoInputSlice)

			// could not find a way to read these into a jpeg :/
			// if videoInput.FourCC == gondi.FourCCTypeUYVY {
			// 	log.Println("preview UYVY", len(frame))
			// }
			// if videoInput.FourCC == gondi.FourCCTypeBGRA {
			// 	log.Println("preview BGRA", len(frame))
			// }
			// if videoInput.FourCC == gondi.FourCCTypeUYVA {
			// 	log.Println("preview UYVA", len(frame))
			// }
			// if videoInput.FourCC == gondi.FourCCTypeBGRX {
			// 	log.Println("preview BGRX", len(frame))
			// }
			if videoInput.FourCC == gondi.FourCCTypeRGBA {
				// log.Println("preview RGBA", stream.Name, len(frame), int(videoInput.Xres), int(videoInput.Yres))
				PreviewIMG.Pix = frame
			}
			if videoInput.FourCC == gondi.FourCCTypeRGBX {
				// log.Println("preview RGBX", stream.Name, len(frame))
				PreviewIMG.Pix = frame
			}

			sender.SendVideoFrame(videoInput)
			receiver.FreeVideoV2(videoInput)
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
		ColorFormat:       gondi.RecvColorFormatRGBXRGBA, // this is the only format I could read on golang to a jpeg
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
	go captureAndReadVideo(receiver, sender)

	// Show info
	go func() {
		for {
			totals, dropped := receiver.GetPerformance()
			fmt.Printf("Total video frames received: %d, total dropped: %d\n", totals.VideoFrames, dropped.VideoFrames)
			time.Sleep(1 * time.Second)
		}
	}()

	stream := mjpeg.Handler{
		Next: func() (image.Image, error) {
			auxIMG := image.NewRGBA(image.Rect(0, 0, PreviewIMG.Bounds().Dx(), PreviewIMG.Bounds().Dy()))
			copy(auxIMG.Pix, PreviewIMG.Pix)
			return auxIMG, nil
		},
		Options: &jpeg.Options{Quality: 80},
	}

	mux := http.NewServeMux()
	mux.Handle("/stream", stream)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
