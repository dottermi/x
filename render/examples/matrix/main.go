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

var matrixChars = []rune("アイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワヲン0123456789")

var palette = []render.Color{
	render.RGB(0, 40, 0),      // 0 - very dark
	render.RGB(0, 70, 0),      // 1
	render.RGB(0, 100, 0),     // 2
	render.RGB(0, 130, 0),     // 3
	render.RGB(0, 160, 0),     // 4
	render.RGB(0, 190, 0),     // 5
	render.RGB(0, 220, 0),     // 6
	render.RGB(0, 255, 0),     // 7 - bright green
	render.RGB(180, 255, 180), // 8 - white-green (head)
}

type droplet struct {
	col    int
	head   int     // current head position
	tail   int     // current tail position
	endRow int     // where head stops
	length int     // trail length
	speed  float64 // rows per second

	headCrawling bool
	tailCrawling bool
	alive        bool

	lingerTime time.Duration // delay before tail starts
	headStopAt time.Time

	headAccum float64 // accumulated fractional movement
	tailAccum float64

	chars []rune
}

type columnState struct {
	canSpawn    bool
	numDroplets int
}

type matrixState struct {
	width, height   int
	droplets        []*droplet
	columns         []columnState
	maxPerColumn    int
	dropletsPerSec  float64
	spawnAccum      float64
	lastSpawnUpdate time.Time
}

func newMatrixState(width, height int) *matrixState {
	m := &matrixState{
		width:           width,
		height:          height,
		droplets:        make([]*droplet, 0, width*2),
		columns:         make([]columnState, width),
		maxPerColumn:    3,
		dropletsPerSec:  float64(width) * 0.8, // ~80% of columns spawn per second
		lastSpawnUpdate: time.Now(),
	}
	for i := range m.columns {
		m.columns[i].canSpawn = true
	}
	return m
}

func (m *matrixState) resize(width, height int) {
	m.width = width
	m.height = height
	m.columns = make([]columnState, width)
	for i := range m.columns {
		m.columns[i].canSpawn = true
	}
	m.dropletsPerSec = float64(width) * 0.8
	// Keep existing droplets that fit
	alive := m.droplets[:0]
	for _, d := range m.droplets {
		if d.col < width {
			alive = append(alive, d)
		}
	}
	m.droplets = alive
}

func (m *matrixState) spawnDroplet(col int) *droplet {
	length := rand.Intn(15) + 8
	chars := make([]rune, length)
	for i := range chars {
		chars[i] = matrixChars[rand.Intn(len(matrixChars))]
	}

	// Vary speed: 8-20 rows per second
	speed := 8.0 + rand.Float64()*12.0

	d := &droplet{
		col:          col,
		head:         -1,
		tail:         -1 - length,
		endRow:       m.height + rand.Intn(m.height/2), // can go past screen
		length:       length,
		speed:        speed,
		headCrawling: true,
		tailCrawling: false,
		alive:        true,
		lingerTime:   time.Duration(rand.Intn(300)+100) * time.Millisecond,
		chars:        chars,
	}

	m.droplets = append(m.droplets, d)
	m.columns[col].canSpawn = false
	m.columns[col].numDroplets++
	return d
}

func (m *matrixState) spawnDroplets() {
	now := time.Now()
	elapsed := now.Sub(m.lastSpawnUpdate).Seconds()
	m.lastSpawnUpdate = now

	m.spawnAccum += elapsed * m.dropletsPerSec
	toSpawn := int(m.spawnAccum)
	m.spawnAccum -= float64(toSpawn)

	// Find available columns
	available := make([]int, 0, m.width)
	for col, cs := range m.columns {
		if cs.canSpawn && cs.numDroplets < m.maxPerColumn {
			available = append(available, col)
		}
	}

	// Shuffle and spawn
	rand.Shuffle(len(available), func(i, j int) {
		available[i], available[j] = available[j], available[i]
	})

	for i := 0; i < toSpawn && i < len(available); i++ {
		m.spawnDroplet(available[i])
	}
}

func (d *droplet) update(elapsed float64, now time.Time, height int, col *columnState) {
	if !d.alive {
		return
	}

	// Randomly mutate a character
	if rand.Intn(10) == 0 && len(d.chars) > 0 {
		d.chars[rand.Intn(len(d.chars))] = matrixChars[rand.Intn(len(matrixChars))]
	}

	if d.headCrawling {
		d.headAccum += d.speed * elapsed
		advance := int(d.headAccum)
		d.headAccum -= float64(advance)

		d.head += advance
		if d.head >= d.endRow {
			d.head = d.endRow
			d.headCrawling = false
			d.headStopAt = now
		}
		// Re-enable column spawn when head passes 1/3 of screen
		if d.head > height/3 && !col.canSpawn {
			col.canSpawn = true
		}
	}

	if !d.headCrawling && !d.tailCrawling {
		// Wait for linger time
		if now.Sub(d.headStopAt) >= d.lingerTime {
			d.tailCrawling = true
		}
	}

	if d.tailCrawling {
		d.tailAccum += d.speed * elapsed
		advance := int(d.tailAccum)
		d.tailAccum -= float64(advance)

		d.tail += advance
		if d.tail >= d.head {
			d.alive = false
			col.numDroplets--
			if col.numDroplets < 0 {
				col.numDroplets = 0
			}
		}
	}
}

func (m *matrixState) update(elapsed float64) {
	now := time.Now()

	// Spawn new droplets
	m.spawnDroplets()

	// Update existing droplets
	for _, d := range m.droplets {
		d.update(elapsed, now, m.height, &m.columns[d.col])
	}

	// Remove dead droplets
	alive := m.droplets[:0]
	for _, d := range m.droplets {
		if d.alive {
			alive = append(alive, d)
		}
	}
	m.droplets = alive
}

func (m *matrixState) render() *render.Buffer {
	buf := render.NewBuffer(m.width, m.height)

	for _, d := range m.droplets {
		for row := max(0, d.tail); row <= min(d.head, m.height-1); row++ {
			if row < 0 {
				continue
			}

			// Distance from head determines brightness
			distFromHead := d.head - row
			brightness := len(palette) - 2 - distFromHead
			if brightness < 0 {
				brightness = 0
			}

			// Head is brightest
			if row == d.head && d.headCrawling {
				brightness = len(palette) - 1
			}

			charIdx := (row - d.tail) % len(d.chars)
			if charIdx < 0 {
				charIdx = 0
			}

			buf.Set(d.col, row, render.Cell{
				Char: d.chars[charIdx],
				FG:   palette[brightness],
			})
		}
	}

	return buf
}

func getTermSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80, 24
	}
	return width, height
}

func main() {
	width, height := getTermSize()
	matrix := newMatrixState(width, height)
	terminal := render.NewTerminal(width, height)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGWINCH)

	fmt.Print(render.AltScreenEnter)
	fmt.Print(render.HideCursor())
	fmt.Print(render.ClearScreen)

	defer func() {
		fmt.Print(render.ShowCursor())
		fmt.Print(render.AltScreenLeave)
	}()

	const frameTime = 33 * time.Millisecond // ~60 FPS
	lastFrame := time.Now()

	for {
		frameStart := time.Now()
		elapsed := frameStart.Sub(lastFrame).Seconds()
		lastFrame = frameStart

		select {
		case sig := <-sigChan:
			if sig == syscall.SIGWINCH {
				width, height = getTermSize()
				matrix.resize(width, height)
				terminal.Resize(width, height)
				fmt.Print(render.ClearScreen)
			} else {
				return
			}
		default:
		}

		matrix.update(elapsed)
		buf := matrix.render()
		fmt.Print(terminal.Render(buf))

		if sleep := frameTime - time.Since(frameStart); sleep > 0 {
			time.Sleep(sleep)
		}
	}
}
