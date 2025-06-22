package gore

import (
	"fmt"
	"image"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"
)

type delayedKeyEvent struct {
	event DoomKeyEvent
	ticks int32 // How many game ticks before we trigger this event, since the last one
}

type doomTestHeadless struct {
	start         time.Time
	t             *testing.T
	keys          []delayedKeyEvent
	lastEventTick int32
	outputFile    io.WriteCloser
	lock          sync.Mutex
}

func (d *doomTestHeadless) Close() {
	if err := d.outputFile.Close(); err != nil {
		d.t.Errorf("Error closing output file: %v", err)
	}
}

func ffmpegSaver(filename string) (io.WriteCloser, error) {
	args := []string{
		"ffmpeg",
		"-y", // Overwrite output file if it exists
		"-loglevel", "error",
		"-hide_banner",
		"-f", "rawvideo",
		"-s", fmt.Sprintf("%dx%d", SCREENWIDTH, SCREENHEIGHT),
		"-r", "29", // Frame rate - 35ms ticks?
		"-pix_fmt", "rgba",
		"-i", "-",
		"-crf", "27",
		"-preset", "veryfast",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-flush_packets", "1",
		"-movflags", "+faststart",
		filename,
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return stdin, nil
}

func (d *doomTestHeadless) DrawFrame(frame *image.RGBA) {
	if d.outputFile == nil {
		var err error
		name := fmt.Sprintf("doom_test_%s.mp4", d.t.Name())
		d.outputFile, err = ffmpegSaver(name)
		if err != nil {
			d.t.Fatalf("Error starting ffmpeg: %v", err)
		}
		d.t.Logf("Saving output to %s", name)
	}
	d.outputFile.Write(frame.Pix)
}

func (d *doomTestHeadless) SetTitle(title string) {
	d.t.Logf("SetTitle called with: %s", title)
}

func (d *doomTestHeadless) GetKey(event *DoomKeyEvent) bool {
	if len(d.keys) == 0 {
		return false
	}
	d.lock.Lock()
	defer d.lock.Unlock()
	now := I_GetTimeMS()
	delta := now - d.lastEventTick
	if d.keys[0].ticks <= delta {
		*event = d.keys[0].event
		//d.t.Logf("Key event: %#v, delta=%d (%d remaining)", *event, delta, len(d.keys)-1)
		d.keys = d.keys[1:]
		d.lastEventTick = now
		return true
	}
	return false
}

// InsertKey simulates an immediate key press and release event in the game.
func (d *doomTestHeadless) InsertKey(Key uint8) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.keys = append(d.keys, delayedKeyEvent{
		event: DoomKeyEvent{
			Pressed: true,
			Key:     Key,
		},
		ticks: 0, // Insert immediately
	},
		delayedKeyEvent{
			event: DoomKeyEvent{
				Pressed: false,
				Key:     Key,
			},
			ticks: 2, // Release after 1 tick
		},
	)
}

func (d *doomTestHeadless) InsertKeyChange(Key uint8, pressed bool) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.keys = append(d.keys, delayedKeyEvent{
		event: DoomKeyEvent{
			Pressed: pressed,
			Key:     Key,
		},
		ticks: 0, // Insert immediately
	})
}

// Run the demo at super speed to make sure it all goes ok
func TestDoomDemo(t *testing.T) {
	dg_speed_ratio = 100.0 // Run at 50x speed
	game := &doomTestHeadless{
		t:     t,
		start: time.Now(),
		keys: []delayedKeyEvent{
			{DoomKeyEvent{Pressed: true, Key: KEY_ESCAPE}, 60_000},
			{DoomKeyEvent{Pressed: false, Key: KEY_ESCAPE}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_UPARROW1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_UPARROW1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: true, Key: 'y'}, 100},
			{DoomKeyEvent{Pressed: false, Key: 'y'}, 100},
		},
	}
	defer game.Close()
	Run(game, []string{"-iwad", "doom1.wad"})
}

func TestLoadSave(t *testing.T) {
	dg_speed_ratio = 100.0 // Run at 50x speed
	game := &doomTestHeadless{
		t:     t,
		start: time.Now(),
		keys: []delayedKeyEvent{
			// Start a new game
			{DoomKeyEvent{Pressed: true, Key: KEY_ESCAPE}, 5_000},
			{DoomKeyEvent{Pressed: false, Key: KEY_ESCAPE}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_ENTER}, 100},
			// Move around in the game a bit
			{DoomKeyEvent{Pressed: true, Key: KEY_UPARROW1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_UPARROW1}, 1000},
			{DoomKeyEvent{Pressed: true, Key: KEY_LEFTARROW1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_LEFTARROW1}, 1000},
			// Go to the menu and save
			{DoomKeyEvent{Pressed: true, Key: KEY_ESCAPE}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_ESCAPE}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_DOWNARROW1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_DOWNARROW1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_DOWNARROW1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_DOWNARROW1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_DOWNARROW1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_DOWNARROW1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_ENTER}, 100},
			// Clear the previous name and enter a new one
			{DoomKeyEvent{Pressed: true, Key: KEY_BACKSPACE1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_BACKSPACE1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_BACKSPACE1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_BACKSPACE1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_BACKSPACE1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_BACKSPACE1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_BACKSPACE1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_BACKSPACE1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_BACKSPACE1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_BACKSPACE1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_BACKSPACE1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_BACKSPACE1}, 100},
			{DoomKeyEvent{Pressed: true, Key: 't'}, 100},
			{DoomKeyEvent{Pressed: false, Key: 't'}, 100},
			{DoomKeyEvent{Pressed: true, Key: 'e'}, 100},
			{DoomKeyEvent{Pressed: false, Key: 'e'}, 100},
			{DoomKeyEvent{Pressed: true, Key: 's'}, 100},
			{DoomKeyEvent{Pressed: false, Key: 's'}, 100},
			{DoomKeyEvent{Pressed: true, Key: 't'}, 100},
			{DoomKeyEvent{Pressed: false, Key: 't'}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_ENTER}, 100},
			// Start a new game
			{DoomKeyEvent{Pressed: true, Key: KEY_ESCAPE}, 1000},
			{DoomKeyEvent{Pressed: false, Key: KEY_ESCAPE}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_UPARROW1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_UPARROW1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_UPARROW1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_UPARROW1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_UPARROW1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_UPARROW1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_ENTER}, 100},
			// Load the saved game
			{DoomKeyEvent{Pressed: true, Key: KEY_ESCAPE}, 1000},
			{DoomKeyEvent{Pressed: false, Key: KEY_ESCAPE}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_DOWNARROW1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_DOWNARROW1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_DOWNARROW1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_DOWNARROW1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_ENTER}, 100},
			// Quit
			{DoomKeyEvent{Pressed: true, Key: KEY_ESCAPE}, 1000},
			{DoomKeyEvent{Pressed: false, Key: KEY_ESCAPE}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_DOWNARROW1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_DOWNARROW1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_DOWNARROW1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_DOWNARROW1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_DOWNARROW1}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_DOWNARROW1}, 100},
			{DoomKeyEvent{Pressed: true, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: false, Key: KEY_ENTER}, 100},
			{DoomKeyEvent{Pressed: true, Key: 'y'}, 100},
			{DoomKeyEvent{Pressed: false, Key: 'y'}, 100},
		},
	}
	defer game.Close()
	Run(game, []string{"-iwad", "doom1.wad"})
}

func TestDoomRandom(t *testing.T) {
	dg_speed_ratio = 100.0 // Run at 50x speed
	game := &doomTestHeadless{
		t:     t,
		start: time.Now(),
		keys:  nil,
	}
	defer game.Close()
	go func() {
		// Let things get settled
		time.Sleep(20 * time.Millisecond)
		// Start a game
		game.InsertKey(KEY_ESCAPE) // Open menu
		time.Sleep(2 * time.Millisecond)
		game.InsertKey(KEY_ENTER)
		time.Sleep(2 * time.Millisecond)
		game.InsertKey(KEY_ENTER)
		time.Sleep(2 * time.Millisecond)
		game.InsertKey(KEY_ENTER) // Start new game

		time.Sleep(10 * time.Millisecond)
		keys := []uint8{
			KEY_UPARROW1, KEY_DOWNARROW1, KEY_LEFTARROW1,
			KEY_RIGHTARROW1, KEY_FIRE1, KEY_USE1,
		}
		// Press shift to run
		game.InsertKeyChange(0x80+0x36, true)
		// Do 10 seconds of random movement
		for i := 0; i < 1000; i++ {
			change := rand.Intn(len(keys))
			key := keys[change]
			game.InsertKeyChange(key, true)
			time.Sleep(2 * time.Millisecond)
			game.InsertKeyChange(key, false)
			time.Sleep(1 * time.Millisecond)
		}

		// Exit
		game.InsertKey(KEY_ESCAPE) // Open menu
		time.Sleep(10 * time.Millisecond)
		game.InsertKey(KEY_UPARROW1) // Go to quit
		time.Sleep(10 * time.Millisecond)
		game.InsertKey(KEY_ENTER) // Confirm quit
		time.Sleep(10 * time.Millisecond)
		game.InsertKey('y') // Confirm exit
		time.Sleep(10 * time.Millisecond)
	}()
	Run(game, []string{"-iwad", "doom1.wad"})
}
