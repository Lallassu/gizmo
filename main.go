//=============================================================
// main.go
//-------------------------------------------------------------
//=============================================================
package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"time"
)

//=============================================================
// Main
//=============================================================
func main() {
	pixelgl.Run(run)
}

//=============================================================
// Setup game window etc.
//=============================================================
func run() {
	cfg := pixelgl.WindowConfig{
		Title:       GameTitle,
		Bounds:      pixel.R(0, 0, 1024, 768),
		VSync:       Vsync,
		Undecorated: Undecorated,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Setup world etc.
	setup()

	// Start game loop
	gameLoop(win)
}

//=============================================================
// Setup map, world, player etc.
//=============================================================
func setup() {
	// Init map

	// Init

	//
}

//=============================================================
// Game loop
//=============================================================
func gameloop(win *pixelgl.Window) {
	last := time.Now()
	//frames := 0
	fps := time.Tick(time.Second / 60)
	second := time.Tick(time.Second)

	for !win.Closed() {
		win.Update()
		<-fps
		//updateFPSDisplay(win, &frames, second)
	}
}
