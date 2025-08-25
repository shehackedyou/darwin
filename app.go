package darwin

import (
	"fmt"
	"unsafe"
	"visualizer/platform/cocoa"
)

func SetupApplication(appName string, delegate uintptr, menu cocoa.MenuBuilder) (Object, error) {
	app, err := NSApp()
	if err != nil {
		return Object{}, err
	}
	appPtr := uintptr(app.Ptr)

	Objc_sendMsg[uintptr](appPtr, Sel_setActivationPolicy, 0) // NSApplicationActivationPolicyRegular

	mainMenu := Objc_alloc_init(Class_NSMenu)

	buildTopLevelMenu := func(title string, items []cocoa.MenuItem) {
		menuItem := Objc_alloc_init(Class_NSMenuItem)
		Objc_sendMsg[uintptr](mainMenu, Sel_addItem, menuItem)
		subMenu := Objc_alloc_init(Class_NSMenu)
		Objc_sendMsg[uintptr](subMenu, Sel_setTitle, NSString_WithUTF8String(title).Ptr)
		buildMenu(subMenu, items)
		Objc_sendMsg[uintptr](menuItem, Sel_setSubmenu, subMenu)
	}

	appMenuItem := Objc_alloc_init(Class_NSMenuItem)
	Objc_sendMsg[uintptr](mainMenu, Sel_addItem, appMenuItem)
	appMenu := Objc_alloc_init(Class_NSMenu)
	Objc_sendMsg[uintptr](appMenu, Sel_setTitle, NSString_WithUTF8String(appName).Ptr)

	buildMenu(appMenu, menu.AppItems)
	Objc_sendMsg[uintptr](appMenuItem, Sel_setSubmenu, appMenu)

	buildTopLevelMenu("File", menu.FileItems)
	buildTopLevelMenu("Edit", menu.EditItems)
	buildTopLevelMenu("Window", menu.WindowItems)

	servicesMenu := Objc_alloc_init(Class_NSMenu)
	Objc_sendMsg[uintptr](appPtr, Sel_setServicesMenu, servicesMenu)

	Objc_sendMsg[uintptr](appPtr, Sel_setMainMenu, mainMenu)
	Objc_sendMsg[uintptr](appPtr, Sel_setDelegate, delegate)

	return app, nil
}

func buildMenu(nativeMenu uintptr, items []cocoa.MenuItem) {
	for _, item := range items {
		addMenuItem(nativeMenu, item)
	}
}

func addMenuItem(menu uintptr, item cocoa.MenuItem) {
	if item.IsSeparator {
		sep := Objc_sendMsg[uintptr](Class_NSMenuItem, Sel_getUid("separatorItem"))
		Objc_sendMsg[uintptr](menu, Sel_addItem, sep)
		return
	}

	titleStr := NSString_WithUTF8String(item.Title)
	keyStr := NSString_WithUTF8String(item.Key)

	var submenu uintptr
	if item.Submenu != nil {
		submenu = Objc_alloc_init(Class_NSMenu)
		buildMenu(submenu, item.Submenu.AppItems)
	}

	menuItem := Objc_sendMsg[uintptr](Class_NSMenuItem, Sel_alloc)
	menuItem = Objc_sendMsg[uintptr](menuItem, Sel_getUid("initWithTitle:action:keyEquivalent:"), titleStr.Ptr, item.Action, keyStr.Ptr)

	if submenu != 0 {
		Objc_sendMsg[uintptr](menuItem, Sel_setSubmenu, submenu)
	}

	mask := 0
	for _, flag := range item.ModifierFlags {
		mask |= int(flag)
	}
	if mask > 0 {
		Objc_sendMsg[uintptr](menuItem, Sel_getUid("setKeyEquivalentModifierMask:"), uintptr(mask))
	}

	Objc_sendMsg[uintptr](menu, Sel_addItem, menuItem)
}

func NSApp() (Object, error) {
	app := Objc_sendMsg[uintptr](Class_NSApplication, Sel_sharedApplication)
	if app == 0 {
		return Object{}, fmt.Errorf("shared NSApplication instance is nil")
	}
	return Object{unsafe.Pointer(app)}, nil
}

func RunApplication(app Object) {
	Objc_sendMsg[uintptr](uintptr(app.Ptr), Sel_run)
}

func ActivateIgnoringOtherApps(app Object) {
	Objc_sendMsg[uintptr](uintptr(app.Ptr), Sel_activateIgnoringOtherApps, 1)
}
