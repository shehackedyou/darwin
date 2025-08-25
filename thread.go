package darwin

import (
	"runtime"
	"sync"

	"github.com/ebitengine/purego"
)

var (
	goCallbackFuncs      = make(map[uintptr]func())
	goCallbackFuncsMtx   sync.Mutex
	goCallbackFuncsIndex uintptr
	classGoCallback      uintptr
)

func MainThread(f func()) {
	// If we are already on the main thread, execute the function directly to avoid deadlock.
	runtime.LockOSThread()
	isMain := isMainThread()
	runtime.UnlockOSThread()

	if isMain {
		f()
		return
	}

	// If we are on a different goroutine, dispatch the function to the main
	// thread and wait for it to complete.
	var wg sync.WaitGroup
	wg.Add(1)
	dispatch(func() {
		defer wg.Done()
		f()
	})
	wg.Wait()
}

func isMainThread() bool {
	return Objc_sendMsg[bool](Class_NSThread, Sel_isMainThread)
}

func dispatch(f func()) {
	goCallbackFuncsMtx.Lock()
	goCallbackFuncsIndex++
	idx := goCallbackFuncsIndex
	goCallbackFuncs[idx] = f
	goCallbackFuncsMtx.Unlock()

	cb := Objc_sendMsg[uintptr](classGoCallback, Sel_alloc)
	cb = Objc_sendMsg[uintptr](cb, Sel_init)

	// Wrap the primitive uintptr in an NSNumber object.
	nsIdx := Objc_sendMsg[uintptr](Class_NSNumber, Sel_numberWithInt, idx)

	Objc_sendMsg[uintptr](cb, Sel_performSelectorOnMainThread, Sel_call, nsIdx, 1)

	// We no longer need the manual retain on 'cb'. The system handles it.
	// The balancing release for 'cb' is still in goCallback.
}

func goCallback(id, sel, arg uintptr) {
	// 'arg' is now an NSNumber. We need to extract the integer value.
	idx := Objc_sendMsg[uintptr](arg, Sel_unsignedLongLongValue)

	goCallbackFuncsMtx.Lock()
	f, ok := goCallbackFuncs[idx]
	if ok {
		delete(goCallbackFuncs, idx)
	}
	goCallbackFuncsMtx.Unlock()

	if ok {
		f()
	}

	// Release the callback object now that we're done with it.
	Objc_sendMsg[uintptr](id, Sel_release)
}

func setupGoCallbackClass() {
	className := "GoCallback"
	class := objc_allocateClassPair(Class_NSObject, className, 0)
	if class == 0 {
		panic("failed to allocate GoCallback class")
	}
	classGoCallback = class

	ok := class_addMethod(class, Sel_call, purego.NewCallback(goCallback), "v@:@")
	if !ok {
		panic("failed to add method 'call' to GoCallback")
	}
	objc_registerClassPair(class)
}
