package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

//=============================================================
// Game constants
//=============================================================
// General constants
const (
	GameTitle   = "GoD - Game or Die"
	GameVersion = "0.0"
)

var gSoundFiles = []string{
	"jump.wav",
	"shot.mp3",
}

// World constants
const (
	wMaxInvFPS           = 1 / 60.0
	wShadowLength        = 5
	wShadowDepth         = 1.5
	wPixelsPerChunk      = 64
	wPixelSize           = 1
	wBorderSize          = 4
	wStaticBorderSize    = 0
	wStaticColor32       = 0xFFFFFFFE
	wStaticColor8        = 0xFE
	wFloodFill8          = 0xFC
	wFloodFill32         = 0xFFFFFFFC
	wFloodFillVisited8   = 0xFB
	wFloodFillVisited32  = 0xFFFFFFFB
	wBackground32        = 0xFFFFFF8F
	wBackground8         = 0x8F
	wLadder8             = 0xAF
	wLadder32            = 0xFFFFFFAF
	wShadow8             = 0xBF
	wShadow32            = 0xFFFFFFBF
	wViewMax             = 450
	wParticleDefaultLife = 5
	wGravity             = -9.82
	wParticlesMax        = 10000
	wAmmoMax             = 2000
	wDoorLen             = 30
	wDoorHeight          = 40
	wLightsMax           = 10
)

//=============================================================
// Global variables
//=============================================================
type Global struct {
	gWindowHeight   int
	gWindowWidth    int
	gVsync          bool
	gUndecorated    bool
	gWorld          *world
	gCamera         *camera
	gParticleEngine *particleEngine
	gClearColor     pixel.RGBA
	gWin            *pixelgl.Window
	gController     *controller
	gAmmoEngine     *ammoEngine
	gTextures       *textures
	gRand           *fRand
	gPlayer         *mob
	gLights         *lights
	gSounds         *sound
	gUI             *UI
}

var global = &Global{
	gWindowHeight:   768,
	gWindowWidth:    1024,
	gVsync:          false,
	gUndecorated:    false,
	gWorld:          &world{},
	gCamera:         &camera{},
	gController:     &controller{},
	gParticleEngine: &particleEngine{},
	gClearColor:     pixel.RGBA{0, 0, 0, 1.0},
	gWin:            &pixelgl.Window{},
	gAmmoEngine:     &ammoEngine{},
	gTextures:       &textures{},
	gRand:           &fRand{},
	gLights:         &lights{},
	gSounds:         &sound{},
	gUI:             &UI{},
	gPlayer: &mob{
		sheetFile:   "assets/mobs/player.png",
		walkFrames:  []int{8, 9, 10, 11, 12, 13, 14},
		idleFrames:  []int{0, 2, 3, 4, 5, 6},
		shootFrames: []int{26},
		jumpFrames:  []int{15, 16, 17, 18, 19, 20},
		climbFrames: []int{1, 7},
		frameWidth:  12.0,
		life:        100.0,
		mobType:     entityPlayer,
	},
}
