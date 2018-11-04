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
	CenterWindow(gWin)
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
		sheetFile:   "test.png",
		walkFrames:  []int{8, 9, 10, 11, 12, 13, 14},
		idleFrames:  []int{0, 2, 3, 4, 5, 6},
		shootFrames: []int{26},
		jumpFrames:  []int{15, 16, 17, 18, 19, 20},
		climbFrames: []int{1, 7},
		frameWidth:  12.0,
		life:        100.0,
		mobType:     entityPlayer,
	}
	test.create(100, 70)
	global.gController.setActiveEntity(&test)
	global.gCamera.setFollow(&test)

	for !global.gWin.Closed() && !global.gController.quit {
		dt := time.Since(last).Seconds()
		last = time.Now()

		global.gWin.Clear(global.gClearColor)

		global.gController.update(dt)
		global.gCamera.update(dt)
		global.gWorld.Draw(dt)
		global.gParticleEngine.update(dt)

		global.gWin.Update()
	}
}
