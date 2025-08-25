package darwin

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	libAppKit, libFoundation, libCoreGraphics, libCoreOpenGL, libIOKit, libobjc, libCoreVideo uintptr
)

var (
	Class_NSApplication, Class_NSString, Class_NSWindow, Class_NSPasteboard, Class_NSOpenGLContext, Class_NSObject, Class_NSCursor, Class_NSImage, Class_NSBitmapImageRep, Class_cocoaWindowDelegate, Class_NSMenu, Class_NSMenuItem, Class_appDelegate, Class_NSOpenGLView, Class_NSAutoreleasePool, Class_NSThread, Class_NSOpenGLPixelFormat, Class_NSScreen, Class_NSRunLoop, Class_NSDictionary, Class_NSArray, Class_NSNumber, Class_NSTrackingArea, Class_NSColor, Class_NSImageView uintptr
)

var (
	// Application & Lifecycle Selectors
	Sel_registerName, Sel_alloc, Sel_init, Sel_release, Sel_retain, Sel_autorelease, Sel_drain, Sel_sharedApplication, Sel_setDelegate, Sel_delegate, Sel_setActivationPolicy, Sel_run, Sel_terminate, Sel_stop, Sel_activateIgnoringOtherApps, Sel_applicationDidFinishLaunching, Sel_applicationShouldTerminateAfterLastWindowClosed, Sel_applicationWillTerminate, Sel_new, Sel_isMainThread, Sel_performSelectorOnMainThread, Sel_call, Sel_mainRunLoop, Sel_class,
	// Menu Selectors
	Sel_setMainMenu, Sel_addItem, Sel_setSubmenu, Sel_addItemWithTitleActionKeyEquivalent, Sel_hide, Sel_hideOtherApplications, Sel_unhideAllApplications, Sel_miniaturize, Sel_performMiniaturize, Sel_performClose, Sel_selectAll, Sel_copy, Sel_paste, Sel_cut, Sel_undo, Sel_redo, Sel_setServicesMenu,
	// Window & View Selectors
	Sel_initWithContentRectStyleMaskBackingDefer, Sel_setTitle, Sel_setContentView, Sel_contentView, Sel_setOpenGLContext, Sel_makeCurrentContext, Sel_update, Sel_prepareOpenGL, Sel_clearCurrentContext, Sel_flushBuffer, Sel_CGLContextObj, Sel_close, Sel_backingScaleFactor, Sel_isKeyWindow, Sel_makeKeyAndOrderFront, Sel_toggleFullScreen, Sel_styleMask, Sel_setAutoresizingMask, Sel_initWithFrame, Sel_frame, Sel_setFrameTopLeftPoint, Sel_nextEventMatchingMaskUntilDateInModeDequeue, Sel_sendEvent, Sel_window, Sel_windowShouldClose, Sel_windowDidResize, Sel_object, Sel_setWantsBestResolutionOpenGLSurface, Sel_makeFirstResponder, Sel_acceptsFirstResponder, Sel_updateTrackingAreas, Sel_addTrackingArea, Sel_initWithRectOptionsOwnerUserInfo, Sel_initWithAttributes, Sel_screen, Sel_mainScreen, Sel_set, Sel_unhide, Sel_viewDidMoveToWindow, Sel_setBackgroundColor, Sel_colorWithSRGB, Sel_setTitlebarAppearsTransparent, Sel_setTitleVisibility, Sel_setWindowLevel, Sel_setCollectionBehavior,
	// Event Selectors
	Sel_keyCode, Sel_modifierFlags, Sel_characters, Sel_locationInWindow, Sel_scrollingDeltaX, Sel_scrollingDeltaY, Sel_buttonNumber, Sel_clickCount, Sel_phase, Sel_magnification, Sel_rotation, Sel_mouseMoved, Sel_mouseDragged, Sel_mouseDown, Sel_mouseUp, Sel_scrollWheel, Sel_keyDown, Sel_keyUp, Sel_flagsChanged, Sel_magnifyWithEvent, Sel_rotateWithEvent, Sel_swipeWithEvent, Sel_deltaX, Sel_deltaY,
	// Drag and Drop Selectors
	Sel_registerForDraggedTypes, Sel_draggingEntered, Sel_performDragOperation, Sel_concludeDragOperation, Sel_draggingPasteboard,
	// Pasteboard & String Selectors
	Sel_generalPasteboard, Sel_stringForType, Sel_setStringForType, Sel_UTF8String, Sel_initWithUTF8String, Sel_pboardTypeFileURL,
	// Collection & Number Selectors
	Sel_count, Sel_objectAtIndex, Sel_numberWithInt, Sel_dictionaryWithObjectsForKeysCount, Sel_arrayWithObjects, Sel_unsignedLongLongValue, Sel_deviceDescription Selector
)

var (
	objc_msgSend, class_getSuperclass_ptr uintptr
)

var (
	NSPasteboardTypeFileURL uintptr
)

var (
	objc_allocateClassPair_ptr, objc_registerClassPair_ptr, sel_getName_ptr, class_addMethod_ptr, object_setInstanceVariable_ptr, object_getInstanceVariable_ptr, class_addIvar_ptr uintptr
	_CGWarpMouseCursorPosition, _CGLFlushDrawable, _CFStringCreateWithCString, _CFNumberCreate, _IOHIDManagerCreate, _CFDictionaryCreateMutable, _IOHIDManagerSetDeviceMatchingMultiple, _IOHIDManagerRegisterDeviceMatchingCallback, _IOHIDManagerRegisterDeviceRemovalCallback, _IOHIDManagerScheduleWithRunLoop, _IOHIDManagerOpen uintptr
	_IOHIDDeviceGetProperty, _IOHIDDeviceCopyMatchingElements, _CFRelease, _CFArrayGetCount, _CFArrayGetValueAtIndex                                                                                                                                                                                                                  uintptr
	_IOHIDElementGetUsagePage, _IOHIDElementGetUsage, _IOHIDElementGetType, _IOHIDElementGetLogicalMin, _IOHIDElementGetLogicalMax, _IOHIDDeviceGetValue, _IOHIDValueGetIntegerValue                                                                                                                                                  uintptr
)

var initOnce sync.Once

func Initialize() {
	initOnce.Do(func() {
		if runtime.GOOS != "darwin" {
			return
		}
		runtime.LockOSThread()

		mustLoadLibraries()
		mustLoadFunctions()
		mustLoadClasses()
		mustRegisterSelectors()
		mustLoadConstants()

		setupAppDelegateClass()
		setupWindowDelegateClass()
		setupGoCallbackClass()
	})
}

func mustLoadLibraries() {
	var openErr error
	libAppKit, openErr = purego.Dlopen("/System/Library/Frameworks/AppKit.framework/AppKit", purego.RTLD_LAZY)
	if openErr != nil {
		panic("failed to load AppKit: " + openErr.Error())
	}
	libFoundation, openErr = purego.Dlopen("/System/Library/Frameworks/Foundation.framework/Foundation", purego.RTLD_LAZY)
	if openErr != nil {
		panic("failed to load Foundation: " + openErr.Error())
	}
	libCoreGraphics, openErr = purego.Dlopen("/System/Library/Frameworks/CoreGraphics.framework/CoreGraphics", purego.RTLD_LAZY)
	if openErr != nil {
		panic("failed to load CoreGraphics: " + openErr.Error())
	}
	libCoreOpenGL, openErr = purego.Dlopen("/System/Library/Frameworks/OpenGL.framework/OpenGL", purego.RTLD_LAZY)
	if openErr != nil {
		panic("failed to load OpenGL: " + openErr.Error())
	}
	libIOKit, openErr = purego.Dlopen("/System/Library/Frameworks/IOKit.framework/IOKit", purego.RTLD_LAZY)
	if openErr != nil {
		panic("failed to load IOKit: " + openErr.Error())
	}
	libobjc, openErr = purego.Dlopen("/usr/lib/libobjc.A.dylib", purego.RTLD_LAZY)
	if openErr != nil {
		panic("failed to load libobjc: " + openErr.Error())
	}
	libCoreVideo, openErr = purego.Dlopen("/System/Library/Frameworks/CoreVideo.framework/CoreVideo", purego.RTLD_LAZY)
	if openErr != nil {
		panic("failed to load CoreVideo: " + openErr.Error())
	}
}

func mustLoadFunctions() {
	load := func(lib uintptr, name string) uintptr {
		ptr, err := purego.Dlsym(lib, name)
		if err != nil {
			panic(fmt.Sprintf("critical failure: could not load function '%s': %v", name, err))
		}
		return ptr
	}

	objc_msgSend = load(libobjc, "objc_msgSend")
	Sel_registerName = Selector(load(libobjc, "sel_registerName"))
	sel_getName_ptr = load(libobjc, "sel_getName")
	objc_allocateClassPair_ptr = load(libobjc, "objc_allocateClassPair")
	objc_registerClassPair_ptr = load(libobjc, "objc_registerClassPair")
	class_addMethod_ptr = load(libobjc, "class_addMethod")
	class_getSuperclass_ptr = load(libobjc, "class_getSuperclass")

	object_setInstanceVariable_ptr = load(libobjc, "object_setInstanceVariable")
	object_getInstanceVariable_ptr = load(libobjc, "object_getInstanceVariable")
	class_addIvar_ptr = load(libobjc, "class_addIvar")

	_CGWarpMouseCursorPosition = load(libCoreGraphics, "CGWarpMouseCursorPosition")
	_CGLFlushDrawable = load(libCoreOpenGL, "CGLFlushDrawable")
	_CFStringCreateWithCString = load(libFoundation, "CFStringCreateWithCString")
	_CFNumberCreate = load(libFoundation, "CFNumberCreate")
	_CFDictionaryCreateMutable = load(libFoundation, "CFDictionaryCreateMutable")
	_CFRelease = load(libFoundation, "CFRelease")
	_CFArrayGetCount = load(libFoundation, "CFArrayGetCount")
	_CFArrayGetValueAtIndex = load(libFoundation, "CFArrayGetValueAtIndex")
	_IOHIDManagerCreate = load(libIOKit, "IOHIDManagerCreate")
	_IOHIDManagerSetDeviceMatchingMultiple = load(libIOKit, "IOHIDManagerSetDeviceMatchingMultiple")
	_IOHIDManagerRegisterDeviceMatchingCallback = load(libIOKit, "IOHIDManagerRegisterDeviceMatchingCallback")
	_IOHIDManagerRegisterDeviceRemovalCallback = load(libIOKit, "IOHIDManagerRegisterDeviceRemovalCallback")
	_IOHIDManagerScheduleWithRunLoop = load(libIOKit, "IOHIDManagerScheduleWithRunLoop")
	_IOHIDManagerOpen = load(libIOKit, "IOHIDManagerOpen")
	_IOHIDDeviceGetProperty = load(libIOKit, "IOHIDDeviceGetProperty")
	_IOHIDDeviceCopyMatchingElements = load(libIOKit, "IOHIDDeviceCopyMatchingElements")
	_IOHIDElementGetUsagePage = load(libIOKit, "IOHIDElementGetUsagePage")
	_IOHIDElementGetUsage = load(libIOKit, "IOHIDElementGetUsage")
	_IOHIDElementGetType = load(libIOKit, "IOHIDElementGetType")
	_IOHIDElementGetLogicalMin = load(libIOKit, "IOHIDElementGetLogicalMin")
	_IOHIDElementGetLogicalMax = load(libIOKit, "IOHIDElementGetLogicalMax")
	_IOHIDDeviceGetValue = load(libIOKit, "IOHIDDeviceGetValue")
	_IOHIDValueGetIntegerValue = load(libIOKit, "IOHIDValueGetIntegerValue")

	_CVDisplayLinkCreateWithCGDisplay = load(libCoreVideo, "CVDisplayLinkCreateWithCGDisplay")
	_CVDisplayLinkSetOutputCallback = load(libCoreVideo, "CVDisplayLinkSetOutputCallback")
	_CVDisplayLinkSetCurrentCGDisplay = load(libCoreVideo, "CVDisplayLinkSetCurrentCGDisplay")
	_CVDisplayLinkStart = load(libCoreVideo, "CVDisplayLinkStart")
	_CVDisplayLinkStop = load(libCoreVideo, "CVDisplayLinkStop")
	_CVDisplayLinkRelease = load(libCoreVideo, "CVDisplayLinkRelease")
}

func mustLoadClasses() {
	objc_getClass_ptr, err := purego.Dlsym(libobjc, "objc_getClass")
	if err != nil {
		panic("critical failure: could not load objc_getClass")
	}
	getClass := func(name string) uintptr {
		class, _, _ := purego.SyscallN(objc_getClass_ptr, uintptr(unsafe.Pointer(NewCString(name))))
		if class == 0 {
			panic(fmt.Sprintf("critical failure: failed to get Objective-C class: %s", name))
		}
		return class
	}
	Class_NSApplication = getClass("NSApplication")
	Class_NSString = getClass("NSString")
	Class_NSWindow = getClass("NSWindow")
	Class_NSPasteboard = getClass("NSPasteboard")
	Class_NSOpenGLContext = getClass("NSOpenGLContext")
	Class_NSObject = getClass("NSObject")
	Class_NSCursor = getClass("NSCursor")
	Class_NSImage = getClass("NSImage")
	Class_NSBitmapImageRep = getClass("NSBitmapImageRep")
	Class_NSMenu = getClass("NSMenu")
	Class_NSMenuItem = getClass("NSMenuItem")
	Class_NSOpenGLView = getClass("NSOpenGLView")
	Class_NSThread = getClass("NSThread")
	Class_NSAutoreleasePool = getClass("NSAutoreleasePool")
	Class_NSOpenGLPixelFormat = getClass("NSOpenGLPixelFormat")
	Class_NSScreen = getClass("NSScreen")
	Class_NSRunLoop = getClass("NSRunLoop")
	Class_NSDictionary = getClass("NSDictionary")
	Class_NSArray = getClass("NSArray")
	Class_NSNumber = getClass("NSNumber")
	Class_NSTrackingArea = getClass("NSTrackingArea")
	Class_NSColor = getClass("NSColor")
	Class_NSImageView = getClass("NSImageView")
}

func mustLoadConstants() {
	loadConstant := func(lib uintptr, name string) uintptr {
		ptr, err := purego.Dlsym(lib, name)
		if err != nil {
			panic(fmt.Sprintf("failed to load constant '%s': %v", name, err))
		}
		return *(*uintptr)(unsafe.Pointer(ptr))
	}
	NSDefaultRunLoopMode = loadConstant(libFoundation, "NSDefaultRunLoopMode")
	NSPasteboardTypeFileURL = loadConstant(libAppKit, "NSPasteboardTypeFileURL")
}

func mustRegisterSelectors() {
	Sel_alloc = Sel_getUid("alloc")
	Sel_init = Sel_getUid("init")
	Sel_release = Sel_getUid("release")
	Sel_retain = Sel_getUid("retain")
	Sel_autorelease = Sel_getUid("autorelease")
	Sel_drain = Sel_getUid("drain")
	Sel_sharedApplication = Sel_getUid("sharedApplication")
	Sel_setDelegate = Sel_getUid("setDelegate:")
	Sel_delegate = Sel_getUid("delegate")
	Sel_setActivationPolicy = Sel_getUid("setActivationPolicy:")
	Sel_setMainMenu = Sel_getUid("setMainMenu:")
	Sel_addItem = Sel_getUid("addItem:")
	Sel_setSubmenu = Sel_getUid("setSubmenu:")
	Sel_addItemWithTitleActionKeyEquivalent = Sel_getUid("addItemWithTitle:action:keyEquivalent:")
	Sel_run = Sel_getUid("run")
	Sel_terminate = Sel_getUid("terminate:")
	Sel_stop = Sel_getUid("stop:")
	Sel_hide = Sel_getUid("hide:")
	Sel_hideOtherApplications = Sel_getUid("hideOtherApplications:")
	Sel_unhideAllApplications = Sel_getUid("unhideAllApplications:")
	Sel_miniaturize = Sel_getUid("miniaturize:")
	Sel_performMiniaturize = Sel_getUid("performMiniaturize:")
	Sel_performClose = Sel_getUid("performClose:")
	Sel_selectAll = Sel_getUid("selectAll:")
	Sel_copy = Sel_getUid("copy:")
	Sel_paste = Sel_getUid("paste:")
	Sel_cut = Sel_getUid("cut:")
	Sel_undo = Sel_getUid("undo:")
	Sel_redo = Sel_getUid("redo:")
	Sel_setServicesMenu = Sel_getUid("setServicesMenu:")
	Sel_activateIgnoringOtherApps = Sel_getUid("activateIgnoringOtherApps:")
	Sel_applicationDidFinishLaunching = Sel_getUid("applicationDidFinishLaunching:")
	Sel_applicationShouldTerminateAfterLastWindowClosed = Sel_getUid("applicationShouldTerminateAfterLastWindowClosed:")
	Sel_applicationWillTerminate = Sel_getUid("applicationWillTerminate:")
	Sel_new = Sel_getUid("new")
	Sel_initWithContentRectStyleMaskBackingDefer = Sel_getUid("initWithContentRect:styleMask:backing:defer:")
	Sel_setTitle = Sel_getUid("setTitle:")
	Sel_setContentView = Sel_getUid("setContentView:")
	Sel_contentView = Sel_getUid("contentView")
	Sel_setOpenGLContext = Sel_getUid("setOpenGLContext:")
	Sel_makeCurrentContext = Sel_getUid("makeCurrentContext")
	Sel_update = Sel_getUid("update")
	Sel_prepareOpenGL = Sel_getUid("prepareOpenGL")
	Sel_clearCurrentContext = Sel_getUid("clearCurrentContext")
	Sel_flushBuffer = Sel_getUid("flushBuffer")
	Sel_CGLContextObj = Sel_getUid("CGLContextObj")
	Sel_close = Sel_getUid("close")
	Sel_backingScaleFactor = Sel_getUid("backingScaleFactor")
	Sel_isKeyWindow = Sel_getUid("isKeyWindow")
	Sel_makeKeyAndOrderFront = Sel_getUid("makeKeyAndOrderFront:")
	Sel_toggleFullScreen = Sel_getUid("toggleFullScreen:")
	Sel_styleMask = Sel_getUid("styleMask")
	Sel_setAutoresizingMask = Sel_getUid("setAutoresizingMask:")
	Sel_initWithFrame = Sel_getUid("initWithFrame:")
	Sel_frame = Sel_getUid("frame")
	Sel_setFrameTopLeftPoint = Sel_getUid("setFrameTopLeftPoint:")
	Sel_nextEventMatchingMaskUntilDateInModeDequeue = Sel_getUid("nextEventMatchingMask:untilDate:inMode:dequeue:")
	Sel_sendEvent = Sel_getUid("sendEvent:")
	Sel_window = Sel_getUid("window")
	Sel_windowShouldClose = Sel_getUid("windowShouldClose:")
	Sel_windowDidResize = Sel_getUid("windowDidResize:")
	Sel_object = Sel_getUid("object")
	Sel_isMainThread = Sel_getUid("isMainThread")
	Sel_performSelectorOnMainThread = Sel_getUid("performSelectorOnMainThread:withObject:waitUntilDone:")
	Sel_call = Sel_getUid("call")
	Sel_setWantsBestResolutionOpenGLSurface = Sel_getUid("setWantsBestResolutionOpenGLSurface:")
	Sel_generalPasteboard = Sel_getUid("generalPasteboard")
	Sel_stringForType = Sel_getUid("stringForType:")
	Sel_setStringForType = Sel_getUid("setString:forType:")
	Sel_UTF8String = Sel_getUid("UTF8String")
	Sel_initWithUTF8String = Sel_getUid("initWithUTF8String:")
	Sel_keyCode = Sel_getUid("keyCode")
	Sel_modifierFlags = Sel_getUid("modifierFlags")
	Sel_characters = Sel_getUid("characters")
	Sel_locationInWindow = Sel_getUid("locationInWindow")
	Sel_scrollingDeltaX = Sel_getUid("scrollingDeltaX")
	Sel_scrollingDeltaY = Sel_getUid("scrollingDeltaY")
	Sel_buttonNumber = Sel_getUid("buttonNumber")
	Sel_clickCount = Sel_getUid("clickCount")
	Sel_phase = Sel_getUid("phase")
	Sel_magnification = Sel_getUid("magnification")
	Sel_rotation = Sel_getUid("rotation")
	Sel_class = Sel_getUid("class")
	Sel_mouseMoved = Sel_getUid("mouseMoved:")
	Sel_mouseDragged = Sel_getUid("mouseDragged:")
	Sel_mouseDown = Sel_getUid("mouseDown:")
	Sel_mouseUp = Sel_getUid("mouseUp:")
	Sel_scrollWheel = Sel_getUid("scrollWheel:")
	Sel_keyDown = Sel_getUid("keyDown:")
	Sel_keyUp = Sel_getUid("keyUp:")
	Sel_flagsChanged = Sel_getUid("flagsChanged:")
	Sel_mainScreen = Sel_getUid("mainScreen")
	Sel_mainRunLoop = Sel_getUid("mainRunLoop")
	Sel_numberWithInt = Sel_getUid("numberWithInt:")
	Sel_dictionaryWithObjectsForKeysCount = Sel_getUid("dictionaryWithObjects:forKeys:count:")
	Sel_arrayWithObjects = Sel_getUid("arrayWithObjects:")
	Sel_unsignedLongLongValue = Sel_getUid("unsignedLongLongValue")
	Sel_makeFirstResponder = Sel_getUid("makeFirstResponder:")
	Sel_acceptsFirstResponder = Sel_getUid("acceptsFirstResponder")
	Sel_count = Sel_getUid("count")
	Sel_objectAtIndex = Sel_getUid("objectAtIndex:")
	Sel_deltaX = Sel_getUid("deltaX")
	Sel_deltaY = Sel_getUid("deltaY")
	Sel_addTrackingArea = Sel_getUid("addTrackingArea:")
	Sel_initWithRectOptionsOwnerUserInfo = Sel_getUid("initWithRect:options:owner:userInfo:")
	Sel_updateTrackingAreas = Sel_getUid("updateTrackingAreas")
	Sel_magnifyWithEvent = Sel_getUid("magnifyWithEvent:")
	Sel_rotateWithEvent = Sel_getUid("rotateWithEvent:")
	Sel_swipeWithEvent = Sel_getUid("swipeWithEvent:")
	Sel_screen = Sel_getUid("screen")
	Sel_deviceDescription = Sel_getUid("deviceDescription")
	Sel_registerForDraggedTypes = Sel_getUid("registerForDraggedTypes:")
	Sel_draggingEntered = Sel_getUid("draggingEntered:")
	Sel_performDragOperation = Sel_getUid("performDragOperation:")
	Sel_concludeDragOperation = Sel_getUid("concludeDragOperation:")
	Sel_draggingPasteboard = Sel_getUid("draggingPasteboard")
	Sel_initWithAttributes = Sel_getUid("initWithAttributes:")
	Sel_set = Sel_getUid("set")
	Sel_unhide = Sel_getUid("unhide")
	Sel_viewDidMoveToWindow = Sel_getUid("viewDidMoveToWindow")
	Sel_setBackgroundColor = Sel_getUid("setBackgroundColor:")
	Sel_colorWithSRGB = Sel_getUid("colorWithSRGBRed:green:blue:alpha:")
	Sel_setTitlebarAppearsTransparent = Sel_getUid("setTitlebarAppearsTransparent:")
	Sel_setTitleVisibility = Sel_getUid("setTitleVisibility:")
	Sel_setWindowLevel = Sel_getUid("setLevel:")
  Sel_setCollectionBehavior = Sel_getUid("setCollectionBehavior:")
}

func object_setInstanceVariable(obj uintptr, name string, value unsafe.Pointer) {
	_, _, _ = purego.SyscallN(object_setInstanceVariable_ptr, obj, uintptr(unsafe.Pointer(NewCString(name))), uintptr(value))
}

func object_getInstanceVariable(obj uintptr, name string, value unsafe.Pointer) {
	_, _, _ = purego.SyscallN(object_getInstanceVariable_ptr, obj, uintptr(unsafe.Pointer(NewCString(name))), uintptr(value))
}

func class_addIvar(class uintptr, name string, size, alignment uintptr, types string) bool {
	ret, _, _ := purego.SyscallN(class_addIvar_ptr, class, uintptr(unsafe.Pointer(NewCString(name))), size, alignment, uintptr(unsafe.Pointer(NewCString(types))))
	return ret != 0
}
