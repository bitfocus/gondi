package gondi

import (
	"errors"
	"unsafe"
)

// Set up a sender instance using the specified name and string.
// Syncronous calls will block on either audio or video frames, or both, depending on the clockVideo and clockAudio parameters, to make sure that the frames are sent at the correct time.
func NewSendInstance(name string, groups string, clockVideo bool, clockAudio bool) (*SendInstance, error) {
	assertLibrary()

	settings := &sendCreateSettings{cString(name), cString(groups), clockVideo, clockAudio}
	instance := ndilib_send_create_v2(uintptr(unsafe.Pointer(settings)))
	if instance == 0 {
		return nil, errors.New("unable to create send instance")
	}

	return &SendInstance{instance, settings}, nil
}

// Remember to call Destroy() on the instance when you are done with it. This will free up resources and unregister the sender.
func (p *SendInstance) Destroy() error {
	assertLibrary()

	ndilib_send_destroy(p.ndiInstance)

	return nil
}

// Allocate a new NDI audio frame object
func NewAudioFrameV2() *AudioFrameV2 {
	af := &AudioFrameV2{}

	af.SampleRate = 0
	af.NumChannels = 0
	af.NumSamples = 0
	af.Timecode = SendTimecodeSynthesize
	af.Data = nil
	af.ChannelStride = 0
	af.Metadata = nil
	af.Timestamp = SendTimecodeEmpty

	return af
}

// Allocate a new NDI audio frame object with preallocated data for
// holding numChannels * numSamples samples.
func NewAudioFrameV2Preallocated(numChannels int32, numSamples int32) *AudioFrameV2 {
	af := &AudioFrameV2{}
	data := make([]float32, numChannels*numSamples)

	af.SampleRate = 0
	af.NumChannels = 0
	af.NumSamples = 0
	af.Timecode = SendTimecodeSynthesize
	af.Data = &data[0]
	af.ChannelStride = 0
	af.Metadata = nil
	af.Timestamp = SendTimecodeEmpty

	return af
}

// Allocate a new NDI video frame with defaults
func NewVideoFrameV2() *VideoFrameV2 {
	frame := &VideoFrameV2{}

	frame.Xres = 0
	frame.Yres = 0
	frame.FourCC = FourCCTypeBGRX
	frame.FrameRateN = 25
	frame.FrameRateD = 1
	frame.PictureAspectRatio = 0
	frame.FrameFormatType = FrameFormatProgressive
	frame.Timecode = SendTimecodeSynthesize
	frame.Data = nil
	frame.LineStride = 0
	frame.Metadata = nil
	frame.Timestamp = SendTimecodeEmpty

	return frame
}

// Send a video frame. This call is syncronous and will block until the frame has been sent if you specified clockVideo=true in NewNDISendInstance().
func (p *SendInstance) SendVideoFrame(frame *VideoFrameV2) {
	assertLibrary()

	ndilib_send_send_video_v2(p.ndiInstance, uintptr(unsafe.Pointer(frame)))
}

// Send video asynchronously, this call will return immediately, and you need to keep the video frame memory resident until a
// synchronizing event has been received.
// Syncronizing events are:
// - A call to frame.SendVideoFrame()
// - A call to frame.SendVideoFrameAsync() with a different video frame
// - A call to frame.SendVideoFrame(nil)
// - A call to frame.Destroy()
func (p *SendInstance) SendVideoFrameAsync(frame *VideoFrameV2) {
	assertLibrary()

	ndilib_send_send_video_async_v2(p.ndiInstance, uintptr(unsafe.Pointer(frame)))
}

// Send a metadata frame
func (p *SendInstance) SendMetadataFrame(frame *MetadataFrame) {
	assertLibrary()

	ndilib_send_send_metadata(p.ndiInstance, uintptr(unsafe.Pointer(frame)))
}

// This method lets you receive metadata from the other end of the connection.
// Remember that there might be multiple connections to your sender instance.
func (p *SendInstance) Capture(metadata *MetadataFrame, timeoutMs uint32) FrameType {
	assertLibrary()

	return FrameType(ndilib_send_capture(p.ndiInstance, uintptr(unsafe.Pointer(metadata)), timeoutMs))
}

// Add a connection metadata string to the list of what is sent on each new connection. If someone is already connected then
// this string will be sent to them immediately.
func (p *SendInstance) AddConnectionMetadata(metadata *MetadataFrame) {
	assertLibrary()

	ndilib_send_add_connection_metadata(p.ndiInstance, uintptr(unsafe.Pointer(metadata)))
}

// Connection based metadata is data that is sent automatically each time a new connection is received. You queue all of these
// up and they are sent on each connection. To reset them you need to clear them all and set them up again.
func (p *SendInstance) ClearConnectionMetadata() {
	assertLibrary()

	ndilib_send_clear_connection_metadata(p.ndiInstance)
}

// Get the current number of receivers connected to this source. This can be used to avoid even rendering when nothing is connected to the video source.
// which can significantly improve the efficiency if you want to make a lot of sources available on the network. If you specify a timeout that is not
// 0 then it will wait until there are connections for this amount of time.
func (p *SendInstance) GetNumberOfConnections(timeoutMs uint32) int32 {
	assertLibrary()

	return ndilib_send_get_no_connections(p.ndiInstance, timeoutMs)
}

// Determine the current tally sate. If you specify a timeout then it will wait until it has changed, otherwise it will simply poll it
// and return the current tally immediately. The boolean return value is whether anything has actually changed (true) or whether it timed out (false)
func (p *SendInstance) GetTally(timeoutMs uint32) (*Tally, bool) {
	assertLibrary()
	tally := &Tally{}

	changed := ndilib_send_get_tally(p.ndiInstance, uintptr(unsafe.Pointer(tally)), timeoutMs)

	return tally, changed
}

// Send an audio frame. This call is syncronous and will block until the frame has been sent, if you specified clockAudio=true in NewNDISendInstance().
func (p *SendInstance) SendAudioFrame(frame *AudioFrameV2) {
	assertLibrary()

	ndilib_send_send_audio_v2(p.ndiInstance, uintptr(unsafe.Pointer(frame)))
}

// This will assign a new fail-over source for this video source. What this means is that if this video source was to fail
// any receivers would automatically switch over to use this source, unless this source then came back online. You can specify
// nil to clear the source.
func (p *SendInstance) SetFailover(source *Source) {
	assertLibrary()

	ndilib_send_set_failover(p.ndiInstance, uintptr(unsafe.Pointer(source)))
}
