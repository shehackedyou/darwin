package darwin

import (
	"fmt"
	"image"
	"image/draw"
	"unsafe"

	"github.com/ebitengine/purego"
)

func NewSplashWindow(img image.Image) (NSWindow, error) {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	rect := NSRect{Size: NSSize{Width: float64(width), Height: float64(height)}}

	styleMask := NSWindowStyleMaskBorderless

	winAlloc := Objc_sendMsg[uintptr](Class_NSWindow, Sel_alloc)
	if winAlloc == 0 {
		return NSWindow{}, fmt.Errorf("darwin: failed to allocate splash NSWindow")
	}

	var initWithContentRect func(uintptr, Selector, NSRect, uintptr, uintptr, bool) uintptr
	purego.RegisterLibFunc(&initWithContentRect, libobjc, "objc_msgSend")
	win := initWithContentRect(winAlloc, Sel_initWithContentRectStyleMaskBackingDefer, rect, uintptr(styleMask), uintptr(NSBackingStoreBuffered), true)
	if win == 0 {
		return NSWindow{}, fmt.Errorf("darwin: failed to initialize splash NSWindow")
	}

	nsWin := NSWindow{Object{unsafe.Pointer(win)}}

	screenFrame := mainNSScreen().Frame()
	originX := (screenFrame.Size.Width - float64(width)) / 2
	originY := (screenFrame.Size.Height - float64(height)) / 2
	var setFrameOrigin func(uintptr, Selector, NSPoint)
	purego.RegisterLibFunc(&setFrameOrigin, libobjc, "objc_msgSend")
	setFrameOrigin(uintptr(nsWin.Ptr), Sel_getUid("setFrameOrigin:"), NSPoint{X: originX, Y: originY})

	nsImage, err := nsImageFromGoImage(img)
	if err != nil {
		nsWin.Release()
		return NSWindow{}, fmt.Errorf("failed to create NSImage for splash: %w", err)
	}

	imageViewAlloc := Objc_sendMsg[uintptr](Class_NSImageView, Sel_alloc)
	imageView := Objc_sendMsg[uintptr](imageViewAlloc, Sel_initWithFrame, rect)
	Objc_sendMsg[uintptr](imageView, Sel_getUid("setImage:"), uintptr(nsImage.Ptr))
	nsImage.Release()

	SetContentView(nsWin, Object{unsafe.Pointer(imageView)})
	return nsWin, nil
}

func NewNSWindow(title string, width, height int) (NSWindow, error) {
	rect := NSRect{Size: NSSize{Width: float64(width), Height: float64(height)}}
	styleMask := NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskResizable

	winAlloc := Objc_sendMsg[uintptr](Class_NSWindow, Sel_alloc)
	if winAlloc == 0 {
		return NSWindow{}, fmt.Errorf("darwin: failed to allocate NSWindow")
	}

	var initWithContentRect func(uintptr, Selector, NSRect, uintptr, uintptr, bool) uintptr
	purego.RegisterLibFunc(&initWithContentRect, libobjc, "objc_msgSend")
	win := initWithContentRect(winAlloc, Sel_initWithContentRectStyleMaskBackingDefer, rect, uintptr(styleMask), uintptr(NSBackingStoreBuffered), true)
	if win == 0 {
		return NSWindow{}, fmt.Errorf("darwin: failed to initialize NSWindow with content rect")
	}

	nsWin := NSWindow{Object{unsafe.Pointer(win)}}
	nsWin.SetTitle(title)

	Objc_sendMsg[uintptr](uintptr(nsWin.Ptr), Sel_setCollectionBehavior, NSWindowCollectionBehaviorFullScreenPrimary)

	screenFrame := mainNSScreen().Frame()
	originX := (screenFrame.Size.Width - float64(width)) / 2
	originY := (screenFrame.Size.Height - float64(height)) / 2
	var setFrameOrigin func(uintptr, Selector, NSPoint)
	purego.RegisterLibFunc(&setFrameOrigin, libobjc, "objc_msgSend")
	setFrameOrigin(uintptr(nsWin.Ptr), Sel_getUid("setFrameOrigin:"), NSPoint{X: originX, Y: originY})

	nsWin.SetBackgroundColor(0.2, 0.3, 0.3, 1.0)
	return nsWin, nil
}

func NewNSWindowOpenGL(title string, width, height int, major, minor int) (NSWindow, NSOpenGLView, NSOpenGLContext, error) {
	win, err := NewNSWindow(title, width, height)
	if err != nil {
		return NSWindow{}, NSOpenGLView{}, NSOpenGLContext{}, err
	}

	attrs := []uint32{
		NSOpenGLPFAOpenGLProfile, NSOpenGLProfileVersion4_1Core,
		NSOpenGLPFADoubleBuffer,
		NSOpenGLPFAAccelerated,
		NSOpenGLPFADepthSize, 24,
		0,
	}

	pixelFormatAlloc := Objc_sendMsg[uintptr](Class_NSOpenGLPixelFormat, Sel_alloc)
	pixelFormatPtr := Objc_sendMsg[uintptr](pixelFormatAlloc, Sel_initWithAttributes, unsafe.Pointer(&attrs[0]))
	if pixelFormatPtr == 0 {
		return NSWindow{}, NSOpenGLView{}, NSOpenGLContext{}, fmt.Errorf("darwin: failed to initialize NSOpenGLPixelFormat")
	}
	pixelFormat := NSOpenGLPixelFormat{Object{unsafe.Pointer(pixelFormatPtr)}}
	defer pixelFormat.Release()

	frame := NSRect{Origin: NSPoint{X: 0, Y: 0}, Size: NSSize{Width: float64(width), Height: float64(height)}}
	view, err := NewCustomOpenGLView(frame, pixelFormat)
	if err != nil {
		return NSWindow{}, NSOpenGLView{}, NSOpenGLContext{}, fmt.Errorf("darwin: failed to create custom OpenGL view: %w", err)
	}

	types := Objc_sendMsg[uintptr](Class_NSArray, Sel_arrayWithObjects, NSPasteboardTypeFileURL, 0)
	Objc_sendMsg[uintptr](uintptr(win.Ptr), Sel_registerForDraggedTypes, types)
	Objc_sendMsg[uintptr](uintptr(view.Ptr), Sel_setWantsBestResolutionOpenGLSurface, true)

	initWithPixelFormatSel := Sel_getUid("initWithFormat:shareContext:")
	ctxAlloc := Objc_sendMsg[uintptr](Class_NSOpenGLContext, Sel_alloc)
	ctxPtr := Objc_sendMsg[uintptr](ctxAlloc, initWithPixelFormatSel, pixelFormatPtr, nil)
	if ctxPtr == 0 {
		view.Release()
		return NSWindow{}, NSOpenGLView{}, NSOpenGLContext{}, fmt.Errorf("darwin: failed to init NSOpenGLContext")
	}
	ctx := NSOpenGLContext{Object{unsafe.Pointer(ctxPtr)}}
	SetContentView(win, view.Object)
	SetOpenGLContext(view, ctx)
	Objc_sendMsg[uintptr](uintptr(view.Ptr), Sel_prepareOpenGL)

	return win, view, ctx, nil
}

func NewCustomOpenGLView(frame NSRect, pixelFormat NSOpenGLPixelFormat) (NSOpenGLView, error) {
	viewAlloc := Objc_sendMsg[uintptr](Class_cocoaWindowDelegate, Sel_alloc)
	if viewAlloc == 0 {
		return NSOpenGLView{}, fmt.Errorf("darwin: failed to allocate CustomOpenGLView")
	}

	initWithFramePixelFormatSel := Sel_getUid("initWithFrame:pixelFormat:")
	var initWithFramePixelFormat func(uintptr, Selector, NSRect, uintptr) uintptr
	purego.RegisterLibFunc(&initWithFramePixelFormat, libobjc, "objc_msgSend")
	viewPtr := initWithFramePixelFormat(viewAlloc, initWithFramePixelFormatSel, frame, uintptr(pixelFormat.Ptr))
	if viewPtr == 0 {
		return NSOpenGLView{}, fmt.Errorf("darwin: failed to initialize CustomOpenGLView")
	}
	view := NSOpenGLView{Object{unsafe.Pointer(viewPtr)}}
	Objc_sendMsg[uintptr](viewPtr, Sel_setAutoresizingMask, NSViewWidthSizable|NSViewHeightSizable)

	return view, nil
}

func (w NSWindow) SetTitle(title string) {
	nsTitle := NSString_WithUTF8String(title)
	Objc_sendMsg[uintptr](uintptr(w.Ptr), Sel_setTitle, uintptr(nsTitle.Ptr))
}

func (w NSWindow) SetBackgroundColor(r, g, b, a float64) {
	color := Objc_sendMsg[uintptr](Class_NSColor, Sel_colorWithSRGB, r, g, b, a)
	Objc_sendMsg[uintptr](uintptr(w.Ptr), Sel_setBackgroundColor, color)
}

func (w NSWindow) SetTitlebarAppearsTransparent(transparent bool) {
	Objc_sendMsg[uintptr](uintptr(w.Ptr), Sel_setTitlebarAppearsTransparent, transparent)
}

func (w NSWindow) SetTitleVisibility(visible bool) {
	visibility := NSWindowTitleVisible
	if !visible {
		visibility = NSWindowTitleHidden
	}
	Objc_sendMsg[uintptr](uintptr(w.Ptr), Sel_setTitleVisibility, visibility)
}

func (w NSWindow) SetWindowLevel(level int) {
	Objc_sendMsg[uintptr](uintptr(w.Ptr), Sel_setWindowLevel, level)
}

func (w NSWindow) ContentSize() (width, height int, scale float64) {
	contentView := Objc_sendMsg[uintptr](uintptr(w.Ptr), Sel_contentView)
	if contentView == 0 {
		return 0, 0, 1.0
	}
	var frameFunc func(uintptr, Selector) NSRect
	purego.RegisterLibFunc(&frameFunc, libobjc, "objc_msgSend")
	frame := frameFunc(contentView, Sel_frame)
	scale = Objc_sendMsg[float64](uintptr(w.Ptr), Sel_backingScaleFactor)
	return int(frame.Size.Width), int(frame.Size.Height), scale
}

func SetContentView(w NSWindow, v Object) {
	Objc_sendMsg[uintptr](uintptr(w.Ptr), Sel_setContentView, uintptr(v.Ptr))
}

func SetOpenGLContext(v NSOpenGLView, ctx NSOpenGLContext) {
	Objc_sendMsg[uintptr](uintptr(v.Ptr), Sel_setOpenGLContext, uintptr(ctx.Ptr))
}

func SetDelegateAndLinkGo(w NSWindow, delegateAsView NSOpenGLView, goWindow any) {
	ptrID := StoreGoPointer(goWindow)
	object_setInstanceVariable(uintptr(delegateAsView.Ptr), "goWindowPtr", unsafe.Pointer(&ptrID))
	Objc_sendMsg[uintptr](uintptr(w.Ptr), Sel_setDelegate, uintptr(delegateAsView.Ptr))
}

func MakeKeyAndOrderFront(w NSWindow) {
	Objc_sendMsg[uintptr](uintptr(w.Ptr), Sel_makeKeyAndOrderFront, 0)
}

func CloseWindow(w NSWindow) {
	Objc_sendMsg[uintptr](uintptr(w.Ptr), Sel_close)
}

func MakeCurrentOpenGLContext(ctx NSOpenGLContext) {
	if ctx.Ptr != nil {
		Objc_sendMsg[uintptr](uintptr(ctx.Ptr), Sel_makeCurrentContext)
	}
}

func FlushBuffer(ctx NSOpenGLContext) {
	if ctx.Ptr != nil {
		Objc_sendMsg[uintptr](uintptr(ctx.Ptr), Sel_flushBuffer)
	}
}

func IsKeyWindow(w NSWindow) bool {
	return Objc_sendMsg[bool](uintptr(w.Ptr), Sel_isKeyWindow)
}

func SetWindowFrameTopLeftPoint(w NSWindow, x, y int) {
	screen := mainNSScreen()
	screenFrame := screen.Frame()
	windowFrame := w.Frame()
	newY := screenFrame.Size.Height - float64(y) - windowFrame.Size.Height

	point := NSPoint{X: float64(x), Y: newY}
	var setFrameTopLeftPoint func(uintptr, Selector, NSPoint)
	purego.RegisterLibFunc(&setFrameTopLeftPoint, libobjc, "objc_msgSend")
	setFrameTopLeftPoint(uintptr(w.Ptr), Sel_setFrameTopLeftPoint, point)
}

func WindowFrameTopLeftPoint(w NSWindow) (int, int) {
	screen := mainNSScreen()
	screenFrame := screen.Frame()
	windowFrame := w.Frame()
	topLeftY := screenFrame.Size.Height - windowFrame.Origin.Y - windowFrame.Size.Height
	return int(windowFrame.Origin.X), int(topLeftY)
}

func IsWindowFullscreen(w NSWindow) bool {
	style := Objc_sendMsg[uintptr](uintptr(w.Ptr), Sel_styleMask)
	return (style & NSWindowStyleMaskFullScreen) != 0
}

func ToggleWindowFullScreen(w NSWindow) {
	Objc_sendMsg[uintptr](uintptr(w.Ptr), Sel_toggleFullScreen, 0)
}

func SetCursor(cursor NSCursor) {
	Objc_sendMsg[uintptr](uintptr(cursor.Ptr), Sel_set)
}

func SetCursorMode(mode DarwinCursorMode) {
	if mode == DarwinCursorHidden || mode == DarwinCursorDisabled {
		Objc_sendMsg[uintptr](Class_NSCursor, Sel_hide)
	} else {
		Objc_sendMsg[uintptr](Class_NSCursor, Sel_unhide)
	}
}

func WarpMouseCursorToPoint(x, y float64) {
	point := NSPoint{X: x, Y: y}
	purego.SyscallN(_CGWarpMouseCursorPosition, uintptr(unsafe.Pointer(&point)))
}

func CreateCustomCursor(img image.Image, hotX, hotY int) (NSCursor, error) {
	return NSCursor{}, fmt.Errorf("CreateCustomCursor is not yet implemented")
}

func SetApplicationIconImageFromImage(img image.Image) error {
	nsImg, err := nsImageFromGoImage(img)
	if err != nil {
		return err
	}
	defer nsImg.Release()

	app, err := NSApp()
	if err != nil {
		return err
	}

	Objc_sendMsg[uintptr](uintptr(app.Ptr), Sel_getUid("setApplicationIconImage:"), uintptr(nsImg.Ptr))
	return nil
}

func nsImageFromGoImage(img image.Image) (NSImage, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	pixels := rgba.Pix

	repAlloc := Objc_sendMsg[uintptr](Class_NSBitmapImageRep, Sel_alloc)

	var initWithBitmapDataPlanes func(self uintptr, sel Selector, planes *uintptr, width, height, bps, spp int, hasAlpha, isPlanar bool, colorSpaceName uintptr, bytesPerRow, bitsPerPixel int) uintptr
	purego.RegisterLibFunc(&initWithBitmapDataPlanes, libobjc, "objc_msgSend")

	planes := uintptr(unsafe.Pointer(&pixels[0]))
	colorSpace := NSString_WithUTF8String("NSCalibratedRGBColorSpace")

	rep := initWithBitmapDataPlanes(
		repAlloc,
		Sel_getUid("initWithBitmapDataPlanes:pixelsWide:pixelsHigh:bitsPerSample:samplesPerPixel:hasAlpha:isPlanar:colorSpaceName:bytesPerRow:bitsPerPixel:"),
		&planes, width, height, 8, 4, true, false, uintptr(colorSpace.Ptr), 4*width, 32,
	)
	if rep == 0 {
		return NSImage{}, fmt.Errorf("failed to create NSBitmapImageRep")
	}

	nsImgAlloc := Objc_sendMsg[uintptr](Class_NSImage, Sel_alloc)
	nsImgPtr := Objc_sendMsg[uintptr](nsImgAlloc, Sel_init)
	nsImage := NSImage{Object{unsafe.Pointer(nsImgPtr)}}

	Objc_sendMsg[uintptr](uintptr(nsImage.Ptr), Sel_getUid("addRepresentation:"), rep)
	Objc_sendMsg[uintptr](rep, Sel_release)

	return nsImage, nil
}

func (w NSWindow) Frame() NSRect {
	var frameFunc func(uintptr, Selector) NSRect
	purego.RegisterLibFunc(&frameFunc, libobjc, "objc_msgSend")
	return frameFunc(uintptr(w.Ptr), Sel_frame)
}

func mainNSScreen() NSScreen {
	screen := Objc_sendMsg[uintptr](Class_NSScreen, Sel_mainScreen)
	return NSScreen{Object{unsafe.Pointer(screen)}}
}

func (s NSScreen) Frame() NSRect {
	var frameFunc func(uintptr, Selector) NSRect
	purego.RegisterLibFunc(&frameFunc, libobjc, "objc_msgSend")
	return frameFunc(uintptr(s.Ptr), Sel_frame)
}
