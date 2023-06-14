package gondi

import (
	"errors"
	"unsafe"
)

// Setup a routed destination, specified by name and groups.
// The groups property may be empty, and it will use the default from NDI access manager.
func NewRoutingInstance(name string, groups string) (*RoutingInstance, error) {
	assertLibrary()

	settings := &routingCreateSettings{
		name:   cString(name),
		groups: cString(groups),
	}

	if groups == "" {
		settings.groups = nil
	}

	inst := ndilib_routing_create(uintptr(unsafe.Pointer(settings)))
	if inst == 0 {
		return nil, errors.New("unable to create routing instance")
	}

	instance := &RoutingInstance{inst, settings, name, groups}

	return instance, nil
}

// Get the set name of this instance. To change it, destroy this routing instance and create a new with the correct settings.
func (p *RoutingInstance) Name() string {
	return p.name
}

// Get the set groups for this instance, to change it, destroy this routing instance and create a new with the correct settings.
func (p *RoutingInstance) Groups() string {
	return p.groups
}

// Change the source this routing instance is connected to.
func (p *RoutingInstance) Change(source *Source) {
	assertLibrary()

	ndilib_routing_change(p.ndiInstance, uintptr(unsafe.Pointer(source)))
}

// Clear the current source this routing instance is connected to. Should return black to watchers.
func (p *RoutingInstance) Clear() {
	assertLibrary()

	ndilib_routing_clear(p.ndiInstance)
}

// Destroy this routing instance.
func (p *RoutingInstance) Destroy() {
	assertLibrary()

	ndilib_routing_destroy(p.ndiInstance)
}
