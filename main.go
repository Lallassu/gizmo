//=============================================================
// main.go
//-------------------------------------------------------------
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
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
	global.gController.create()
	global.gWorld.Init()
	global.gWorld.NewMap(mapEasy)
	global.gParticleEngine.create()
	global.gCamera.setPosition(0, 0)
	global.gCamera.zoom = 2
	global.gWin.SetSmooth(false)

}

//=============================================================
// Game loop
//=============================================================
func gameLoop() {
	last := time.Now()
	//frames := 0
	//fps := time.Tick(time.Second / 60)
	//second := time.Tick(time.Second)

	test := mob{
		sheetFile:   "/Users/nergal/Dropbox/dev/golang/src/GoD/test.png",
		walkFrames:  []int{0, 1},
		idleFrames:  []int{2, 3},
		shootFrames: []int{4, 5},
		jumpFrames:  []int{6, 7},
		climbFrames: []int{1, 2},
		frameWidth:  32.0,
		life:        100.0,
		pos:         pixel.Vec{X: 100, Y: 100},
	}
	test.create()

	for !global.gWin.Closed() && !global.gController.quit {
		dt := time.Since(last).Seconds()
		last = time.Now()

		global.gWin.Clear(global.gClearColor)

		global.gController.update(dt)
		global.gCamera.update(dt)
		global.gWorld.Draw(dt)
		global.gParticleEngine.update(dt)
		test.draw(dt)

		global.gWin.Update()
	}
}
