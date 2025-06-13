package main

import (
	"log"
	"time"

	"modernc.org/libc"
)

var (
	start time.Time
)

func DG_Init(tls *libc.TLS) {
	log.Printf("DG_Init called\n")
	start = time.Now()
}

func DG_DrawFrame(tls *libc.TLS) {
	log.Printf("DG_DrawFrame called\n")
}

func DG_SleepMs(tls *libc.TLS, ms uint32) {
	log.Printf("DG_SleepMs called with %d ms\n", ms)
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func DG_GetTicksMs(tls *libc.TLS) (r int64) {
	return time.Since(start).Milliseconds()
}

func DG_GetKey(tls *libc.TLS, pressed uintptr, doomKey uintptr) (r int32) {
	log.Printf("DG_GetKey called with pressed: %d, doomKey: %d\n", pressed, doomKey)
	// This is a stub; actual key handling would depend on the platform and input system.
	return 0 // Return 0 to indicate no key pressed
}

func DG_SetWindowTitle(tls *libc.TLS, title uintptr) {
	log.Printf("DG_SetWindowTitle called with title: %s\n", libc.GoString(title))
	// This is a stub; actual window title setting would depend on the platform and windowing system.
}
