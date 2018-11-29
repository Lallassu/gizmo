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

// World constants
const (
	wShadowLength        = 5
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
	wViewMax             = 400
	wParticleDefaultLife = 5
	wGravity             = -9.82
	wParticlesMax        = 10000
	wAmmoMax             = 2000
	wDoorLen             = 30
	wDoorHeight          = 40
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
}

var global = &Global{
	gWindowHeight:   768,
	gWindowWidth:    1024,
	gVsync:          true,
	gUndecorated:    false,
	gWorld:          &world{},
	gCamera:         &camera{},
	gController:     &controller{},
	gParticleEngine: &particleEngine{},
	gClearColor:     pixel.RGBA{0.4, 0.4, 0.4, 1.0},
	gWin:            &pixelgl.Window{},
	gAmmoEngine:     &ammoEngine{},
}
