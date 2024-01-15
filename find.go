package gondi

import (
	"errors"
	"unsafe"
)

// Setup a finder instance, initialized groups and extraIPs.
// ShowLocalSources will control whether local sources are shown or not.
// The groups property may be empty, and it will use the default from NDI access manager.
// The extraIPs is only used to manually find sources from known ips on different subnets and is comma separated.
func NewFindInstance(showLocalSources bool, groups string, extraIPs string) (*FindInstance, error) {
	assertLibrary()
	settings := &findCreateSettings{showLocalSources, cString(groups), cString(extraIPs)}
	inst := &FindInstance{}

	if groups == "" {
		settings.groups = nil
	}
	if extraIPs == "" {
		settings.extraIPs = nil
	}

	inst.createSettings = settings
	inst.ndiInstance = ndilib_find_create_v2(uintptr(unsafe.Pointer(settings)))
	if inst.ndiInstance == 0 {
		return nil, errors.New("unable to create finder instance")
	}

	return inst, nil
}

// Get the current sources from this finder instance. It is recomended to call WaitForSources before this.
// If you have a UI element to change the source, you should call this function before showing the user the list of sources,
// to always have the latest list of sources.
func (p *FindInstance) GetCurrentSources() []*Source {
	assertLibrary()

	var numSources uint32
	ret := ndilib_find_get_current_sources(p.ndiInstance, uintptr(unsafe.Pointer(&numSources)))

	sources := make([]*Source, numSources)

	// We take the address and then dereference it to trick go vet from creating a possible misuse of unsafe.Pointer
	blockp := *(*unsafe.Pointer)(unsafe.Pointer(&ret))

	for i := range sources {
		sources[i] = (*Source)(blockp)
		// Increment pointer
		blockp = unsafe.Add(blockp, unsafe.Sizeof(Source{}))
	}

	return sources
}

// This allows you to wait until the sources on the network have changed.
// It will return true if the list of sourcess has changed within the timeout, false otherwise.
// You are not required to call this function, but it is helpful for getting an initial list of sources,
// and to detect when the list of sources has changed.
func (p *FindInstance) WaitForSources(timeoutMs uint32) bool {
	assertLibrary()

	return ndilib_find_wait_for_sources(p.ndiInstance, timeoutMs)
}

// Destroy this finder instance.
func (p *FindInstance) Destroy() {
	assertLibrary()

	ndilib_find_destroy(p.ndiInstance)
}
