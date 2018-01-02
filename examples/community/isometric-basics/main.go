package main

import (
	"image"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

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

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

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

	levels()

	var newPos = walkerCartesianPos
	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyUp) {
			newPos = walkerCartesianPos.Add(pixel.V(1, 0))
		} else if win.JustPressed(pixelgl.KeyDown) {
			newPos = walkerCartesianPos.Sub(pixel.V(1, 0))
		} else if win.JustPressed(pixelgl.KeyRight) {
			newPos = walkerCartesianPos.Sub(pixel.V(0, 1))
		} else if win.JustPressed(pixelgl.KeyLeft) {
			newPos = walkerCartesianPos.Add(pixel.V(0, 1))
		}
		if isWalkable(0, newPos) {
			walkerCartesianPos = newPos
		}
		levels()
		win.Update()
	}
}

// Draw level data tiles to window, from farthest to closest.
// In order to achieve the depth effect, we need to render tiles up to down, being lower
// closer to the viewer (see painter's algorithm). To do that, we need to process levelData in reverse order,
// so its first row is rendered last, as OpenGL considers its origin to be in the lower left corner of the display.
func depthSort(floor int) {
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
				walkerIsoCoords := cartesianToIso(walkerCartesianPos)
				mat = pixel.IM.Moved(offset.Add(walkerIsoCoords))
				walkerSprite.Draw(win, mat)
			}
		}
	}
}

func levels() {
	for n := 0; n < len(levelData); n++ {
		depthSort(n)
	}
}

func isWalkable(floor uint, pos pixel.Vec) bool {
	if pos.X < 0 || pos.Y < 0 {
		return false
	}
	if levelData[floor][uint(pos.X)][uint(pos.Y)] != w {
		return true
	}
	return false
}

func cartesianToIso(pt pixel.Vec) pixel.Vec {
	return pixel.V((pt.X-pt.Y)*(tileSize/2), (pt.X+pt.Y)*(tileSize/4))
}

func main() {
	pixelgl.Run(run)
}
