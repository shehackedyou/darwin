package darwin

import (
	"fmt"
	"os"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

type joystick struct {
	device  IOHIDDeviceRef
	name    string
	axes    []ioHIDElement
	buttons []ioHIDElement
	hats    []ioHIDElement
}

type ioHIDElement struct {
	element    uintptr
	logicalMin int64
	logicalMax int64
}

var (
	joystickMtx     sync.Mutex
	joystickManager IOHIDManagerRef
	joysticks       []joystick
)

const (
	kIOHIDPageGenericDesktop       = 0x01
	kIOHIDPageButton               = 0x09
	kIOHIDUsageJoystick            = 0x04
	kIOHIDUsageGamepad             = 0x05
	kIOHIDUsageMultiAxisController = 0x08
	kIOHIDUsageHatSwitch           = 0x39
	kIOHIDElementTypeAxis          = 2
	kIOHIDElementTypeButton        = 3
	kIOHIDElementTypeHatswitch     = 4
	kCFStringEncodingUTF8          = 0x08000100
)

func SetupJoysticks() {
	pool := NewAutoreleasePool()
	defer pool.Drain()

	joystickMtx.Lock()
	defer joystickMtx.Unlock()

	if joystickManager != 0 {
		return
	}

	mgr, _, _ := purego.SyscallN(_IOHIDManagerCreate, 0, uintptr(0))
	if mgr == 0 {
		fmt.Fprintf(os.Stderr, "darwin: IOHIDManagerCreate failed\n")
		return
	}
	joystickManager = IOHIDManagerRef(mgr)

	match := createDeviceMatchingArray()
	purego.SyscallN(_IOHIDManagerSetDeviceMatchingMultiple, uintptr(joystickManager), match)

	purego.SyscallN(_IOHIDManagerRegisterDeviceMatchingCallback, uintptr(joystickManager), purego.NewCallback(deviceMatchingCallback), 0)
	purego.SyscallN(_IOHIDManagerRegisterDeviceRemovalCallback, uintptr(joystickManager), purego.NewCallback(deviceRemovalCallback), 0)

	runLoop := Objc_sendMsg[uintptr](Class_NSRunLoop, Sel_mainRunLoop)
	purego.SyscallN(_IOHIDManagerScheduleWithRunLoop, uintptr(joystickManager), runLoop, NSDefaultRunLoopMode)

	purego.SyscallN(_IOHIDManagerOpen, uintptr(joystickManager), uintptr(0))
}

func createDeviceMatchingArray() uintptr {
	usagePageKey := NSString_WithUTF8String("UsagePage")
	usageKey := NSString_WithUTF8String("Usage")

	newNumber := func(v int) uintptr {
		return Objc_sendMsg[uintptr](Class_NSNumber, Sel_numberWithInt, v)
	}

	newDict := func(usage int) uintptr {
		keys := [2]uintptr{uintptr(usagePageKey.Ptr), uintptr(usageKey.Ptr)}
		values := [2]uintptr{newNumber(kIOHIDPageGenericDesktop), newNumber(usage)}
		return Objc_sendMsg[uintptr](Class_NSDictionary, Sel_dictionaryWithObjectsForKeysCount, &values[0], &keys[0], len(keys))
	}

	joystickDict := newDict(kIOHIDUsageJoystick)
	gamepadDict := newDict(kIOHIDUsageGamepad)
	multiAxisDict := newDict(kIOHIDUsageMultiAxisController)

	return Objc_sendMsg[uintptr](Class_NSArray, Sel_arrayWithObjects, joystickDict, gamepadDict, multiAxisDict, 0)
}

func deviceMatchingCallback(ctx, result, sender, device uintptr) {
	joystickMtx.Lock()
	defer joystickMtx.Unlock()
	addJoystick(IOHIDDeviceRef(device))
}

func deviceRemovalCallback(ctx, result, sender, device uintptr) {
	joystickMtx.Lock()
	defer joystickMtx.Unlock()
	removeJoystick(IOHIDDeviceRef(device))
}

func addJoystick(devRef IOHIDDeviceRef) {
	var j joystick
	j.device = devRef
	j.name = "Unknown Joystick"

	productName, _, _ := purego.SyscallN(_IOHIDDeviceGetProperty, uintptr(devRef), uintptr(NSString_WithUTF8String("Product").Ptr))
	if productName != 0 {
		j.name = GoString(Objc_sendMsg[uintptr](productName, Sel_UTF8String))
	}

	elements, _, _ := purego.SyscallN(_IOHIDDeviceCopyMatchingElements, uintptr(devRef), 0, 0)
	if elements == 0 {
		return
	}
	defer purego.SyscallN(_CFRelease, elements)

	count, _, _ := purego.SyscallN(_CFArrayGetCount, elements)
	for i := uintptr(0); i < count; i++ {
		elem, _, _ := purego.SyscallN(_CFArrayGetValueAtIndex, elements, i)
		if elem == 0 {
			continue
		}
		usagePage, _, _ := purego.SyscallN(_IOHIDElementGetUsagePage, elem)
		usage, _, _ := purego.SyscallN(_IOHIDElementGetUsage, elem)
		elemType, _, _ := purego.SyscallN(_IOHIDElementGetType, elem)

		hidElem := ioHIDElement{element: elem}
		min, _, _ := purego.SyscallN(_IOHIDElementGetLogicalMin, elem)
		hidElem.logicalMin = int64(min)
		max, _, _ := purego.SyscallN(_IOHIDElementGetLogicalMax, elem)
		hidElem.logicalMax = int64(max)

		if usagePage == kIOHIDPageGenericDesktop {
			if elemType == kIOHIDElementTypeAxis {
				j.axes = append(j.axes, hidElem)
			} else if usage == kIOHIDUsageHatSwitch {
				j.hats = append(j.hats, hidElem)
			}
		} else if usagePage == kIOHIDPageButton {
			j.buttons = append(j.buttons, hidElem)
		}
	}

	joysticks = append(joysticks, j)
	fmt.Fprintf(os.Stderr, "Joystick connected: %s\n", j.name)
}

func removeJoystick(devRef IOHIDDeviceRef) {
	for i, j := range joysticks {
		if j.device == devRef {
			name := j.name
			joysticks = append(joysticks[:i], joysticks[i+1:]...)
			fmt.Fprintf(os.Stderr, "Joystick disconnected: %s\n", name)
			return
		}
	}
}

func IsJoystickPresent(joy int) bool {
	joystickMtx.Lock()
	defer joystickMtx.Unlock()
	return joy >= 0 && joy < len(joysticks)
}

func GetJoystickName(joy int) string {
	joystickMtx.Lock()
	defer joystickMtx.Unlock()
	if joy < 0 || joy >= len(joysticks) {
		return ""
	}
	return joysticks[joy].name
}

func GetJoystickAxes(joy int) ([]float32, error) {
	joystickMtx.Lock()
	defer joystickMtx.Unlock()
	if joy < 0 || joy >= len(joysticks) {
		return nil, fmt.Errorf("joystick index %d is out of bounds", joy)
	}

	j := joysticks[joy]
	axes := make([]float32, len(j.axes))
	for i, elem := range j.axes {
		var hidValue uintptr
		purego.SyscallN(_IOHIDDeviceGetValue, uintptr(j.device), elem.element, uintptr(unsafe.Pointer(&hidValue)))
		if hidValue == 0 {
			continue
		}
		val, _, _ := purego.SyscallN(_IOHIDValueGetIntegerValue, hidValue)
		if elem.logicalMax != elem.logicalMin {
			normalized := 2.0*float32(int64(val)-elem.logicalMin)/float32(elem.logicalMax-elem.logicalMin) - 1.0
			axes[i] = normalized
		}
	}
	return axes, nil
}

func GetJoystickButtons(joy int) ([]byte, error) {
	joystickMtx.Lock()
	defer joystickMtx.Unlock()
	if joy < 0 || joy >= len(joysticks) {
		return nil, fmt.Errorf("joystick index %d is out of bounds", joy)
	}

	j := joysticks[joy]
	buttons := make([]byte, len(j.buttons))
	for i, elem := range j.buttons {
		var hidValue uintptr
		purego.SyscallN(_IOHIDDeviceGetValue, uintptr(j.device), elem.element, uintptr(unsafe.Pointer(&hidValue)))
		if hidValue == 0 {
			continue
		}
		val, _, _ := purego.SyscallN(_IOHIDValueGetIntegerValue, hidValue)
		if val != 0 {
			buttons[i] = 1
		}
	}
	return buttons, nil
}

func GetJoystickHats(joy int) ([]byte, error) {
	joystickMtx.Lock()
	defer joystickMtx.Unlock()
	if joy < 0 || joy >= len(joysticks) {
		return nil, fmt.Errorf("joystick index %d is out of bounds", joy)
	}

	j := joysticks[joy]
	hats := make([]byte, len(j.hats))
	for i, elem := range j.hats {
		var hidValue uintptr
		purego.SyscallN(_IOHIDDeviceGetValue, uintptr(j.device), elem.element, uintptr(unsafe.Pointer(&hidValue)))
		if hidValue == 0 {
			continue
		}
		val, _, _ := purego.SyscallN(_IOHIDValueGetIntegerValue, hidValue)
		hats[i] = byte(val)
	}
	return hats, nil
}
