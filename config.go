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
	wPixelsPerChunk     = 64
	wPixelSize          = 1
	wBorderSize         = 5
	wStaticBorderSize   = 0
	wStaticColor32      = 0xFFFFFFFE
	wStaticColor8       = 0xFE
	wFloodFill8         = 0xFC
	wFloodFill32        = 0xFFFFFFFC
	wFloodFillVisited8  = 0xFB
	wFloodFillVisited32 = 0xFFFFFFFB
	wBackground32       = 0xFFFFFF8F
	wBackground8        = 0x8F
	wLadder8            = 0xAF
	wLadder32           = 0xFFFFFFAF
	wShadow8            = 0xBF
	wShadow32           = 0xFFFFFFBF
	wViewMax            = 400
)

//=============================================================
// Global variables
//=============================================================
type Global struct {
	gWindowHeight int
	gWindowWidth  int
	gVsync        bool
	gUndecorated  bool
	gWorld        *world
	gCamera       *camera
	gParticles    *particleEngine
	gClearColor   pixel.RGBA
	gWin          *pixelgl.Window
}

var global = &Global{
	gWindowHeight: 768,
	gWindowWidth:  1024,
	gVsync:        true,
	gUndecorated:  false,
	gWorld:        &world{},
	gCamera:       &camera{},
	gParticles:    &particleEngine{},
	gClearColor:   pixel.RGBA{0.4, 0.4, 0.4, 1.0},
	gWin:          &pixelgl.Window{},
}
