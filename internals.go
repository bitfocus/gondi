package gondi

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/ebitengine/purego"
)

var (
	ndi_shared_library uintptr
	ndilib_load        func() uintptr
	ndilib_initialize  func() bool
	ndilib_version     func() uintptr

	ndilib_util_audio_from_interleaved_32f_v2 func(src uintptr, dst uintptr)
	ndilib_util_audio_to_interleaved_32f_v2   func(src uintptr, dst uintptr)

	ndilib_send_create_v2                 func(settings uintptr) uintptr
	ndilib_send_destroy                   func(instance uintptr)
	ndilib_send_send_video_v2             func(instance uintptr, frame uintptr)
	ndilib_send_send_video_async_v2       func(instance uintptr, frame uintptr)
	ndilib_send_send_audio_v2             func(instance uintptr, frame uintptr)
	ndilib_send_send_metadata             func(instance uintptr, frame uintptr)
	ndilib_send_get_tally                 func(instance uintptr, tally uintptr, timeout uint32) bool
	ndilib_send_capture                   func(instance uintptr, metadata uintptr, timeout uint32) int32
	ndilib_send_free_metadata             func(instance uintptr, metadata uintptr)
	ndilib_send_add_connection_metadata   func(instance uintptr, metadata uintptr)
	ndilib_send_clear_connection_metadata func(instance uintptr)
	ndilib_send_set_failover              func(instance uintptr, source uintptr)
	ndilib_send_get_no_connections        func(instance uintptr, timeout uint32) int32

	ndilib_find_create_v2           func(settings uintptr) uintptr
	ndilib_find_destroy             func(instance uintptr)
	ndilib_find_get_current_sources func(instance uintptr, numSources uintptr) uintptr
	ndilib_find_wait_for_sources    func(instance uintptr, timeout uint32) bool

	ndilib_recv_create_v3                 func(settings uintptr) uintptr
	ndilib_recv_destroy                   func(instance uintptr)
	ndilib_recv_free_video_v2             func(instance uintptr, frame uintptr)
	ndilib_recv_free_audio_v2             func(instance uintptr, frame uintptr)
	ndilib_recv_free_metadata             func(instance uintptr, frame uintptr)
	ndilib_recv_capture_v2                func(instance uintptr, videoFrame uintptr, audioFrame uintptr, metadataFrame uintptr, timeout uint32) int32
	ndilib_recv_get_performance           func(instance uintptr, total uintptr, dropped uintptr)
	ndilib_recv_set_tally                 func(instance uintptr, tally uintptr) bool
	ndilib_recv_send_metadata             func(instance uintptr, metadata uintptr) bool
	ndilib_recv_add_connection_metadata   func(instance uintptr, metadata uintptr) bool
	ndilib_recv_clear_connection_metadata func(instance uintptr)

	ndilib_routing_create  func(settings uintptr) uintptr
	ndilib_routing_destroy func(instance uintptr)
	ndilib_routing_change  func(instance uintptr, source uintptr) bool
	ndilib_routing_clear   func(instance uintptr) bool
)

// Windows is not supported by go-purego
func getLibraryPath() string {
	switch runtime.GOOS {
	case "darwin":
		return "/usr/local/lib/libndi.dylib"
	case "linux":
		return "libndi.so"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}

func assertLibrary() {
	if ndi_shared_library == 0 {
		panic("library not initialized, use gondi.InitLibrary()")
	}
}

// Initialize the NDI Library, the libraryPath argument is optional and will be
// automatically set to the default library path for the current platform if
// empty. This function will panic if it does not find the library, or if it is
// unable to initialize NDI with the given library. But return error if NDI reports an error initializing.
func InitLibrary(libraryPath string) error {
	var err error

	if ndi_shared_library == 0 {
		if libraryPath == "" {
			libraryPath = getLibraryPath()
		}

		ndi_shared_library, err = purego.Dlopen(libraryPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			panic(err)
		}

		// Register all used NDI Library functions
		purego.RegisterLibFunc(&ndilib_load, ndi_shared_library, "NDIlib_v3_load")
		purego.RegisterLibFunc(&ndilib_initialize, ndi_shared_library, "NDIlib_initialize")
		purego.RegisterLibFunc(&ndilib_version, ndi_shared_library, "NDIlib_version")

		purego.RegisterLibFunc(&ndilib_util_audio_from_interleaved_32f_v2, ndi_shared_library, "NDIlib_util_audio_from_interleaved_32f_v2")
		purego.RegisterLibFunc(&ndilib_util_audio_to_interleaved_32f_v2, ndi_shared_library, "NDIlib_util_audio_to_interleaved_32f_v2")

		purego.RegisterLibFunc(&ndilib_send_create_v2, ndi_shared_library, "NDIlib_send_create_v2")
		purego.RegisterLibFunc(&ndilib_send_destroy, ndi_shared_library, "NDIlib_send_destroy")
		purego.RegisterLibFunc(&ndilib_send_send_video_v2, ndi_shared_library, "NDIlib_send_send_video_v2")
		purego.RegisterLibFunc(&ndilib_send_send_video_async_v2, ndi_shared_library, "NDIlib_send_send_video_async_v2")
		purego.RegisterLibFunc(&ndilib_send_send_audio_v2, ndi_shared_library, "NDIlib_send_send_audio_v2")
		purego.RegisterLibFunc(&ndilib_send_get_tally, ndi_shared_library, "NDIlib_send_get_tally")
		purego.RegisterLibFunc(&ndilib_send_capture, ndi_shared_library, "NDIlib_send_capture")
		purego.RegisterLibFunc(&ndilib_send_free_metadata, ndi_shared_library, "NDIlib_send_free_metadata")
		purego.RegisterLibFunc(&ndilib_send_send_metadata, ndi_shared_library, "NDIlib_send_send_metadata")
		purego.RegisterLibFunc(&ndilib_send_add_connection_metadata, ndi_shared_library, "NDIlib_send_add_connection_metadata")
		purego.RegisterLibFunc(&ndilib_send_clear_connection_metadata, ndi_shared_library, "NDIlib_send_clear_connection_metadata")
		purego.RegisterLibFunc(&ndilib_send_set_failover, ndi_shared_library, "NDIlib_send_set_failover")
		purego.RegisterLibFunc(&ndilib_send_get_no_connections, ndi_shared_library, "NDIlib_send_get_no_connections")

		purego.RegisterLibFunc(&ndilib_find_create_v2, ndi_shared_library, "NDIlib_find_create_v2")
		purego.RegisterLibFunc(&ndilib_find_get_current_sources, ndi_shared_library, "NDIlib_find_get_current_sources")
		purego.RegisterLibFunc(&ndilib_find_wait_for_sources, ndi_shared_library, "NDIlib_find_wait_for_sources")
		purego.RegisterLibFunc(&ndilib_find_destroy, ndi_shared_library, "NDIlib_find_destroy")

		purego.RegisterLibFunc(&ndilib_recv_create_v3, ndi_shared_library, "NDIlib_recv_create_v3")
		purego.RegisterLibFunc(&ndilib_recv_destroy, ndi_shared_library, "NDIlib_recv_destroy")
		purego.RegisterLibFunc(&ndilib_recv_free_metadata, ndi_shared_library, "NDIlib_recv_free_metadata")
		purego.RegisterLibFunc(&ndilib_recv_free_video_v2, ndi_shared_library, "NDIlib_recv_free_video_v2")
		purego.RegisterLibFunc(&ndilib_recv_free_audio_v2, ndi_shared_library, "NDIlib_recv_free_audio_v2")
		purego.RegisterLibFunc(&ndilib_recv_capture_v2, ndi_shared_library, "NDIlib_recv_capture_v2")
		purego.RegisterLibFunc(&ndilib_recv_get_performance, ndi_shared_library, "NDIlib_recv_get_performance")
		purego.RegisterLibFunc(&ndilib_recv_set_tally, ndi_shared_library, "NDIlib_recv_set_tally")
		purego.RegisterLibFunc(&ndilib_recv_send_metadata, ndi_shared_library, "NDIlib_recv_send_metadata")
		purego.RegisterLibFunc(&ndilib_recv_add_connection_metadata, ndi_shared_library, "NDIlib_recv_add_connection_metadata")
		purego.RegisterLibFunc(&ndilib_recv_clear_connection_metadata, ndi_shared_library, "NDIlib_recv_clear_connection_metadata")

		purego.RegisterLibFunc(&ndilib_routing_create, ndi_shared_library, "NDIlib_routing_create")
		purego.RegisterLibFunc(&ndilib_routing_destroy, ndi_shared_library, "NDIlib_routing_destroy")
		purego.RegisterLibFunc(&ndilib_routing_change, ndi_shared_library, "NDIlib_routing_change")
		purego.RegisterLibFunc(&ndilib_routing_clear, ndi_shared_library, "NDIlib_routing_clear")

		result := ndilib_load()
		if result == 0 {
			return errors.New("the NDIlib_v3_load function did not return a valid pointer")
		}

		loaded := ndilib_initialize()
		if !loaded {
			return errors.New("the NDIlib_initialize function returned false")
		}
	}

	return nil
}
