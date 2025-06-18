# 🔥 GORE 🔥
## A Minimal Platform-Agnostic Go Port of doomgeneric

```
    ██████╗  ██████╗  ██████╗ ███╗   ███╗
    ██╔══██╗██╔═══██╗██╔═══██╗████╗ ████║
    ██║  ██║██║   ██║██║   ██║██╔████╔██║
    ██║  ██║██║   ██║██║   ██║██║╚██╔╝██║
    ██████╔╝╚██████╔╝╚██████╔╝██║ ╚═╝ ██║
    ╚═════╝  ╚═════╝  ╚═════╝ ╚═╝     ╚═╝
                    .GO
```

**"Rip and tear... in Go!"**

> *The demons thought they were safe when they corrupted the C codebase.  
> They were wrong.  
> The Doom Slayer has learned Go.*

---

## 💀 WHAT FRESH HELL IS THIS?

This is a **minimal, platform-agnostic Go port** of the legendary DOOM engine, transpiled from the `doomgeneric` codebase. No CGo. No platform dependencies. Just pure, unadulterated demon-slaying action powered by the glory of Go's cross-compilation.

The original C code has been mechanically converted to Go using [advanced transpilation sorcery](https://gitlab.com/cznic/doomgeneric.git)

## 🔫 FEATURES

- ✅ **Platform Agnostic**: Runs anywhere Go runs
- ✅ **Minimal Dependencies**: Only requires Go standard library
- ✅ **Multiple DOOM Versions**: Supports DOOM, DOOM II, Ultimate DOOM, Final DOOM
- ✅ **WAD File Support**: Bring your own demons via WAD files
- ✅ **Memory Safe**: Go's GC protects you from buffer overflows (but not from Cacodemons)
- ✅ **Cross Compilation**: Build for any target from any platform

## 🚀 INSTALLATION

### Prerequisites
- Go 1.24+ (The demons fear modern Go)
- A WAD file containing the forces of Hell

### Running the examples
#### Web based
```bash
git clone <this-repo>
cd gore
go run ./example/webserver
```
Now browse to http://localhost:8080 to play

#### Ebitengine
```bash
go run ./example/ebitengine
```
The window should pop up to run Doom

### Getting WAD Files
You need the game data files (WAD) to run DOOM:
- **Shareware**: Download `doom1.wad` (free)
- **Retail**: Use your legally owned copy of DOOM.WAD or doom2.wad
- **Ultimate DOOM**: doom.wad from Ultimate DOOM
- **Final DOOM**: tnt.wad or plutonia.wad

## 🔧 PLATFORM IMPLEMENTATION

Similar to `doomgeneric`, the actual input/output is provided externally. The following interface is required:
```go
type DoomFrontend interface {
	DrawFrame(img *image.RGBA)
	SetTitle(title string)
	GetKey(event *DoomKeyEvent) bool
}
```

| Function | Purpose |
|----------|---------|
| `DrawFrame()` | Render the frame to your display |
| `SetTitle()` | Set the window title as appropriate to the given WAD |
| `GetKey()` | Handle keyboard input |

## 🛠️ EXTENDING THIS PORT

Want to make this actually playable? Here's what you need:

1. **Display**: Replace PNG output with SDL2, OpenGL, or terminal output
2. **Input**: Implement proper keyboard/mouse handling in `DG_GetKey()`
3. **Audio**: Add sound system (optional, purists play in silence)
4. **Packaging**: Bundle with shareware WAD for easy distribution

## ⚡ TECHNICAL NOTES

This port uses:
- **Transpiled C Code**: Mechanical conversion from original DOOM source
- **libc Compatibility**: `modernc.org/libc` for C standard library functions
- **Memory Management**: Go's garbage collector handles memory (safer than malloc/free)
- **Type Safety**: Go's type system prevents many classic C vulnerabilities

## 📜 LICENSE

DOOM source code is released under the GNU General Public License.  
This Go port maintains the same licensing terms.

---

## 🔥 FINAL WORDS

*"In the first age, in the first battle, when the shadows first lengthened, one stood. He chose the path of perpetual torment. In his ravenous hatred he found no peace, and with boiling blood he scoured the umbral plains seeking vengeance against the dark lords who had wronged him. And those that tasted the bite of his sword named him... **The Doom Slayer**."*

Now go forth and **RIP AND TEAR** in Go! 🚀

---

*Built with ❤️ and excessive violence*
