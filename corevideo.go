package darwin

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

type CVDisplayLinkRef uintptr
type CGDirectDisplayID uint32

const (
	KCVReturnSuccess = 0
)

var (
	_CVDisplayLinkCreateWithCGDisplay,
	_CVDisplayLinkSetOutputCallback,
	_CVDisplayLinkSetCurrentCGDisplay,
	_CVDisplayLinkStart,
	_CVDisplayLinkStop,
	_CVDisplayLinkRelease uintptr
)

func CVDisplayLinkCreateWithCGDisplay(displayID CGDirectDisplayID, displayLinkOut *CVDisplayLinkRef) int32 {
	ret, _, _ := purego.SyscallN(_CVDisplayLinkCreateWithCGDisplay, uintptr(displayID), uintptr(unsafe.Pointer(displayLinkOut)))
	return int32(ret)
}

func CVDisplayLinkSetOutputCallback(displayLink CVDisplayLinkRef, callback uintptr, userInfo unsafe.Pointer) int32 {
	ret, _, _ := purego.SyscallN(_CVDisplayLinkSetOutputCallback, uintptr(displayLink), callback, uintptr(userInfo))
	return int32(ret)
}

func CVDisplayLinkSetCurrentCGDisplay(displayLink CVDisplayLinkRef, displayID CGDirectDisplayID) int32 {
	ret, _, _ := purego.SyscallN(_CVDisplayLinkSetCurrentCGDisplay, uintptr(displayLink), uintptr(displayID))
	return int32(ret)
}

func CVDisplayLinkStart(displayLink CVDisplayLinkRef) int32 {
	ret, _, _ := purego.SyscallN(_CVDisplayLinkStart, uintptr(displayLink))
	return int32(ret)
}

func CVDisplayLinkStop(displayLink CVDisplayLinkRef) int32 {
	ret, _, _ := purego.SyscallN(_CVDisplayLinkStop, uintptr(displayLink))
	return int32(ret)
}

func CVDisplayLinkRelease(displayLink CVDisplayLinkRef) {
	purego.SyscallN(_CVDisplayLinkRelease, uintptr(displayLink))
}
