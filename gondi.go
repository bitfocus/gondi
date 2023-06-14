/*
Package gondi provides a wrapper for the NDI SDK for Linux, macOS, using purego.

The NDI SDK is available at https://www.ndi.tv/sdk/

The NDI SDK is licensed under the NewTek SDK EULA, which can be found at https://www.ndi.tv/sdk/

NDI® is a registered trademark of Vizrt Group.
*/
package gondi

import (
	"unsafe"
)

// Get the version of the NDI library as string
func GetVersion() string {
	assertLibrary()

	mystrptr := ndilib_version()
	if mystrptr == 0 {
		return "N/A"
	}

	return goString(mystrptr)
}

// If your audio frames are interleaved, you can use this function to convert them to planar format.
// You can also use the frame.SetFromInterleavedArray(data) function to automatically convert an array of float32s in interleaved format to planar format as NDI likes it.
func ConvertAudioFromInterleaved(pSrc *AudioFrameV2, pDst *AudioFrameV2) {
	assertLibrary()

	ndilib_util_audio_from_interleaved_32f_v2(uintptr(unsafe.Pointer(pSrc)), uintptr(unsafe.Pointer(pDst)))
}

// If your want your audio frames to be interleaved, you can use this function to convert them from planar format.
// You can also use the frame.GetInterleavedArray() function to get a coverted array of float32s in interleaved format.
func ConvertAudioToInterleaved(pSrc *AudioFrameV2, pDst *AudioFrameV2) {
	assertLibrary()

	ndilib_util_audio_to_interleaved_32f_v2(uintptr(unsafe.Pointer(pSrc)), uintptr(unsafe.Pointer(pDst)))
}

// Allocate a new NDIMetadataFrame and initialize it with the specified utf-8 data string.
func NewMetadataFrame(data string) *MetadataFrame {
	// I am afraid that the data parameter might be garbage collected before the C code is done with it, though.
	ret := &MetadataFrame{
		Length:   int32(len(data)),
		Timecode: SendTimecodeSynthesize,
		Data:     cString(data),
	}

	return ret
}

// Get the data of the metadata frame as a utf8-string
func (p *MetadataFrame) GetData() string {
	return goString(uintptr(unsafe.Pointer(p.Data)))
}

// Name of the source
func (s *Source) Name() string {
	if s.name == nil {
		return ""
	}
	return goString(uintptr(unsafe.Pointer(s.name)))
}

// Address of the source
func (s *Source) Address() string {
	if s.address == nil {
		return ""
	}
	return goString(uintptr(unsafe.Pointer(s.address)))
}

// Set the name and address of the source object
func (s *Source) Set(name string, address string) {
	s.name = cString(name)
	s.address = cString(address)
}

// Get the audio frames as an array of float32
// This is usually stored as planar audio, so the first NumSamples values are the first channel, the next NumSamples values are the second channel, etc.
// If you need to work with interleaved audio, you can use the GetInterleavedArray() function instead.
func (p *AudioFrameV2) GetArray() []float32 {
	return (*[1 << 30]float32)(unsafe.Pointer(p.Data))[0 : p.NumSamples*p.NumChannels]
}

// Get the audio frames as an array of float32
// This function converts the audio to interleaved audio, so each sample is stored as a single value, and the channels are interleaved.
func (p *AudioFrameV2) GetInterleavedArray() []float32 {
	dst := make([]float32, p.NumSamples*p.NumChannels)
	tempFrame := &AudioFrameV2{
		NumSamples:  p.NumSamples,
		SampleRate:  p.SampleRate,
		NumChannels: p.NumChannels,
		Data:        &dst[0],
	}

	ndilib_util_audio_to_interleaved_32f_v2(uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(tempFrame)))

	return dst
}

// This function converts the interleaved audio from the parameter, to planar audio, and stores it in the Data field of the AudioFrameV2.
// The Data field of the frame needs to be preallocated.
func (p *AudioFrameV2) SetFromInterleavedArray(audio []float32) {
	if p.Data == nil {
		panic("AudioFrameV2.Data is nil")
	}
	tempFrame := &AudioFrameV2{
		NumSamples:  p.NumSamples,
		SampleRate:  p.SampleRate,
		NumChannels: p.NumChannels,
		Data:        &audio[0],
	}
	ndilib_util_audio_from_interleaved_32f_v2(uintptr(unsafe.Pointer(tempFrame)), uintptr(unsafe.Pointer(p)))
}

// Set the audio frames from an array of float32
// The Data field of the frame needs to be preallocated.
func (p *AudioFrameV2) SetArray(audio []float32) {
	if p.Data == nil {
		panic("AudioFrameV2.Data is nil")
	}
	copy((*[1 << 30]float32)(unsafe.Pointer(p.Data))[0:p.NumSamples*p.NumChannels], audio)
}
