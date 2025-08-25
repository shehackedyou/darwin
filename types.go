package darwin

import (
	"sync"
	"unsafe"
)

type Selector uintptr
type IOHIDManagerRef uintptr
type IOHIDDeviceRef uintptr
type CFStringRef uintptr
type CFArrayRef uintptr
type CFRunLoopRef uintptr

// Object is a wrapper for a raw Objective-C object pointer.
type Object struct {
	Ptr unsafe.Pointer
}

type (
	NSString            struct{ Object }
	NSWindow            struct{ Object }
	NSPasteboard        struct{ Object }
	NSEvent             struct{ Object }
	NSCursor            struct{ Object }
	NSOpenGLContext     struct{ Object }
	NSOpenGLView        struct{ Object }
	NSImage             struct{ Object }
	NSMenu              struct{ Object }
	NSMenuItem          struct{ Object }
	NSAutoreleasePool   struct{ Object }
	NSOpenGLPixelFormat struct{ Object }
	NSScreen            struct{ Object }
	NSRunLoop           struct{ Object }
	NSDictionary        struct{ Object }
	NSArray             struct{ Object }
	NSNumber            struct{ Object }
	NSTrackingArea      struct{ Object }
	NSImageView         struct{ Object }
	NSColor             struct{ Object }
)

type (
	NSPoint struct{ X, Y float64 }
	NSSize  struct{ Width, Height float64 }
	NSRect  struct {
		Origin NSPoint
		Size   NSSize
	}
)

type objc_super struct {
	Receiver   uintptr
	SuperClass uintptr
}

type (
	NSUInteger    uintptr
	CGLContextObj uintptr
)

const (
	NSOpenGLPFADoubleBuffer       = 5
	NSOpenGLPFADepthSize          = 12
	NSOpenGLPFAAccelerated        = 73
	NSOpenGLPFAOpenGLProfile      = 99
	NSOpenGLProfileVersion4_1Core = 0x4100
)

const (
	NSBackingStoreBuffered      = 2
	NSWindowStyleMaskBorderless = 0
	NSWindowStyleMaskTitled     = 1 << 0
	NSWindowStyleMaskClosable   = 1 << 1
	NSWindowStyleMaskResizable  = 1 << 3
	NSWindowStyleMaskFullScreen = 1 << 14
	NSEventMaskAny              = 0xFFFFFFFF
	NSViewWidthSizable          = 2
	NSViewHeightSizable         = 16
  NSWindowCollectionBehaviorFullScreenPrimary = 1 << 7
)

const (
	NSWindowTitleVisible = 0
	NSWindowTitleHidden  = 1
)

const (
	NSEventModifierFlagShift   = 1 << 17
	NSEventModifierFlagControl = 1 << 18
	NSEventModifierFlagOption  = 1 << 19 // Alt key
	NSEventModifierFlagCommand = 1 << 20
)

const (
	NSTrackingMouseEnteredAndExited = 0x01
	NSTrackingMouseMoved            = 0x02
	NSTrackingActiveInKeyWindow     = 0x20
)

type DarwinCursorMode int

const (
	DarwinCursorNormal DarwinCursorMode = iota
	DarwinCursorHidden
	DarwinCursorDisabled
)

var (
	windowMapMtx sync.RWMutex
	WindowMap    = make(map[NSWindow]any)
)

var NSDefaultRunLoopMode uintptr

type MenuItem struct {
	Title         string
	Action        Selector
	Key           string
	IsSeparator   bool
	Submenu       *ApplicationMenu
	ModifierFlags []uintptr
}

type ApplicationMenu struct {
	AppItems    []MenuItem
	FileItems   []MenuItem
	EditItems   []MenuItem
	WindowItems []MenuItem
}
