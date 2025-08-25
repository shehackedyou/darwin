package darwin

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

func EventLocationInWindow(event NSEvent) (float64, float64) {
	var loc NSPoint
	var locationInWindow func(receiver uintptr, selector Selector) NSPoint
	purego.RegisterLibFunc(&locationInWindow, libobjc, "objc_msgSend")
	loc = locationInWindow(uintptr(event.Ptr), Sel_locationInWindow)
	return loc.X, loc.Y
}

func EventScrollingDeltaX(event NSEvent) float64 {
	return Objc_sendMsg[float64](uintptr(event.Ptr), Sel_scrollingDeltaX)
}

func EventScrollingDeltaY(event NSEvent) float64 {
	return Objc_sendMsg[float64](uintptr(event.Ptr), Sel_scrollingDeltaY)
}

func EventButtonNumber(event NSEvent) int {
	return int(Objc_sendMsg[int64](uintptr(event.Ptr), Sel_buttonNumber))
}

func EventClickCount(event NSEvent) int {
	return int(Objc_sendMsg[int64](uintptr(event.Ptr), Sel_clickCount))
}

func EventKeyCode(event NSEvent) int {
	return int(Objc_sendMsg[uint16](uintptr(event.Ptr), Sel_keyCode))
}

func EventModifierFlags(event NSEvent) uintptr {
	return Objc_sendMsg[uintptr](uintptr(event.Ptr), Sel_modifierFlags)
}

func EventCharacters(event NSEvent) string {
	nsString := Objc_sendMsg[uintptr](uintptr(event.Ptr), Sel_characters)
	return (NSString{Object{unsafe.Pointer(nsString)}}).String()
}

// EventFilePathsFromPasteboard extracts file paths from a pasteboard object.
func EventFilePathsFromPasteboard(pasteboard uintptr) []string {
	// The pasteboard contains an array of items. For file drops, it's typically
	// an array of NSURL objects represented as strings.
	array := Objc_sendMsg[uintptr](pasteboard, Sel_getUid("propertyListForType:"), NSPasteboardTypeFileURL)
	if array == 0 {
		return nil
	}
	return eventFilePaths(Object{unsafe.Pointer(array)})
}

// eventFilePaths is a helper to convert an NSArray of NSStrings into a Go slice.
func eventFilePaths(array Object) []string {
	count := Objc_sendMsg[uintptr](uintptr(array.Ptr), Sel_count)
	if count == 0 {
		return nil
	}
	paths := make([]string, count)
	for i := uintptr(0); i < count; i++ {
		path := Objc_sendMsg[uintptr](uintptr(array.Ptr), Sel_objectAtIndex, i)
		paths[i] = (NSString{Object{unsafe.Pointer(path)}}).String()
	}
	return paths
}

func EventMagnification(event NSEvent) float64 {
	return Objc_sendMsg[float64](uintptr(event.Ptr), Sel_magnification)
}

func EventRotation(event NSEvent) float64 {
	return Objc_sendMsg[float64](uintptr(event.Ptr), Sel_rotation)
}

func EventPhase(event NSEvent) uintptr {
	return Objc_sendMsg[uintptr](uintptr(event.Ptr), Sel_phase)
}

func EventTranslationX(event NSEvent) float64 {
	return Objc_sendMsg[float64](uintptr(event.Ptr), Sel_deltaX)
}

func EventTranslationY(event NSEvent) float64 {
	return Objc_sendMsg[float64](uintptr(event.Ptr), Sel_deltaY)
}
