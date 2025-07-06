package main

import (
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"time"

	"github.com/AndreRenaud/gore"
	"github.com/nfnt/resize"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

type termDoom struct {
	outstandingKeys map[uint8]time.Time
}

// Characters with progressively fewer filled pixels, to simulate brightness
const brightChars = "$@B%#*+=\"~^;:..."

func ascii(img *image.RGBA, writer io.Writer) {
	bound := img.Bounds()
	height, width := bound.Max.Y, bound.Max.X
	var lastcolor string

	for y := bound.Min.Y; y < height; y++ {
		for x := bound.Min.X; x < width; x++ {
			offset := y*img.Stride + x*4
			r := img.Pix[offset] & 0xf0
			g := img.Pix[offset+1] & 0xf0
			b := img.Pix[offset+2] & 0xf0
			colStr := fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
			if colStr != lastcolor {
				writer.Write([]byte(colStr))
				lastcolor = colStr
			}

			brightness := int(r+g+b) * (len(brightChars) - 1) / (3 * 255) // Normalize to 0-1 range
			writer.Write([]byte{brightChars[brightness]})
		}
		writer.Write([]byte("\r\n"))
	}
}

func (t *termDoom) DrawFrame(frame *image.RGBA) {
	// Fit the frame to the terminal size
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 120
	}
	// Make sure we're not going to scroll off the bottom of the visible terminal
	width = min(height*320/200*2, width)

	height = (width * 200 / 320) / 2 // fixed width fonts are typically twice as high as wide
	smaller := resize.Resize(uint(width), uint(height), frame, resize.Lanczos3)
	// Go back to 0,0
	fmt.Print("\033[0;0H")
	rgba, ok := smaller.(*image.RGBA)
	if !ok {
		log.Printf("Error: resized image is not of type *image.RGBA")
		return
	}
	ascii(rgba, os.Stdout)
	os.Stdout.Sync()
}

func (t *termDoom) GetEvent(event *gore.DoomEvent) bool {
	for key, lastTime := range t.outstandingKeys {
		if time.Since(lastTime) > 100*time.Millisecond {
			delete(t.outstandingKeys, key)
			event.Type = gore.Ev_keyup
			event.Key = key
			return true
		}
	}
	var buf [3]byte
	fds := []unix.PollFd{
		{Fd: int32(os.Stdin.Fd()), Events: unix.POLLIN},
	}
	if n, err := unix.Poll(fds, 0); err != nil || n <= 0 {
		return false // No event handled
	}
	n, err := os.Stdin.Read(buf[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
	}
	switch string(buf[:n]) {
	case "\x1b[A": // Up arrow
		event.Type = gore.Ev_keydown
		event.Key = gore.KEY_UPARROW1
	case "\x1b[B": // Down arrow
		event.Type = gore.Ev_keydown
		event.Key = gore.KEY_DOWNARROW1
	case "\x1b[C": // Right arrow
		event.Type = gore.Ev_keydown
		event.Key = gore.KEY_RIGHTARROW1
	case "\x1b[D": // Left arrow
		event.Type = gore.Ev_keydown
		event.Key = gore.KEY_LEFTARROW1
	case " ":
		event.Type = gore.Ev_keydown
		event.Key = gore.KEY_USE1
	case "\n", "\r":
		event.Type = gore.Ev_keydown
		event.Key = gore.KEY_ENTER
	case "\x1b": // Escape
		event.Type = gore.Ev_keydown
		event.Key = gore.KEY_ESCAPE
	case ",": // We can't use Ctrl, since terminals don't see it. So we use comma instead
		event.Type = gore.Ev_keydown
		event.Key = gore.KEY_FIRE1
	case "y", "n", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
		event.Type = gore.Ev_keydown
		event.Key = buf[0]
	case "\t":
		event.Type = gore.Ev_keydown
		event.Key = gore.KEY_TAB
	default:
		return false

	}
	t.outstandingKeys[event.Key] = time.Now()
	return true
}

func (t *termDoom) SetTitle(title string) {
	fmt.Printf("\033]0;%s\007", title)
}

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Initialize termDoom and start the game loop
	termGame := &termDoom{
		outstandingKeys: make(map[uint8]time.Time),
	}

	gore.Run(termGame, os.Args[1:])
}
