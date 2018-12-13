//=============================================================
// main.go
//-------------------------------------------------------------
//=============================================================
package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	_ "github.com/pkg/profile"
	"math/rand"
	"strings"
	"time"
)

//=============================================================
// Main
//=============================================================
func main() {
	//defer profile.Start().Stop()
	//defer profile.Start(profile.MemProfile).Stop()
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
	global.gRand.create(100000)
	global.gCamera.create()
	global.gController.create()
	global.gWorld.Init()
	global.gWorld.NewMap(mapEasy)
	global.gParticleEngine.create()
	global.gAmmoEngine.create()
	global.gCamera.setPosition(0, 0)
	global.gCamera.zoom = 3
	global.gWin.SetSmooth(false)
	global.gPlayer.create(100, 50)
	global.gController.setActiveEntity(global.gPlayer)
	global.gCamera.setFollow(global.gPlayer)
	global.gTextures.load("packed.json")
}

//=============================================================
// Game loop
//=============================================================
func gameLoop() {
	last := time.Now()
	frameDt := 0.0
	startTime := time.Now()

	// var fragmentShader = `
	//   #version 330 core
	//
	//   in vec2  vTexCoords;
	//   in vec2 vPosition;
	//   in vec4 vColor;
	//
	//   out vec4 fragColor;
	//
	//   uniform float uPosX;
	//   uniform float uPosY;
	//   uniform vec4 uTexBounds;
	//   uniform sampler2D uTexture;
	//
	//   void main() {
	//   	// Get our current screen coordinate
	//   	vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;
	//
	//      vec4 color = vec4(0.0,0.0,0.0,1.0);
	// 	 if (vColor.r == 1.0 && vColor.g == 1.0 && vColor.b == 1.0 && vColor.a == 1.0 ){
	//         color = vec4(texture(uTexture,t).r, texture(uTexture,t).g, texture(uTexture,t).b, texture(uTexture,t).a);
	//     } else {
	//          int val = int(vColor.a*255) & 0xFF;
	//
	//          float distance = sqrt(pow(vPosition.x-uPosX+cos(uPosX/10)*10, 2) + pow(vPosition.y-(uPosY+10)+cos(uPosY/10)*10, 2))/5;

	//          if (val == 0x8F) {
	//        		    color = vec4(vColor.r/(distance/2), vColor.g/distance, vColor.b/distance, vColor.a);
	//          } else {
	//        		 color = vec4(vColor.r, vColor.g, vColor.b, vColor.a);
	//     	 }
	//     }
	//   	fragColor = color;
	//   }
	//   `
	//fps := time.Tick(time.Second / 1000)
	//second := time.Tick(time.Second)
	//frames := 0

	// Load a bunch of weapons
	for _, x := range []string{"ak47_weapon", "p90_weapon", "rocketlauncher_weapon"} { // "shotgun_weapon", "crate_obj"} {
		var otype objectType
		scale := 0.15

		static := false
		if strings.Contains(x, "_weapon") {
			otype = objectWeapon
			static = false
		} else if strings.Contains(x, "_obj") {
			otype = objectCrate
			scale = 1
		}

		for i := 0; i < 15; i++ {
			objTest := object{
				textureFile: fmt.Sprintf("assets/objects/%v.png", x),
				name:        x,
				static:      static,
				entityType:  entityObject,
				objectType:  otype,
				scale:       scale,
			}
			objTest.create(float64(rand.Intn(global.gWorld.width)), float64(rand.Intn(global.gWorld.height)))
		}
	}

	for i := 0; i < 10; i++ {
		test := mob{
			sheetFile:   "assets/mobs/enemy1.png",
			walkFrames:  []int{8, 9, 10, 11, 12, 13, 14},
			idleFrames:  []int{0, 2, 3, 4, 5, 6},
			shootFrames: []int{26},
			jumpFrames:  []int{15, 16, 17, 18, 19, 20},
			climbFrames: []int{1, 7},
			frameWidth:  12.0,
			life:        100.0,
			speed:       100,
			mobType:     entityEnemy,
			ai:          &AI{},
		}
		test.create(float64(rand.Intn(global.gWorld.width)), float64(rand.Intn(global.gWorld.height)))
	}

	// var uPosX, uPosY float32
	// global.gWin.Canvas().SetUniform("uPosX", &uPosX)
	// global.gWin.Canvas().SetUniform("uPosY", &uPosY)
	// global.gWin.Canvas().SetFragmentShader(fragmentShader)
	elapsed := 0.0

	for !global.gWin.Closed() && !global.gController.quit {
		dt := time.Since(last).Seconds()
		frameDt += dt
		last = time.Now()
		elapsed = float64(time.Now().Unix() - startTime.Unix())

		for {
			if frameDt >= wMaxInvFPS {
				global.gWin.Clear(global.gClearColor)
				global.gController.update(wMaxInvFPS)
				global.gWorld.Draw(wMaxInvFPS, elapsed)
				global.gTextures.update(wMaxInvFPS)
				global.gParticleEngine.update(wMaxInvFPS)
				global.gAmmoEngine.update(wMaxInvFPS)
				global.gCamera.update(wMaxInvFPS)
				global.gWin.Update()
				//uPosX = float32(test.bounds.X)
				//uPosY = float32(test.bounds.Y)
			} else {
				break
			}

			frameDt -= wMaxInvFPS
		}

		//  <-fps
		//  updateFPSDisplay(global.gWin, &frames, second)
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
