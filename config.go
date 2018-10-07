package main

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
)

//=============================================================
// Global variables
//=============================================================
var (
	gWindowHeight = 768
	gWindowWidth  = 1024
	gVsync        = true
	gUndecorated  = false
	gWorld        = &World{}
	gParticles    = &ParticleEngine{}
)
