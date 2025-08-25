package darwin

import (
	"fmt"
	"unsafe"
)

func NSString_WithUTF8String(s string) NSString {
	// This function now returns an autoreleased string. The caller does not
	// need to manually release the returned object, as it will be handled
	// by the current NSAutoreleasePool. This is the idiomatic way to return
	// newly-created objects in Objective-C.
	class := Objc_sendMsg[uintptr](Class_NSString, Sel_alloc)
	cString := NewCString(s)
	nsStringPtr := Objc_sendMsg[uintptr](class, Sel_initWithUTF8String, uintptr(unsafe.Pointer(cString)))
	
	nsStringObj := Object{unsafe.Pointer(nsStringPtr)}
	nsStringObj.Autorelease()

	return NSString{nsStringObj}
}

func (s NSString) String() string {
	if s.Ptr == nil {
		return ""
	}
	cString := Objc_sendMsg[uintptr](uintptr(s.Ptr), Sel_UTF8String)
	return GoString(cString)
}

func GetClipboardString() (string, error) {
	pool := NewAutoreleasePool()
	defer pool.Drain()
	
	pb := Objc_sendMsg[uintptr](Class_NSPasteboard, Sel_generalPasteboard)
	if pb == 0 {
		return "", fmt.Errorf("failed to get general pasteboard")
	}

	typeString := NSString_WithUTF8String("public.utf8-plain-text")
	ret := Objc_sendMsg[uintptr](pb, Sel_stringForType, uintptr(typeString.Ptr))
	if ret == 0 {
		return "", nil
	}
	
	nsString := NSString{Object{unsafe.Pointer(ret)}}
	return nsString.String(), nil
}

func SetClipboardString(value string) error {
	pool := NewAutoreleasePool()
	defer pool.Drain()

	pb := Objc_sendMsg[uintptr](Class_NSPasteboard, Sel_generalPasteboard)
	if pb == 0 {
		return fmt.Errorf("failed to get general pasteboard")
	}

	// The setString:forType: method clears previous contents and then sets the new value.
	// It's a single, atomic operation for this common case.
	nsValue := NSString_WithUTF8String(value)
	typeString := NSString_WithUTF8String("public.utf8-plain-text")

	ok := Objc_sendMsg[bool](pb, Sel_setStringForType, uintptr(nsValue.Ptr), uintptr(typeString.Ptr))
	if !ok {
		return fmt.Errorf("setString:forType: returned false, indicating failure")
	}

	return nil
}
