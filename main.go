//=============================================================
// main.go
//-------------------------------------------------------------
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	_ "github.com/pkg/profile"
	"time"
)

//=============================================================
// Main
//=============================================================
func main() {
	//	defer profile.Start(profile.CPUProfile).Stop()
	pixelgl.Run(run)
}

//=============================================================
// Setup game window etc.
//=============================================================
func run() {
	global.gVariableConfig.LoadConfiguration()

	cfg := pixelgl.WindowConfig{
		Title:       GameTitle,
		Bounds:      pixel.R(0, 0, float64(global.gVariableConfig.WindowWidth), float64(global.gVariableConfig.WindowHeight)),
		VSync:       global.gVariableConfig.Vsync,
		Undecorated: global.gVariableConfig.UndecoratedWindow,
	}
	if global.gVariableConfig.Fullscreen {
		cfg.Monitor = pixelgl.PrimaryMonitor()
	}
	gWin, err := pixelgl.NewWindow(cfg)
	//gWin.SetBounds(pixel.R(0, 0, 800, 600))

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
	global.gFont.create()
	global.gMainMenu.createMain()
	global.gOptionsMenu.createOptions()
	global.gControllerMenu.createController()
	global.gDisplayMenu.createDisplay()
	global.gGameMenu.createGame()

	// Initialize with menu
	global.gActiveMenu = global.gMainMenu

	global.gUI.create()
	global.gMapColor.create()
	global.gRand.create(100000)
	global.gSounds.create()
	global.gCamera.create()
	global.gController.create()
	global.gWorld.Init()
	global.gParticleEngine.create()
	global.gAmmoEngine.create()
	global.gCamera.setPosition(0, 0)
	global.gCamera.zoom = 3
	global.gWin.SetSmooth(false)
	global.gController.setActiveEntity(global.gPlayer)
	global.gCamera.setFollow(global.gPlayer)
	global.gTextures.load("packed.json")
	global.gMap.newMap(1)

	//global.gWin.Canvas().SetUniform("utime", &global.utime)
	global.gWin.Canvas().SetFragmentShader(fragmentShaderFullScreen)

}

//=============================================================
// Game loop
//=============================================================
func gameLoop() {
	last := time.Now()
	frameDt := 0.0

	//fps := time.Tick(time.Second / 1000)
	//second := time.Tick(time.Second)
	//frames := 0

	elapsed := 0.0

	for !global.gWin.Closed() && !global.gController.quit {
		dt := time.Since(last).Seconds()
		frameDt += dt
		last = time.Now()

		for {
			if frameDt >= wMaxInvFPS {
				elapsed += wMaxInvFPS
				global.uTime += wMaxInvFPS

				global.gWin.Clear(global.gClearColor)

				global.gController.update(wMaxInvFPS)
				global.gWorld.Draw(wMaxInvFPS, elapsed)
				global.gTextures.update(wMaxInvFPS)

				global.gParticleEngine.update(wMaxInvFPS)
				global.gAmmoEngine.update(wMaxInvFPS)

				global.gCamera.update(wMaxInvFPS)

				if global.gActiveMenu != nil {
					global.gActiveMenu.draw(wMaxInvFPS, elapsed)
				} else {
					global.gUI.draw(wMaxInvFPS)
				}

				global.gWin.Update()
			} else {
				break
			}
			frameDt -= wMaxInvFPS
			//  <-fps
			//  updateFPSDisplay(global.gWin, &frames, second)
		}

	}
}

func updateFPSDisplay(win *pixelgl.Window, frames *int, second <-chan time.Time) {
	*frames++
	select {
	case <-second:
		//	win.SetTitle(fmt.Sprintf("%s (FPS: %d)", GameTitle, *frames))
		global.gUI.updateFPS(*frames)
		*frames = 0
	default:
	}
}
