package main

import (
	"fmt"
	"math"
	"time"
	"unsafe"

	"github.com/bitfocus/gondi"
)

const (
	width  int32 = 1920
	height int32 = 1080
	fps    int32 = 50
)

var (
	videoFrames int64 = 0
	audioFrames int64 = 0
)

func castToBytes[T any](s []T) []byte {
	if len(s) == 0 {
		return nil
	}

	size := unsafe.Sizeof(s[0])
	return unsafe.Slice((*byte)(unsafe.Pointer(&s[0])), int(size)*len(s))
}

func generateTestVideo(sender *gondi.SendInstance, use75Bars bool, movingBars bool) {
	imageBuffer := make([]uint32, width*height)
	byteBuffer := castToBytes(imageBuffer)

	bars100 := [8]uint32{0xEB80EB80, 0xD292D210, 0xAA10AAA6, 0x91229136, 0x6ADE6ACA, 0x51F0515A, 0x296E29F0, 0x10801080}
	bars75 := [8]uint32{0xB480B480, 0xA888A82C, 0x912C9193, 0x8534853F, 0x3FCC3FC1, 0x33D4336D, 0x1C781CD4, 0x10801080}

	videoFrame := gondi.NewVideoFrameV2()
	videoFrame.FourCC = gondi.FourCCTypeUYVY
	videoFrame.FrameFormatType = gondi.FrameFormatProgressive
	videoFrame.Xres = width
	videoFrame.Yres = height
	videoFrame.LineStride = width * 2 // 2 bytes per pixel
	videoFrame.FrameRateN = fps
	videoFrame.FrameRateD = 1
	videoFrame.Timecode = gondi.SendTimecodeSynthesize
	videoFrame.Timestamp = gondi.SendTimecodeEmpty
	videoFrame.PictureAspectRatio = 0
	videoFrame.Data = &byteBuffer[0]

	var j int32 = 0
	for {
		// Generate video frame
		var (
			i    = 0
			x, y int32
		)

		if use75Bars {
			for y = 0; y < height; y++ {
				for x = 0; x < width; x += 2 {
					imageBuffer[i] = bars75[(((j+x)%width)*8)/width]
					i++
				}
			}
		} else {
			for y = 0; y < height; y++ {
				for x = 0; x < width; x += 2 {
					imageBuffer[i] = bars100[(((j+x)%width)*8)/width]
					i++
				}
			}
		}

		// Will make sure only 50 frames per second are sent because of the video sync setting in NewSendInstance()
		sender.SendVideoFrame(videoFrame)
		videoFrames++
		if movingBars {
			j += 2
		}
	}
}

// Generate a 2 channel 1KHz sine wave
func generateTestAudio(sender *gondi.SendInstance, dbVolume float32) {
	if (dbVolume > 0) || (dbVolume < -144) {
		panic("Volume must be between 0dB and -144dB")
	}

	volume := float32(math.Pow(10, float64(dbVolume/20)))

	audioFrame := gondi.NewAudioFrameV2()
	audioFrame.SampleRate = 48000
	audioFrame.NumChannels = 2
	audioFrame.NumSamples = int32(48000 / fps)
	audioBuffer := make([]float32, audioFrame.SampleRate*audioFrame.NumSamples)
	audioFrame.Data = &audioBuffer[0]
	audioFrame.ChannelStride = audioFrame.NumSamples * int32(unsafe.Sizeof(audioBuffer[0]))

	// ii is used to keep track of the continous sample number across frames
	var i, ii int32
	var frequency float64 = 1000

	sampletime := 1 / float64(audioFrame.SampleRate)
	for {
		for i = 0; i < audioFrame.NumSamples; i++ {
			audioBuffer[i] = float32(math.Sin(2*math.Pi*frequency*float64(ii)*sampletime)) * volume
			// Add same sample to right channel too
			audioBuffer[i+audioFrame.NumSamples] = audioBuffer[i]
			ii = (ii + 1) % audioFrame.SampleRate
		}

		// Will send 50 samples per second because of the audio sync setting in NewSendInstance()
		// and the number of samples per frame sent
		sender.SendAudioFrame(audioFrame)
		audioFrames++
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

	// Set up video goroutine, set the last parameter to false if you do not want the bars to move
	go generateTestVideo(sender, false, true)

	// Set up audio goroutine generating 1Khz test tone, volume is specified in peak-to-peak dB. Mean volume is 3dB lower
	go generateTestAudio(sender, -20)

	// Show info
	for {
		fmt.Printf("Generated %d video frames and %d audio frames\n", videoFrames, audioFrames)
		time.Sleep(1 * time.Second)
	}
}
