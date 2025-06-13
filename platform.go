package main

import (
	"image"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"modernc.org/libc"
)

type keyChange struct {
	Key   int
	State bool
}

var (
	start time.Time

	streamer   MJPEGHandler
	keyChanges []keyChange
	keyLock    sync.Mutex
)

func handleKey(key string, state string) error {
	keyVal, err := strconv.Atoi(key)
	if err != nil {
		return err
	}
	stateVal, err := strconv.Atoi(state)
	if err != nil {
		return err
	}
	stateBool := stateVal != 0

	keyLock.Lock()
	defer keyLock.Unlock()
	keyChanges = append(keyChanges, keyChange{
		Key:   keyVal,
		State: stateBool,
	})
	return nil
}

func DG_Init() {
	log.Printf("DG_Init called\n")
	start = time.Now()
	addr := ":8080"

	mux := http.NewServeMux()

	mux.Handle("GET /stream.mjpg", &streamer)
	mux.HandleFunc("POST /key/{key}/{state}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		state := r.PathValue("state")
		if err := handleKey(key, state); err != nil {
			http.Error(w, "Invalid key or state value", http.StatusBadRequest)
			log.Printf("Error handling key event: %v\n", err)
			return
		}
	})
	mux.Handle("GET /", http.FileServer(http.Dir("./static")))

	go func() {
		log.Printf("Starting HTTP server on %s\n", addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Fatalf("Failed to start HTTP server: %v\n", err)
		}
	}()
}

func DG_DrawFrame(frame *image.RGBA) {
	if _, err := streamer.AddImage(frame); err != nil {
		log.Printf("Error adding image to MJPEG stream: %v\n", err)
	}
}

func DG_SleepMs(ms uint32) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func DG_GetTicksMs() (r int64) {
	since := time.Since(start)
	return since.Milliseconds()
}

func DG_GetKey(event *DoomKeyEvent) bool {
	// This is a stub; actual key handling would depend on the platform and input system.
	keyLock.Lock()
	defer keyLock.Unlock()
	//log.Printf("DG_GetKey called with pressed: %d, doomKey: %d outstanding entries %d\n", pressed, doomKey, len(keyChanges))
	if len(keyChanges) > 0 {
		change := keyChanges[0]
		keyChanges = keyChanges[1:]
		log.Printf("Processing key change: key=%d, state=%t\n", change.Key, change.State)

		var thisDoomKey int32
		switch change.Key {
		case 38:
			thisDoomKey = key_up
		case 40:
			thisDoomKey = key_down
		case 37:
			thisDoomKey = key_left
		case 39:
			thisDoomKey = key_right
		case 17: // Ctrl
			thisDoomKey = key_fire
		case 32:
			thisDoomKey = key_use
		case 13:
			thisDoomKey = key_menu_forward
		case 27:
			thisDoomKey = key_menu_activate
		default:
			log.Printf("Unknown key %d, ignoring", change.Key)
			return false
		}

		event.Pressed = change.State
		event.Key = uint8(thisDoomKey)
		return true
	}
	return false
}

func DG_SetWindowTitle(title uintptr) {
	log.Printf("DG_SetWindowTitle called with title: %s\n", libc.GoString(title))
	// This is a stub; actual window title setting would depend on the platform and windowing system.
}
