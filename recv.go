package gondi

import (
	"errors"
	"unsafe"
)

// Allocate a new Receiver, using a NewRecvInstanceSetting struct as parameters
func NewRecvInstance(settings *NewRecvInstanceSettings) (*RecvInstance, error) {
	assertLibrary()

	var name *byte
	if settings.Name == "" {
		name = nil
	} else {
		name = cString(settings.Name)
	}

	intSettings := &recvCreateSettings{
		sourceToConnectTo: *settings.SourceToConnectTo,
		colorFormat:       settings.ColorFormat,
		bandwidth:         settings.Bandwidth,
		allowVideoFields:  settings.AllowVideoFields,
		name:              name,
	}

	inst := &RecvInstance{
		ndiInstance:    0,
		createSettings: intSettings,
	}

	inst.ndiInstance = ndilib_recv_create_v3(uintptr(unsafe.Pointer(intSettings)))
	if inst.ndiInstance == 0 {
		return nil, errors.New("unable to create receiver instance")
	}

	return inst, nil
}

// This will allow you to receive video, audio and metadata frames from the source you are connected to.
// Any of the frame pointers can be nil, in which case that type of frame will not be captured.
// This call can be called on separate threads, so it is possible to have a separate thread for each of video, audio and metadata.
// This function will return the type of frame that was received, or gondi.FrameTypeNone if no frame was received within the specified timeout.
func (p *RecvInstance) CaptureV2(vf *VideoFrameV2, af *AudioFrameV2, mf *MetadataFrame, timeoutMs uint32) FrameType {
	assertLibrary()

	return FrameType(ndilib_recv_capture_v2(p.ndiInstance, uintptr(unsafe.Pointer(vf)), uintptr(unsafe.Pointer(af)), uintptr(unsafe.Pointer(mf)), timeoutMs))
}

// Get the current amount of total and dropped video, audio and metadata frames. This can be used to determine if
// you have been calling instace.CaptureV2() fast enough to keep up with the incoming stream.
func (p *RecvInstance) GetPerformance() (total *RecvPerformance, dropped *RecvPerformance) {
	assertLibrary()
	total = &RecvPerformance{}
	dropped = &RecvPerformance{}

	ndilib_recv_get_performance(p.ndiInstance, uintptr(unsafe.Pointer(total)), uintptr(unsafe.Pointer(dropped)))

	return total, dropped
}

// Set the up-stream tally notifications. This returns FALSE if we are not currently connected to anything. That
// said, the moment that we do connect to something it will automatically be sent the tally state.
func (p *RecvInstance) SetTally(program bool, preview bool) bool {
	assertLibrary()
	tally := &Tally{program, preview}

	return ndilib_recv_set_tally(p.ndiInstance, uintptr(unsafe.Pointer(tally)))
}

// This function will send a meta frame to the source that we are connected too. This returns FALSE if we are
// not currently connected to anything.
func (p *RecvInstance) SendMetadata(metadata *MetadataFrame) bool {
	assertLibrary()

	return ndilib_recv_send_metadata(p.ndiInstance, uintptr(unsafe.Pointer(metadata)))
}

// Add a connection metadata string to the list of what is sent on each new connection. If someone is already connected then
// this frame will be sent to them immediately.
func (p *RecvInstance) AddConnectionMetadata(metadata *MetadataFrame) {
	assertLibrary()

	ndilib_recv_add_connection_metadata(p.ndiInstance, uintptr(unsafe.Pointer(metadata)))
}

// Connection based metadata is data that is sent automatically each time a new connection is received. You queue all of these
// up and they are sent on each connection. To reset them you need to clear them all and set them up again.
func (p *RecvInstance) ClearConnectionMetadata() {
	assertLibrary()

	ndilib_recv_clear_connection_metadata(p.ndiInstance)
}

// Free the buffers returned by capture for metadata
func (p *RecvInstance) FreeMetadata(metadata *MetadataFrame) {
	ndilib_recv_free_metadata(p.ndiInstance, uintptr(unsafe.Pointer(metadata)))
}

// Free the buffers returned by capture for video
func (p *RecvInstance) FreeVideoV2(vf *VideoFrameV2) {
	ndilib_recv_free_video_v2(p.ndiInstance, uintptr(unsafe.Pointer(vf)))
}

// Free the buffers returned by capture for audio
func (p *RecvInstance) FreeAudioV2(af *AudioFrameV2) {
	ndilib_recv_free_audio_v2(p.ndiInstance, uintptr(unsafe.Pointer(af)))
}

// Destroy a receiver instance
func (p *RecvInstance) Destroy() {
	assertLibrary()

	ndilib_recv_destroy(p.ndiInstance)
}
