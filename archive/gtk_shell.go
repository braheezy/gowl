// Generated by go-wayland-scanner
// https://github.com/rajveermalviya/go-wayland/cmd/go-wayland-scanner
// XML file : https://raw.githubusercontent.com/bugaevc/wl-clipboard/master/src/protocol/gtk-shell.xml
//
// gtk Protocol Copyright:

package main

import "github.com/rajveermalviya/go-wayland/wayland/client"

// GtkShell1 : gtk specific extensions
//
// gtk_shell is a protocol extension providing additional features for
// clients implementing it.
type GtkShell1 struct {
	client.BaseProxy
	capabilitiesHandler GtkShell1CapabilitiesHandlerFunc
}

// NewGtkShell1 : gtk specific extensions
//
// gtk_shell is a protocol extension providing additional features for
// clients implementing it.
func NewGtkShell1(ctx *client.Context) *GtkShell1 {
	gtkShell1 := &GtkShell1{}
	ctx.Register(gtkShell1)
	return gtkShell1
}

// GetGtkSurface :
func (i *GtkShell1) GetGtkSurface(surface *client.Surface) (*GtkSurface1, error) {
	gtkSurface := NewGtkSurface1(i.Context())
	const opcode = 0
	const _reqBufLen = 8 + 4 + 4
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	client.PutUint32(_reqBuf[l:l+4], gtkSurface.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], surface.ID())
	l += 4
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return gtkSurface, err
}

// SetStartupId :
func (i *GtkShell1) SetStartupId(startupId string) error {
	const opcode = 1
	startupIdLen := client.PaddedLen(len(startupId) + 1)
	_reqBufLen := 8 + (4 + startupIdLen)
	_reqBuf := make([]byte, _reqBufLen)
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	client.PutString(_reqBuf[l:l+(4+startupIdLen)], startupId, startupIdLen)
	l += (4 + startupIdLen)
	err := i.Context().WriteMsg(_reqBuf, nil)
	return err
}

// SystemBell :
func (i *GtkShell1) SystemBell(surface *GtkSurface1) error {
	const opcode = 2
	const _reqBufLen = 8 + 4
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	if surface == nil {
		client.PutUint32(_reqBuf[l:l+4], 0)
		l += 4
	} else {
		client.PutUint32(_reqBuf[l:l+4], surface.ID())
		l += 4
	}
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return err
}

// NotifyLaunch :
func (i *GtkShell1) NotifyLaunch(startupId string) error {
	const opcode = 3
	startupIdLen := client.PaddedLen(len(startupId) + 1)
	_reqBufLen := 8 + (4 + startupIdLen)
	_reqBuf := make([]byte, _reqBufLen)
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	client.PutString(_reqBuf[l:l+(4+startupIdLen)], startupId, startupIdLen)
	l += (4 + startupIdLen)
	err := i.Context().WriteMsg(_reqBuf, nil)
	return err
}

func (i *GtkShell1) Destroy() error {
	i.Context().Unregister(i)
	return nil
}

type GtkShell1Capability uint32

// GtkShell1Capability :
const (
	GtkShell1CapabilityGlobalAppMenu GtkShell1Capability = 1
	GtkShell1CapabilityGlobalMenuBar GtkShell1Capability = 2
	GtkShell1CapabilityDesktopIcons  GtkShell1Capability = 3
)

func (e GtkShell1Capability) Name() string {
	switch e {
	case GtkShell1CapabilityGlobalAppMenu:
		return "global_app_menu"
	case GtkShell1CapabilityGlobalMenuBar:
		return "global_menu_bar"
	case GtkShell1CapabilityDesktopIcons:
		return "desktop_icons"
	default:
		return ""
	}
}

func (e GtkShell1Capability) Value() string {
	switch e {
	case GtkShell1CapabilityGlobalAppMenu:
		return "1"
	case GtkShell1CapabilityGlobalMenuBar:
		return "2"
	case GtkShell1CapabilityDesktopIcons:
		return "3"
	default:
		return ""
	}
}

func (e GtkShell1Capability) String() string {
	return e.Name() + "=" + e.Value()
}

// GtkShell1CapabilitiesEvent :
type GtkShell1CapabilitiesEvent struct {
	Capabilities uint32
}
type GtkShell1CapabilitiesHandlerFunc func(GtkShell1CapabilitiesEvent)

// SetCapabilitiesHandler : sets handler for GtkShell1CapabilitiesEvent
func (i *GtkShell1) SetCapabilitiesHandler(f GtkShell1CapabilitiesHandlerFunc) {
	i.capabilitiesHandler = f
}

func (i *GtkShell1) Dispatch(opcode uint32, fd int, data []byte) {
	switch opcode {
	case 0:
		if i.capabilitiesHandler == nil {
			return
		}
		var e GtkShell1CapabilitiesEvent
		l := 0
		e.Capabilities = client.Uint32(data[l : l+4])
		l += 4

		i.capabilitiesHandler(e)
	}
}

// GtkSurface1 :
type GtkSurface1 struct {
	client.BaseProxy
	configureHandler      GtkSurface1ConfigureHandlerFunc
	configureEdgesHandler GtkSurface1ConfigureEdgesHandlerFunc
}

// NewGtkSurface1 :
func NewGtkSurface1(ctx *client.Context) *GtkSurface1 {
	gtkSurface1 := &GtkSurface1{}
	ctx.Register(gtkSurface1)
	return gtkSurface1
}

// SetDbusProperties :
func (i *GtkSurface1) SetDbusProperties(applicationId, appMenuPath, menubarPath, windowObjectPath, applicationObjectPath, uniqueBusName string) error {
	const opcode = 0
	applicationIdLen := client.PaddedLen(len(applicationId) + 1)
	appMenuPathLen := client.PaddedLen(len(appMenuPath) + 1)
	menubarPathLen := client.PaddedLen(len(menubarPath) + 1)
	windowObjectPathLen := client.PaddedLen(len(windowObjectPath) + 1)
	applicationObjectPathLen := client.PaddedLen(len(applicationObjectPath) + 1)
	uniqueBusNameLen := client.PaddedLen(len(uniqueBusName) + 1)
	_reqBufLen := 8 + (4 + applicationIdLen) + (4 + appMenuPathLen) + (4 + menubarPathLen) + (4 + windowObjectPathLen) + (4 + applicationObjectPathLen) + (4 + uniqueBusNameLen)
	_reqBuf := make([]byte, _reqBufLen)
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	client.PutString(_reqBuf[l:l+(4+applicationIdLen)], applicationId, applicationIdLen)
	l += (4 + applicationIdLen)
	client.PutString(_reqBuf[l:l+(4+appMenuPathLen)], appMenuPath, appMenuPathLen)
	l += (4 + appMenuPathLen)
	client.PutString(_reqBuf[l:l+(4+menubarPathLen)], menubarPath, menubarPathLen)
	l += (4 + menubarPathLen)
	client.PutString(_reqBuf[l:l+(4+windowObjectPathLen)], windowObjectPath, windowObjectPathLen)
	l += (4 + windowObjectPathLen)
	client.PutString(_reqBuf[l:l+(4+applicationObjectPathLen)], applicationObjectPath, applicationObjectPathLen)
	l += (4 + applicationObjectPathLen)
	client.PutString(_reqBuf[l:l+(4+uniqueBusNameLen)], uniqueBusName, uniqueBusNameLen)
	l += (4 + uniqueBusNameLen)
	err := i.Context().WriteMsg(_reqBuf, nil)
	return err
}

// SetModal :
func (i *GtkSurface1) SetModal() error {
	const opcode = 1
	const _reqBufLen = 8
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return err
}

// UnsetModal :
func (i *GtkSurface1) UnsetModal() error {
	const opcode = 2
	const _reqBufLen = 8
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return err
}

// Present :
func (i *GtkSurface1) Present(time uint32) error {
	const opcode = 3
	const _reqBufLen = 8 + 4
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(time))
	l += 4
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return err
}

// RequestFocus :
func (i *GtkSurface1) RequestFocus(startupId string) error {
	const opcode = 4
	startupIdLen := client.PaddedLen(len(startupId) + 1)
	_reqBufLen := 8 + (4 + startupIdLen)
	_reqBuf := make([]byte, _reqBufLen)
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	client.PutString(_reqBuf[l:l+(4+startupIdLen)], startupId, startupIdLen)
	l += (4 + startupIdLen)
	err := i.Context().WriteMsg(_reqBuf, nil)
	return err
}

// Release :
func (i *GtkSurface1) Release() error {
	defer i.Context().Unregister(i)
	const opcode = 5
	const _reqBufLen = 8
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return err
}

// TitlebarGesture :
func (i *GtkSurface1) TitlebarGesture(serial uint32, seat *client.Seat, gesture uint32) error {
	const opcode = 6
	const _reqBufLen = 8 + 4 + 4 + 4
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(serial))
	l += 4
	client.PutUint32(_reqBuf[l:l+4], seat.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(gesture))
	l += 4
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return err
}

type GtkSurface1State uint32

// GtkSurface1State :
const (
	GtkSurface1StateTiled       GtkSurface1State = 1
	GtkSurface1StateTiledTop    GtkSurface1State = 2
	GtkSurface1StateTiledRight  GtkSurface1State = 3
	GtkSurface1StateTiledBottom GtkSurface1State = 4
	GtkSurface1StateTiledLeft   GtkSurface1State = 5
)

func (e GtkSurface1State) Name() string {
	switch e {
	case GtkSurface1StateTiled:
		return "tiled"
	case GtkSurface1StateTiledTop:
		return "tiled_top"
	case GtkSurface1StateTiledRight:
		return "tiled_right"
	case GtkSurface1StateTiledBottom:
		return "tiled_bottom"
	case GtkSurface1StateTiledLeft:
		return "tiled_left"
	default:
		return ""
	}
}

func (e GtkSurface1State) Value() string {
	switch e {
	case GtkSurface1StateTiled:
		return "1"
	case GtkSurface1StateTiledTop:
		return "2"
	case GtkSurface1StateTiledRight:
		return "3"
	case GtkSurface1StateTiledBottom:
		return "4"
	case GtkSurface1StateTiledLeft:
		return "5"
	default:
		return ""
	}
}

func (e GtkSurface1State) String() string {
	return e.Name() + "=" + e.Value()
}

type GtkSurface1EdgeConstraint uint32

// GtkSurface1EdgeConstraint :
const (
	GtkSurface1EdgeConstraintResizableTop    GtkSurface1EdgeConstraint = 1
	GtkSurface1EdgeConstraintResizableRight  GtkSurface1EdgeConstraint = 2
	GtkSurface1EdgeConstraintResizableBottom GtkSurface1EdgeConstraint = 3
	GtkSurface1EdgeConstraintResizableLeft   GtkSurface1EdgeConstraint = 4
)

func (e GtkSurface1EdgeConstraint) Name() string {
	switch e {
	case GtkSurface1EdgeConstraintResizableTop:
		return "resizable_top"
	case GtkSurface1EdgeConstraintResizableRight:
		return "resizable_right"
	case GtkSurface1EdgeConstraintResizableBottom:
		return "resizable_bottom"
	case GtkSurface1EdgeConstraintResizableLeft:
		return "resizable_left"
	default:
		return ""
	}
}

func (e GtkSurface1EdgeConstraint) Value() string {
	switch e {
	case GtkSurface1EdgeConstraintResizableTop:
		return "1"
	case GtkSurface1EdgeConstraintResizableRight:
		return "2"
	case GtkSurface1EdgeConstraintResizableBottom:
		return "3"
	case GtkSurface1EdgeConstraintResizableLeft:
		return "4"
	default:
		return ""
	}
}

func (e GtkSurface1EdgeConstraint) String() string {
	return e.Name() + "=" + e.Value()
}

type GtkSurface1Gesture uint32

// GtkSurface1Gesture :
const (
	GtkSurface1GestureDoubleClick GtkSurface1Gesture = 1
	GtkSurface1GestureRightClick  GtkSurface1Gesture = 2
	GtkSurface1GestureMiddleClick GtkSurface1Gesture = 3
)

func (e GtkSurface1Gesture) Name() string {
	switch e {
	case GtkSurface1GestureDoubleClick:
		return "double_click"
	case GtkSurface1GestureRightClick:
		return "right_click"
	case GtkSurface1GestureMiddleClick:
		return "middle_click"
	default:
		return ""
	}
}

func (e GtkSurface1Gesture) Value() string {
	switch e {
	case GtkSurface1GestureDoubleClick:
		return "1"
	case GtkSurface1GestureRightClick:
		return "2"
	case GtkSurface1GestureMiddleClick:
		return "3"
	default:
		return ""
	}
}

func (e GtkSurface1Gesture) String() string {
	return e.Name() + "=" + e.Value()
}

type GtkSurface1Error uint32

// GtkSurface1Error :
const (
	GtkSurface1ErrorInvalidGesture GtkSurface1Error = 0
)

func (e GtkSurface1Error) Name() string {
	switch e {
	case GtkSurface1ErrorInvalidGesture:
		return "invalid_gesture"
	default:
		return ""
	}
}

func (e GtkSurface1Error) Value() string {
	switch e {
	case GtkSurface1ErrorInvalidGesture:
		return "0"
	default:
		return ""
	}
}

func (e GtkSurface1Error) String() string {
	return e.Name() + "=" + e.Value()
}

// GtkSurface1ConfigureEvent :
type GtkSurface1ConfigureEvent struct {
	States []byte
}
type GtkSurface1ConfigureHandlerFunc func(GtkSurface1ConfigureEvent)

// SetConfigureHandler : sets handler for GtkSurface1ConfigureEvent
func (i *GtkSurface1) SetConfigureHandler(f GtkSurface1ConfigureHandlerFunc) {
	i.configureHandler = f
}

// GtkSurface1ConfigureEdgesEvent :
type GtkSurface1ConfigureEdgesEvent struct {
	Constraints []byte
}
type GtkSurface1ConfigureEdgesHandlerFunc func(GtkSurface1ConfigureEdgesEvent)

// SetConfigureEdgesHandler : sets handler for GtkSurface1ConfigureEdgesEvent
func (i *GtkSurface1) SetConfigureEdgesHandler(f GtkSurface1ConfigureEdgesHandlerFunc) {
	i.configureEdgesHandler = f
}

func (i *GtkSurface1) Dispatch(opcode uint32, fd int, data []byte) {
	switch opcode {
	case 0:
		if i.configureHandler == nil {
			return
		}
		var e GtkSurface1ConfigureEvent
		l := 0
		statesLen := int(client.Uint32(data[l : l+4]))
		l += 4
		e.States = make([]byte, statesLen)
		copy(e.States, data[l:l+statesLen])
		l += statesLen

		i.configureHandler(e)
	case 1:
		if i.configureEdgesHandler == nil {
			return
		}
		var e GtkSurface1ConfigureEdgesEvent
		l := 0
		constraintsLen := int(client.Uint32(data[l : l+4]))
		l += 4
		e.Constraints = make([]byte, constraintsLen)
		copy(e.Constraints, data[l:l+constraintsLen])
		l += constraintsLen

		i.configureEdgesHandler(e)
	}
}