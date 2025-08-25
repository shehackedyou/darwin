### `platform/cocoa/darwin`

A low-level, self-contained bridge to the macOS native API using `purego` for 
direct system calls. A private backend for a higher-level, cross-platform PAL

I released it as a self-contained package because none exist for this purpose
without using 'cgo' and I also abstracted joystick functionality which is also
unique

#### **Structure**
The package is structured to provide direct access to the Objective-C runtime
and various Core Foundation frameworks without external dependencies

* **`init.go`**: Dynamic loading of system frameworks, registers Objective-C classes, and maps selectors for message-passing
* **`types.go`**: Go representations of native structs and Objective-C objects
* **`app.go`**: Native `NSApplication` lifecycle and menu bar creation
* **`window.go`**: Functions for creating and manipulating native `NSWindow` and `NSOpenGLView` objects
* **`events.go`**: Wrappers for retrieving data from native `NSEvent` objects
* **`callbacks.go`**: Go functions that receive callbacks from the Objective-C runtime, bridging native events to Go
* **`thread.go`**: `MainThread` function, a critical utility for dispatching code to the main OS thread as required by the AppKit framework
* **`clipboard.go`**: Clipboard access using `NSPasteboard`
* **`joystick.go`**: IOKit framework to handle joystick and gamepad input
* **`memory.go`**: Objective-C memory management calls (`Retain`, `Release`, `Autorelease`)
* **`helpers.go`**: Utility functions for Objective-C message sending and Go pointer management

#### **Usage**
This package is not intended for direct use by end-user applications
