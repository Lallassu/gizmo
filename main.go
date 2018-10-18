//=============================================================
// main.go
//-------------------------------------------------------------
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	_ "math"
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
		VSync:       global.gVsync,
		Undecorated: global.gUndecorated,
	}
	gWin, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	global.gWin = gWin

	PrintMemoryUsage()
	// Setup world etc.
	setup()
	PrintMemoryUsage()

	// Start game loop
	gameLoop()
}

//=============================================================
// Setup map, world, player etc.
//=============================================================
func setup() {
	// Camera setup
	global.gCamera.create()
	// Init map
	global.gWorld.Init()
	global.gWorld.NewMap(mapEasy)
	global.gCamera.setPosition(-200, -200)
	global.gCamera.zoom = 2

}

//=============================================================
// Game loop
//=============================================================
func gameLoop() {
	last := time.Now()
	//frames := 0
	//fps := time.Tick(time.Second / 60)
	//second := time.Tick(time.Second)

	//cs := pixel.Vec{0, 0}
	for !global.gWin.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		if global.gWin.Pressed(pixelgl.KeyDown) {
			global.gCamera.zoom -= 0.1
		}
		if global.gWin.Pressed(pixelgl.KeyUp) {
			global.gCamera.zoom += 0.1
		}
		if global.gWin.Pressed(pixelgl.KeyLeft) {
			global.gCamera.pos.X -= 5.1
		}
		if global.gWin.Pressed(pixelgl.KeyRight) {
			global.gCamera.pos.X += 5.1
		}
		if global.gWin.Pressed(pixelgl.KeyQ) {
			break
		}

		global.gWin.Clear(global.gClearColor)

		global.gCamera.update(dt)

		global.gWorld.Draw(dt)

		global.gWin.Update()
	}
}
