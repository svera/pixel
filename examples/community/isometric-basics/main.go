package main

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	_ "image/png"
)

const (
	windowWidth  = 800
	windowHeight = 800
	// sprite tiles are squared, 64x64 size
	tileSize = 64
	e        = 0 // empty
	f        = 1 // floor identifier
	w        = 2 // wall identifier
	// This is the scrolling speed
	speed = 120
)

var levelData = [][][]uint{
	{
		{f, f, f, f, f, w}, // This row will be rendered in the lower left part of the screen (closer to the viewer)
		{w, f, f, f, f, w},
		{w, f, f, f, f, w},
		{w, f, f, f, f, w},
		{w, f, f, f, f, w},
		{w, w, w, w, w, w}, // And this in the upper right
	},
	{
		{0, 0, 0, 0, 0, w},
		{w, 0, 0, 0, 0, w},
		{w, 0, 0, 0, 0, w},
		{w, 0, 0, 0, 0, w},
		{w, 0, 0, 0, 0, w},
		{w, w, w, w, w, w},
	},
}
var win *pixelgl.Window
var offset = pixel.V(400, 325)
var floorTile, wallTile, walkerSprite *pixel.Sprite
var walkerCartesianPos = pixel.V(1, 1)
var walkerIsoPos = pixel.V(0, 0)

var movementMode = 1

func run() {
	var err error

	cfg := pixelgl.WindowConfig{
		Title:  "Isometric demo",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}
	win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	pic, err := loadPicture("castle.png")
	if err != nil {
		panic(err)
	}
	walker, err := loadPicture("walker2.png")
	if err != nil {
		panic(err)
	}

	wallTile = pixel.NewSprite(pic, pixel.R(0, 448, tileSize, 512))
	floorTile = pixel.NewSprite(pic, pixel.R(0, 128, tileSize, 192))
	walkerSprite = pixel.NewSprite(walker, pixel.R(0, 0, tileSize, tileSize))

	renderAllFloors()
	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		d := speed * dt
		last = time.Now()
		v := pixel.V(0, 0)

		if win.Pressed(pixelgl.KeyW) {
			v.Y = d
			offset = offset.Add(v)
		}
		if win.Pressed(pixelgl.KeyS) {
			v.Y = -d
			offset = offset.Add(v)
		}
		if win.Pressed(pixelgl.KeyD) {
			v.X = d
			offset = offset.Add(v)
		}
		if win.Pressed(pixelgl.KeyA) {
			v.X = -d
			offset = offset.Add(v)
		}
		if win.Pressed(pixelgl.KeyUp) {
			v.Y = d
			walkerIsoPos = walkerIsoPos.Add(v)
		}
		if win.Pressed(pixelgl.KeyDown) {
			v.Y = -d
			walkerIsoPos = walkerIsoPos.Add(v)
		}
		if win.Pressed(pixelgl.KeyRight) {
			v.X = d
			walkerIsoPos = walkerIsoPos.Add(v)
		}
		if win.Pressed(pixelgl.KeyLeft) {
			v.X = -d
			walkerIsoPos = walkerIsoPos.Add(v)
		}

		// Check if screen must be re-rendered
		if v.X != 0 || v.Y != 0 {
			renderAllFloors()
		}

		win.Update()
	}
}

// Draw level data tiles to window, from farthest to closest.
// In order to achieve the depth effect, we need to render tiles up to down, being lower
// closer to the viewer (see painter's algorithm). To do that, we need to process levelData in reverse order,
// so its first row is rendered last, as OpenGL considers its origin to be in the lower left corner of the display.
func renderFloor(floor int) {
	height := floor * (tileSize / 2)
	for x := len(levelData[floor]) - 1; x >= 0; x-- {
		for y := len(levelData[floor][x]) - 1; y >= 0; y-- {
			isoCoords := cartesianToIso(pixel.V(float64(x), float64(y)))
			isoCoords.Y += float64(height)
			mat := pixel.IM.Moved(offset.Add(isoCoords))
			tileType := levelData[floor][x][y]
			if tileType == f {
				floorTile.Draw(win, mat)
			} else if tileType == w {
				wallTile.Draw(win, mat)
			}
			if walkerCartesianPos.X == float64(x) && walkerCartesianPos.Y == float64(y) && floor == 0 {
				walkerSprite.Draw(win, mat.Moved(walkerIsoPos))
			}
		}
	}
}

func renderAllFloors() {
	win.Clear(colornames.Black)
	for n := 0; n < len(levelData); n++ {
		renderFloor(n)
	}
}

func cartesianToIso(pt pixel.Vec) pixel.Vec {
	return pixel.V((pt.X-pt.Y)*(tileSize/2), (pt.X+pt.Y)*(tileSize/4))
}

func isoToCartesian(pt pixel.Vec) pixel.Vec {
	x := (pt.X/(tileSize/2) + pt.Y/(tileSize/4)) / 2
	y := (pt.Y/(tileSize/4) - (pt.X / (tileSize / 2))) / 2
	return pixel.V(x, y)
}

func main() {
	pixelgl.Run(run)
}
