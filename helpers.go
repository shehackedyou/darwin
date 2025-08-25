package darwin

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

func Objc_sendMsg[R any](receiver uintptr, selector Selector, args ...any) R {
	argPtrs := make([]uintptr, len(args))
	for i, arg := range args {
		switch v := arg.(type) {
		case nil:
			argPtrs[i] = 0
		case uintptr:
			argPtrs[i] = v
		case int:
			argPtrs[i] = uintptr(v)
		case bool:
			if v {
				argPtrs[i] = 1
			} else {
				argPtrs[i] = 0
			}
		case float64:
			argPtrs[i] = *(*uintptr)(unsafe.Pointer(&v))
		case Selector:
			argPtrs[i] = uintptr(v)
		case IOHIDManagerRef:
			argPtrs[i] = uintptr(v)
		case unsafe.Pointer:
			argPtrs[i] = uintptr(v)
		case *byte:
			argPtrs[i] = uintptr(unsafe.Pointer(v))
		case *[]uintptr:
			argPtrs[i] = uintptr(unsafe.Pointer(v))
		case Object:
			argPtrs[i] = uintptr(v.Ptr)
		case NSWindow:
			argPtrs[i] = uintptr(v.Ptr)
		case NSOpenGLView:
			argPtrs[i] = uintptr(v.Ptr)
		case NSOpenGLContext:
			argPtrs[i] = uintptr(v.Ptr)
		case NSString:
			argPtrs[i] = uintptr(v.Ptr)
		case NSRunLoop:
			argPtrs[i] = uintptr(v.Ptr)
		default:
			val := reflect.ValueOf(v)
			if val.Kind() == reflect.Ptr {
				argPtrs[i] = val.Pointer()
			} else {
				panic(fmt.Sprintf("unhandled type in Objc_sendMsg: %T", v))
			}
		}
	}

	allArgs := append([]uintptr{receiver, uintptr(selector)}, argPtrs...)
	ret, _, _ := purego.SyscallN(objc_msgSend, allArgs...)

	var result R
	if unsafe.Sizeof(result) > unsafe.Sizeof(ret) {
		// This path is for struct returns.
		// purego.SyscallN returns a pointer to the struct in the 'ret' variable.
		result = *(*R)(unsafe.Pointer(ret))
	} else {
		// This path is for scalar returns (int, uintptr, bool, etc.).
		result = *(*R)(unsafe.Pointer(&ret))
	}
	return result
}

func Objc_alloc_init(class uintptr) uintptr {
	obj := Objc_sendMsg[uintptr](class, Sel_alloc)
	return Objc_sendMsg[uintptr](obj, Sel_init)
}

func NewCString(s string) *byte {
	b := make([]byte, len(s)+1)
	copy(b, s)
	return &b[0]
}

func GoString(s uintptr) string {
	if s == 0 {
		return ""
	}
	var l int
	for *(*byte)(unsafe.Pointer(s + uintptr(l))) != 0 {
		l++
	}
	return string(unsafe.Slice((*byte)(unsafe.Pointer(s)), l))
}

// Sel_getUid registers a selector name with the Objective-C runtime and returns its handle.
func Sel_getUid(name string) Selector {
	ret, _, _ := purego.SyscallN(uintptr(Sel_registerName), uintptr(unsafe.Pointer(NewCString(name))))
	return Selector(ret)
}

func sel_getName(selector Selector) string {
	cStr, _, _ := purego.SyscallN(sel_getName_ptr, uintptr(selector))
	return GoString(cStr)
}

func class_addMethod(class uintptr, selector Selector, imp uintptr, types string) bool {
	ret, _, _ := purego.SyscallN(class_addMethod_ptr, class, uintptr(selector), imp, uintptr(unsafe.Pointer(NewCString(types))))
	return ret != 0
}

func objc_allocateClassPair(superclass uintptr, name string, extraBytes uintptr) uintptr {
	ret, _, _ := purego.SyscallN(objc_allocateClassPair_ptr, superclass, uintptr(unsafe.Pointer(NewCString(name))), extraBytes)
	return ret
}

func objc_registerClassPair(class uintptr) {
	purego.SyscallN(objc_registerClassPair_ptr, class)
}

// goPointers is a map used to associate a simple integer ID with a Go interface{}.
// This allows us to safely pass a reference to a Go object into the C/Objective-C world.
var (
	goPointers    = make(map[uintptr]any)
	goPointersMtx sync.RWMutex
	goPointerID   uintptr
)

func StoreGoPointer(v any) uintptr {
	goPointersMtx.Lock()
	defer goPointersMtx.Unlock()
	goPointerID++
	ptr := goPointerID
	goPointers[ptr] = v
	return ptr
}

func GetGoPointer(ptr uintptr) any {
	goPointersMtx.RLock()
	defer goPointersMtx.RUnlock()
	return goPointers[ptr]
}

func FreeGoPointer(ptr uintptr) {
	goPointersMtx.Lock()
	defer goPointersMtx.Unlock()
	delete(goPointers, ptr)
}
