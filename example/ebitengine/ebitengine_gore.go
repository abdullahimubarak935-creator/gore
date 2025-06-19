package main

import (
	"image"
	"log"
	"os"
	"sync"

	"github.com/AndreRenaud/gore"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type DoomGame struct {
	lastFrame *ebiten.Image

	keyEvents   []gore.DoomKeyEvent
	lock        sync.Mutex
	terminating bool
}

func (g *DoomGame) Update() error {
	keys := map[ebiten.Key]uint8{
		ebiten.KeySpace:   gore.KEY_USE1,
		ebiten.KeyEscape:  gore.KEY_ESCAPE,
		ebiten.KeyUp:      gore.KEY_UPARROW1,
		ebiten.KeyDown:    gore.KEY_DOWNARROW1,
		ebiten.KeyLeft:    gore.KEY_LEFTARROW1,
		ebiten.KeyRight:   gore.KEY_RIGHTARROW1,
		ebiten.KeyEnter:   gore.KEY_ENTER,
		ebiten.KeyControl: gore.KEY_FIRE1,
		ebiten.KeyY:       'y',
		ebiten.KeyN:       'n',
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	for key, doomKey := range keys {
		if inpututil.IsKeyJustPressed(key) {
			var event gore.DoomKeyEvent

			event.Pressed = true
			event.Key = doomKey
			g.keyEvents = append(g.keyEvents, event)
		} else if inpututil.IsKeyJustReleased(key) {
			var event gore.DoomKeyEvent
			event.Pressed = false
			event.Key = doomKey
			g.keyEvents = append(g.keyEvents, event)
		}
	}
	if g.terminating {
		return ebiten.Termination
	}
	return nil
}

func (g *DoomGame) Draw(screen *ebiten.Image) {
	g.lock.Lock()
	defer g.lock.Unlock()
	screen.DrawImage(g.lastFrame, nil)
}

func (g *DoomGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *DoomGame) GetKey(event *gore.DoomKeyEvent) bool {
	g.lock.Lock()
	defer g.lock.Unlock()
	if len(g.keyEvents) > 0 {
		*event = g.keyEvents[0]
		g.keyEvents = g.keyEvents[1:]
		return true
	}
	return false
}

func (g *DoomGame) DrawFrame(frame *image.RGBA) {
	g.lock.Lock()
	op := &ebiten.DrawImageOptions{}
	rect := frame.Bounds()
	yScale := float64(screenHeight) / float64(rect.Dy())
	xScale := float64(screenWidth) / float64(rect.Dx())
	op.GeoM.Scale(xScale, yScale)
	op.GeoM.Translate(xScale, yScale)
	g.lastFrame.DrawImage(ebiten.NewImageFromImage(frame), op)
	g.lock.Unlock()
}

func (g *DoomGame) SetTitle(title string) {
	ebiten.SetWindowTitle(title)
}

func main() {
	game := &DoomGame{}
	game.lastFrame = ebiten.NewImage(screenWidth, screenHeight)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Gamepad (Ebitengine Demo)")
	ebiten.SetFullscreen(true)
	go func() {
		gore.Run(game, os.Args[1:])
		game.terminating = true
	}()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
