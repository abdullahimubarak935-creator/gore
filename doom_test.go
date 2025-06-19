package gore

import (
	"image"
	"testing"
	"time"
)

type delayedKeyEvent struct {
	event DoomKeyEvent
	delay time.Duration
}

type doomTestHeadless struct {
	start time.Time
	t     *testing.T
	keys  []delayedKeyEvent
}

func (d *doomTestHeadless) DrawFrame(frame *image.RGBA) {

}

func (d *doomTestHeadless) SetTitle(title string) {
	d.t.Logf("SetTitle called with: %s", title)
}

func (d *doomTestHeadless) GetKey(event *DoomKeyEvent) bool {
	for i, keyEvent := range d.keys {
		if time.Since(d.start) >= keyEvent.delay {
			*event = keyEvent.event
			d.t.Logf("Key event: %v", keyEvent.event)
			d.keys = append(d.keys[:i], d.keys[i+1:]...)
			return true
		}
	}
	return false
}

// Run the demo at super speed to make sure it all goes ok
func TestDoomDemo(t *testing.T) {
	dg_speed_ratio = 50.0 // Run at 50x speed
	wait := time.Second * 10
	game := &doomTestHeadless{
		t:     t,
		start: time.Now(),
		keys: []delayedKeyEvent{
			{DoomKeyEvent{Pressed: true, Key: KEY_ESCAPE}, wait},
			{DoomKeyEvent{Pressed: false, Key: KEY_ESCAPE},
				wait + 100*time.Millisecond},
			{DoomKeyEvent{Pressed: true, Key: KEY_UPARROW1},
				wait + 200*time.Millisecond},
			{DoomKeyEvent{Pressed: false, Key: KEY_UPARROW1},
				wait + 300*time.Millisecond},
			{DoomKeyEvent{Pressed: true, Key: KEY_ENTER},
				wait + 400*time.Millisecond},
			{DoomKeyEvent{Pressed: false, Key: KEY_ENTER},
				wait + 500*time.Millisecond},
			{DoomKeyEvent{Pressed: true, Key: 'y'},
				wait + 600*time.Millisecond},
			{DoomKeyEvent{Pressed: false, Key: 'y'},
				wait + 700*time.Millisecond},
		},
	}
	Run(game, 1, 0)
}
