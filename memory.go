package darwin

import (
	"unsafe"
)

func NewAutoreleasePool() NSAutoreleasePool {
	pool := Objc_sendMsg[uintptr](Class_NSAutoreleasePool, Sel_alloc)
	pool = Objc_sendMsg[uintptr](pool, Sel_init)
	return NSAutoreleasePool{Object{unsafe.Pointer(pool)}}
}

func (p NSAutoreleasePool) Drain() {
	if p.Ptr != nil {
		Objc_sendMsg[uintptr](uintptr(p.Ptr), Sel_drain)
	}
}

func (o Object) Retain() {
	if o.Ptr != nil {
		Objc_sendMsg[uintptr](uintptr(o.Ptr), Sel_retain)
	}
}

func (o Object) Release() {
	if o.Ptr != nil {
		Objc_sendMsg[uintptr](uintptr(o.Ptr), Sel_release)
	}
}

func (o Object) Autorelease() {
	if o.Ptr != nil {
		Objc_sendMsg[uintptr](uintptr(o.Ptr), Sel_autorelease)
	}
}
