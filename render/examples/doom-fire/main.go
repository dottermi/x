package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"

	"github.com/dottermi/x/render"
)

// Fire palette: black → red → orange → yellow → white
var palette = []render.Color{
	render.RGB(0, 0, 0),       // 0 - terminal default (transparent)
	render.RGB(7, 7, 7),       // 1
	render.RGB(31, 7, 7),      // 2
	render.RGB(47, 15, 7),     // 3
	render.RGB(71, 15, 7),     // 4
	render.RGB(87, 23, 7),     // 5
	render.RGB(103, 31, 7),    // 6
	render.RGB(119, 31, 7),    // 7
	render.RGB(143, 39, 7),    // 8
	render.RGB(159, 47, 7),    // 9
	render.RGB(175, 63, 7),    // 10
	render.RGB(191, 71, 7),    // 11
	render.RGB(199, 71, 7),    // 12
	render.RGB(223, 79, 7),    // 13
	render.RGB(223, 87, 7),    // 14
	render.RGB(223, 87, 7),    // 15
	render.RGB(215, 95, 7),    // 16
	render.RGB(215, 95, 7),    // 17
	render.RGB(215, 103, 15),  // 18
	render.RGB(207, 111, 15),  // 19
	render.RGB(207, 119, 15),  // 20
	render.RGB(207, 127, 15),  // 21
	render.RGB(207, 135, 23),  // 22
	render.RGB(199, 135, 23),  // 23
	render.RGB(199, 143, 23),  // 24
	render.RGB(199, 151, 31),  // 25
	render.RGB(191, 159, 31),  // 26
	render.RGB(191, 159, 31),  // 27
	render.RGB(191, 167, 39),  // 28
	render.RGB(191, 167, 39),  // 29
	render.RGB(191, 175, 47),  // 30
	render.RGB(183, 175, 47),  // 31
	render.RGB(183, 183, 47),  // 32
	render.RGB(183, 183, 55),  // 33
	render.RGB(207, 207, 111), // 34
	render.RGB(223, 223, 159), // 35
	render.RGB(239, 239, 199), // 36
	render.RGB(255, 255, 255), // 37 - white
}

type fireState struct {
	width, height int
	fire          []int
}

func newFireState(width, height int) *fireState {
	// Double height for half-block rendering
	fireHeight := height * 2
	fire := make([]int, width*fireHeight)

	// Initialize bottom row to max intensity
	for x := 0; x < width; x++ {
		fire[(fireHeight-1)*width+x] = len(palette) - 1
	}

	return &fireState{
		width:  width,
		height: fireHeight,
		fire:   fire,
	}
}

func (f *fireState) resize(width, height int) {
	fireHeight := height * 2
	newFire := make([]int, width*fireHeight)

	// Copy existing fire data where possible
	for y := 0; y < min(f.height, fireHeight); y++ {
		for x := 0; x < min(f.width, width); x++ {
			newFire[y*width+x] = f.fire[y*f.width+x]
		}
	}

	// Initialize bottom row to max intensity
	for x := 0; x < width; x++ {
		newFire[(fireHeight-1)*width+x] = len(palette) - 1
	}

	f.width = width
	f.height = fireHeight
	f.fire = newFire
}

func (f *fireState) spread() {
	for x := 0; x < f.width; x++ {
		for y := 1; y < f.height; y++ {
			src := y*f.width + x
			pixel := f.fire[src]

			if pixel == 0 {
				f.fire[src-f.width] = 0
			} else {
				decay := rand.Intn(3)
				drift := rand.Intn(3) - 1

				dstX := x + drift
				if dstX < 0 {
					dstX = 0
				} else if dstX >= f.width {
					dstX = f.width - 1
				}

				dst := (y-1)*f.width + dstX
				newVal := pixel - decay
				if newVal < 0 {
					newVal = 0
				}
				f.fire[dst] = newVal
			}
		}
	}
}

func (f *fireState) render() *render.Buffer {
	bufHeight := f.height / 2
	buf := render.NewBuffer(f.width, bufHeight)

	for y := 0; y < bufHeight; y++ {
		for x := 0; x < f.width; x++ {
			topIdx := (y*2)*f.width + x
			botIdx := (y*2+1)*f.width + x

			topColor := palette[f.fire[topIdx]]
			botColor := palette[f.fire[botIdx]]

			cell := render.Cell{
				Char: '▀',
				FG:   topColor,
				BG:   botColor,
			}
			buf.Set(x, y, cell)
		}
	}

	return buf
}

func getTermSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80, 24 // fallback
	}
	return width, height
}

func main() {
	width, height := getTermSize()
	fire := newFireState(width, height)
	terminal := render.NewTerminal(width, height)

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGWINCH)

	// Enter alternate screen buffer + hide cursor
	fmt.Print(render.AltScreenEnter)
	fmt.Print(render.HideCursor())
	fmt.Print(render.ClearScreen)

	// Cleanup on exit
	defer func() {
		fmt.Print(render.ShowCursor())
		fmt.Print(render.AltScreenLeave)
	}()

	ticker := time.NewTicker(33 * time.Millisecond) // ~30 FPS
	defer ticker.Stop()

	for {
		select {
		case sig := <-sigChan:
			if sig == syscall.SIGWINCH {
				// Terminal resized
				width, height = getTermSize()
				fire.resize(width, height)
				terminal.Resize(width, height)
				fmt.Print(render.ClearScreen)
			} else {
				// SIGINT or SIGTERM
				return
			}
		case <-ticker.C:
			fire.spread()
			buf := fire.render()
			output := terminal.Render(buf)
			fmt.Print(output)
		}
	}
}
