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
	cfg := pixelgl.WindowConfig{
		Title:       GameTitle,
		Bounds:      pixel.R(0, 0, 1024, 768),
		VSync:       global.gVsync,
		Undecorated: global.gUndecorated,
		//	Monitor:     pixelgl.PrimaryMonitor(), // Fullscreen
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
	global.gUI.create()
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

	// TEST
	createLights()

	// Full screen fragment shader
	var fragmentShader = `
             #version 330 core
             
             in vec2  vTexCoords;
             in vec4  vColor;
             
             out vec4 fragColor;
             
             uniform vec4 uTexBounds;
             uniform sampler2D uTexture;
             
             void main() {
				vec4 c = vColor;
				vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;
		  		vec4 tx = texture(uTexture, t);
				if (c.r == 1) {
			       fragColor = vec4(tx.r, tx.g, tx.b, tx.a);
			    } else {
			        if (c.a == 0.1111) {
				       c *= 2;
				       vec3 fc = vec3(1.0, 0.3, 0.1);
	                   vec2 borderSize = vec2(0.1); 

	                   vec2 rectangleSize = vec2(1.0) - borderSize; 

	                   float distanceField = length(max(abs(c.x)-rectangleSize,0.0) / borderSize);

	                   float alpha = 1.0 - distanceField;
			           fc *= abs(0.8 / (sin( c.x + sin(c.y)+ 1.3 ) * 5.0) );
                       fragColor = vec4(fc, alpha*5);
				    } else {
			           fragColor = vColor;
			        }
				}
             }
             
			 `
	//global.gWin.Canvas().SetUniform("utime", &global.utime)
	global.gWin.Canvas().SetFragmentShader(fragmentShader)

}

//=============================================================
// Game loop
//=============================================================
func gameLoop() {
	last := time.Now()
	frameDt := 0.0
	startTime := time.Now()

	//fps := time.Tick(time.Second / 1000)
	//second := time.Tick(time.Second)
	//frames := 0

	elapsed := 0.0

	for !global.gWin.Closed() && !global.gController.quit {
		dt := time.Since(last).Seconds()
		frameDt += dt
		last = time.Now()
		elapsed = float64(time.Now().Unix() - startTime.Unix())
		//	global.utime += float32(dt)

		for {
			if frameDt >= wMaxInvFPS {
				global.gWin.Clear(global.gClearColor)

				//	global.gWin.SetComposeMethod(pixel.ComposeOver)

				go global.gController.update(wMaxInvFPS)
				global.gWorld.Draw(wMaxInvFPS, elapsed)
				go global.gTextures.update(wMaxInvFPS)

				global.gParticleEngine.update(wMaxInvFPS)
				global.gAmmoEngine.update(wMaxInvFPS)

				//	global.gWin.SetComposeMethod(pixel.ComposePlus)
				go global.gCamera.update(wMaxInvFPS)

				//	global.gWin.SetColorMask(pixel.Alpha(0.4))
				//	global.gWin.SetComposeMethod(pixel.ComposePlus)
				//	global.gLights.update(wMaxInvFPS, elapsed)
				//	global.gWin.SetColorMask(pixel.Alpha(1))

				global.gUI.draw(wMaxInvFPS)

				// TEST
				drawLights(wMaxInvFPS)

				global.gWin.Update()
				//uPosX = float32(test.bounds.X)
				//uPosY = float32(test.bounds.Y)
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
