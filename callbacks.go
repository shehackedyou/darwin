package darwin

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/ebitengine/purego"
)

// WindowDelegate is a Go interface that our native callbacks will call.
type WindowDelegate interface {
	WindowShouldClose()
	WindowDidResize(resizedWindow NSWindow)
	KeyDown(event NSEvent)
	KeyUp(event NSEvent)
	MouseDown(event NSEvent)
	MouseUp(event NSEvent)
	MouseMoved(event NSEvent)
	MouseDragged(event NSEvent)
	ScrollWheel(event NSEvent)
	FlagsChanged(event NSEvent)
	MagnifyGesture(event NSEvent)
	RotateGesture(event NSEvent)
	SwipeGesture(event NSEvent)
	FilesDropped(files []string)
}

var appDelegateCallback func()
var appTerminationCallback func()

func SetAppDelegateCallback(f func()) {
	appDelegateCallback = f
}

func SetAppTerminationCallback(f func()) {
	appTerminationCallback = f
}

func applicationDidFinishLaunching(id, sel, notification uintptr) {
	if appDelegateCallback != nil {
		appDelegateCallback()
	}
}

var shouldTerminateAfterLastWindowClosed = true

func applicationShouldTerminateAfterLastWindowClosed(id, sel, sender uintptr) bool {
	return shouldTerminateAfterLastWindowClosed
}

func applicationWillTerminate(id, sel, notification uintptr) {
	if appTerminationCallback != nil {
		appTerminationCallback()
	}
}

func getGoWindowDelegate(viewInstance uintptr) WindowDelegate {
	if viewInstance == 0 {
		return nil
	}
	var goPtr uintptr
	object_getInstanceVariable(viewInstance, "goWindowPtr", unsafe.Pointer(&goPtr))
	goObj := GetGoPointer(goPtr)
	if goObj == nil {
		return nil
	}
	delegate, ok := goObj.(WindowDelegate)
	if !ok {
		fmt.Printf("CRITICAL: Stored Go pointer for view %p is not a WindowDelegate\n", unsafe.Pointer(viewInstance))
		return nil
	}
	return delegate
}

func acceptsFirstResponder(id, sel uintptr) bool {
	// log.Println("[NATIVE] acceptsFirstResponder called")
	return true
}

func viewDidMoveToWindow(id, sel uintptr) {
	log.Println("[NATIVE] viewDidMoveToWindow called")
	window := Objc_sendMsg[uintptr](id, Sel_window)
	if window != 0 {
		log.Println("[NATIVE] View has a window, attempting to make first responder.")
		// This makes our view the target for keyboard and other events.
		Objc_sendMsg[bool](window, Sel_makeFirstResponder, id)
	}
}

func updateTrackingAreas(id, sel uintptr) {
	// log.Println("[NATIVE] updateTrackingAreas called")
	class := Objc_sendMsg[uintptr](id, Sel_class)
	superClass, _, _ := purego.SyscallN(class_getSuperclass_ptr, class)
	super := objc_super{
		Receiver:   id,
		SuperClass: superClass,
	}
	var super_updateTrackingAreas func(*objc_super, Selector)
	purego.RegisterLibFunc(&super_updateTrackingAreas, libobjc, "objc_msgSendSuper")
	super_updateTrackingAreas(&super, Sel_updateTrackingAreas)

	trackingAreas := Objc_sendMsg[uintptr](id, Sel_getUid("trackingAreas"))
	count := Objc_sendMsg[uintptr](trackingAreas, Sel_count)
	for i := uintptr(0); i < count; i++ {
		area := Objc_sendMsg[uintptr](trackingAreas, Sel_objectAtIndex, i)
		Objc_sendMsg[uintptr](id, Sel_getUid("removeTrackingArea:"), area)
	}

	frame := Objc_sendMsg[NSRect](id, Sel_frame)
	options := NSTrackingMouseMoved | NSTrackingActiveInKeyWindow | NSTrackingMouseEnteredAndExited

	trackingAreaAlloc := Objc_sendMsg[uintptr](Class_NSTrackingArea, Sel_alloc)
	var initWithRectOptions func(uintptr, Selector, NSRect, uintptr, uintptr, uintptr) uintptr
	purego.RegisterLibFunc(&initWithRectOptions, libobjc, "objc_msgSend")
	trackingArea := initWithRectOptions(trackingAreaAlloc, Sel_initWithRectOptionsOwnerUserInfo, frame, uintptr(options), id, 0)

	if trackingArea != 0 {
		Objc_sendMsg[uintptr](id, Sel_addTrackingArea, trackingArea)
		Objc_sendMsg[uintptr](trackingArea, Sel_release)
	}
}

func keyDown(id, sel, event uintptr) {
	log.Println("[NATIVE] keyDown called")
	if delegate := getGoWindowDelegate(id); delegate != nil {
		delegate.KeyDown(NSEvent{Object{unsafe.Pointer(event)}})
	}
}

func keyUp(id, sel, event uintptr) {
	log.Println("[NATIVE] keyUp called")
	if delegate := getGoWindowDelegate(id); delegate != nil {
		delegate.KeyUp(NSEvent{Object{unsafe.Pointer(event)}})
	}
}

func mouseDown(id, sel, event uintptr) {
	log.Println("[NATIVE] mouseDown called")
	if delegate := getGoWindowDelegate(id); delegate != nil {
		delegate.MouseDown(NSEvent{Object{unsafe.Pointer(event)}})
	}
}

func mouseUp(id, sel, event uintptr) {
	log.Println("[NATIVE] mouseUp called")
	if delegate := getGoWindowDelegate(id); delegate != nil {
		delegate.MouseUp(NSEvent{Object{unsafe.Pointer(event)}})
	}
}

func mouseMoved(id, sel, event uintptr) {
	// This can be very noisy, so logging is commented out.
	// log.Println("[NATIVE] mouseMoved called")
	if delegate := getGoWindowDelegate(id); delegate != nil {
		delegate.MouseMoved(NSEvent{Object{unsafe.Pointer(event)}})
	}
}

func mouseDragged(id, sel, event uintptr) {
	// log.Println("[NATIVE] mouseDragged called")
	if delegate := getGoWindowDelegate(id); delegate != nil {
		delegate.MouseDragged(NSEvent{Object{unsafe.Pointer(event)}})
	}
}

func scrollWheel(id, sel, event uintptr) {
	log.Println("[NATIVE] scrollWheel called")
	if delegate := getGoWindowDelegate(id); delegate != nil {
		delegate.ScrollWheel(NSEvent{Object{unsafe.Pointer(event)}})
	}
}

func flagsChanged(id, sel, event uintptr) {
	log.Println("[NATIVE] flagsChanged called")
	if delegate := getGoWindowDelegate(id); delegate != nil {
		delegate.FlagsChanged(NSEvent{Object{unsafe.Pointer(event)}})
	}
}

func draggingEntered(id, sel, sender uintptr) uintptr {
	log.Println("[NATIVE] draggingEntered called")
	const NSDragOperationCopy = 1
	return NSDragOperationCopy
}

func performDragOperation(id, sel, sender uintptr) bool {
	log.Println("[NATIVE] performDragOperation called")
	pb := Objc_sendMsg[uintptr](sender, Sel_draggingPasteboard)
	if pb == 0 {
		return false
	}

	paths := EventFilePathsFromPasteboard(pb)
	if len(paths) > 0 {
		if delegate := getGoWindowDelegate(id); delegate != nil {
			delegate.FilesDropped(paths)
			return true
		}
	}
	return false
}

func windowShouldClose(id, sel, window uintptr) bool {
	log.Println("[NATIVE] windowShouldClose called")
	if delegate := getGoWindowDelegate(id); delegate != nil {
		delegate.WindowShouldClose()
	}
	var goPtr uintptr
	object_getInstanceVariable(id, "goWindowPtr", unsafe.Pointer(&goPtr))
	if goPtr != 0 {
		FreeGoPointer(goPtr)
	}
	return true
}

func windowDidResize(id, sel, notification uintptr) {
	log.Println("[NATIVE] windowDidResize called")
	if delegate := getGoWindowDelegate(id); delegate != nil {
		windowObject := Object{unsafe.Pointer(Objc_sendMsg[uintptr](notification, Sel_object))}
		delegate.WindowDidResize(NSWindow{windowObject})
	}
}

func setupCustomOpenGLViewClass() {
	className := "GoCustomOpenGLView"
	class := objc_allocateClassPair(Class_NSOpenGLView, className, 0)
	if class == 0 {
		panic("failed to allocate GoCustomOpenGLView class")
	}

	ok := class_addIvar(class, "goWindowPtr", unsafe.Sizeof(uintptr(0)), 3, "@")
	if !ok {
		panic("failed to add ivar to GoCustomOpenGLView")
	}
	Class_cocoaWindowDelegate = class

	addMethod := func(selector Selector, callback any, types string) {
		imp := purego.NewCallback(callback)
		if !class_addMethod(class, selector, imp, types) {
			panic(fmt.Sprintf("failed to add method %s to GoCustomOpenGLView", sel_getName(selector)))
		}
	}

	addMethod(Sel_acceptsFirstResponder, acceptsFirstResponder, "B@:")
	addMethod(Sel_viewDidMoveToWindow, viewDidMoveToWindow, "v@:")
	addMethod(Sel_updateTrackingAreas, updateTrackingAreas, "v@:")

	addMethod(Sel_keyDown, keyDown, "v@:@")
	addMethod(Sel_keyUp, keyUp, "v@:@")
	addMethod(Sel_mouseDown, mouseDown, "v@:@")
	addMethod(Sel_mouseUp, mouseUp, "v@:@")
	addMethod(Sel_mouseMoved, mouseMoved, "v@:@")
	addMethod(Sel_mouseDragged, mouseDragged, "v@:@")
	addMethod(Sel_scrollWheel, scrollWheel, "v@:@")
	addMethod(Sel_flagsChanged, flagsChanged, "v@:@")

	addMethod(Sel_draggingEntered, draggingEntered, "I@:@")
	addMethod(Sel_performDragOperation, performDragOperation, "B@:@")

	addMethod(Sel_windowShouldClose, windowShouldClose, "B@:@")
	addMethod(Sel_windowDidResize, windowDidResize, "v@:@")

	objc_registerClassPair(class)
}

func setupWindowDelegateClass() {
	setupCustomOpenGLViewClass()
}

func setupAppDelegateClass() {
	className := "GoAppDelegate"
	class := objc_allocateClassPair(Class_NSObject, className, 0)
	if class == 0 {
		panic("failed to allocate GoAppDelegate class")
	}
	Class_appDelegate = class
	addMethod := func(selector Selector, callback any, types string) {
		imp := purego.NewCallback(callback)
		if !class_addMethod(class, selector, imp, types) {
			panic(fmt.Sprintf("failed to add method %s to GoAppDelegate", sel_getName(selector)))
		}
	}
	addMethod(Sel_applicationDidFinishLaunching, applicationDidFinishLaunching, "v@:@")
	addMethod(Sel_applicationShouldTerminateAfterLastWindowClosed, applicationShouldTerminateAfterLastWindowClosed, "B@:@")
	addMethod(Sel_applicationWillTerminate, applicationWillTerminate, "v@:@")
	objc_registerClassPair(class)
}
