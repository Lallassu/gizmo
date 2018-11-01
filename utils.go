//=============================================================
// utils.go
//-------------------------------------------------------------
// Utility functions
//=============================================================
package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"runtime"
)

func distance(p1, p2 pixel.Vec) float64 {
	return math.Sqrt(math.Pow(float64(p1.X-p2.X), 2) + math.Pow(float64(p1.Y-p2.Y), 2))
}

func PrintMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	Debug(fmt.Sprintf("Alloc = %v MiB", m.Alloc/1024/1024))
	Debug(fmt.Sprintf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024))
	Debug(fmt.Sprintf("\tSys = %v MiB", m.Sys/1024/1024))
	Debug(fmt.Sprintf("\tNumGC = %v\n", m.NumGC))
}

// Without this window is black, bug after mojave update for osx?
func CenterWindow(win *pixelgl.Window) {
	x, y := pixelgl.PrimaryMonitor().Size()
	width, height := win.Bounds().Size().XY()
	win.SetPos(
		pixel.V(
			x/2-width/2,
			y/2-height/2,
		),
	)
}
