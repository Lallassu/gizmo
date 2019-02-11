package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"runtime"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Load texture file (could also be spritesheets) (png format)
// Size is the largest of height/width used for 1D pixel arrays
func loadTexture(file string) (img image.Image, width, height, size float64) {
	width = 0.0
	height = 0.0
	size = 0.0
	img = nil
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	imgfile, err := os.Open(file)
	if err != nil {
		Error(fmt.Sprintf("Failed to open file %v", file))
		return
	}

	defer imgfile.Close()

	imgCfg, _, err := image.DecodeConfig(imgfile)
	if err != nil {
		Error(fmt.Sprintf("Failed to decode file %v: %v", file, err))
		return
	}

	imgfile.Seek(0, 0)
	img, _, _ = image.Decode(imgfile)

	height = float64(imgCfg.Height)
	width = float64(imgCfg.Width)
	size = width
	if width < height {
		size = height
	}
	return
}

// Distance between two points in 2D space
func distance(p1, p2 pixel.Vec) float64 {
	return math.Sqrt(math.Pow(float64(p1.X-p2.X), 2) + math.Pow(float64(p1.Y-p2.Y), 2))
}

// Just print current memory usage
func printMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	Debug(fmt.Sprintf("Alloc = %v MiB", m.Alloc/1024/1024))
	Debug(fmt.Sprintf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024))
	Debug(fmt.Sprintf("\tSys = %v MiB", m.Sys/1024/1024))
	Debug(fmt.Sprintf("\tNumGC = %v\n", m.NumGC))
}

// Without this window is black, bug after mojave update for
// osx?
// Must move window a bit in order to make it draw the first
// time.
func centerWindow(win *pixelgl.Window) {
	x, y := pixelgl.PrimaryMonitor().Size()
	width, height := win.Bounds().Size().XY()
	win.SetPos(
		pixel.V(
			x/2-width/2,
			y/2-height/2,
		),
	)
}
