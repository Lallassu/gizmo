//=============================================================
// main.go
//-------------------------------------------------------------
//=============================================================
package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/profile"
	"math/rand"
	"time"
)

//=============================================================
// Main
//=============================================================
func main() {
	//defer profile.Start().Stop()
	defer profile.Start(profile.MemProfile).Stop()
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
		//	Monitor:     pixelgl.PrimaryMonitor(), // Fullscreen
	}
	gWin, err := pixelgl.NewWindow(cfg)

	if err != nil {
		panic(err)
	}
	CenterWindow(gWin)
	global.gWin = gWin

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
	global.gAmmoEngine.create()
	global.gCamera.setPosition(0, 0)
	global.gCamera.zoom = 2
	global.gWin.SetSmooth(false)

}

//=============================================================
// Game loop
//=============================================================
func gameLoop() {
	last := time.Now()
	fps := time.Tick(time.Second / 1000)
	second := time.Tick(time.Second)
	frames := 0

	// Load a bunch of weapons
	for _, x := range []string{"ak47.png", "p90.png", "rocketlauncher.png", "shotgun.png"} {
		objTest := object{
			textureFile: fmt.Sprintf("assets/objects/%v", x),
			entityType:  entityObject,
			scale:       0.15,
		}
		objTest.create(100+float64(rand.Intn(50)), 180+float64(rand.Intn(50)))
	}

	// Load a player
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
	test.create(100, 120)
	global.gController.setActiveEntity(&test)
	global.gCamera.setFollow(&test)

	for !global.gWin.Closed() && !global.gController.quit {
		dt := time.Since(last).Seconds()
		last = time.Now()

		global.gWin.Clear(global.gClearColor)

		global.gController.update(dt)
		global.gWorld.Draw(dt)
		global.gParticleEngine.update(dt)
		global.gAmmoEngine.update(dt)
		global.gCamera.update(dt)

		global.gWin.Update()

		<-fps
		updateFPSDisplay(global.gWin, &frames, second)
	}
}

func updateFPSDisplay(win *pixelgl.Window, frames *int, second <-chan time.Time) {
	*frames++
	select {
	case <-second:
		win.SetTitle(fmt.Sprintf("%s (FPS: %d)", GameTitle, *frames))
		*frames = 0
	default:
	}
}
