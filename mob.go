//=============================================================
// mob.go
//-------------------------------------------------------------
// Anything that can move/be destroyed etc.
//=============================================================
package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"image"
	"image/png"
	"os"
)

type mob struct {
	sheetFile   string
	life        float64
	pos         pixel.Vec
	walkFrames  []int
	jumpFrames  []int
	climbFrames []int
	shootFrames []int
	idleFrames  []int
	frameWidth  int
	frameHeight int
	currentAnim animationType
	canvas      []*pixelgl.Canvas
	frames      map[int][]uint32
	bounds      Bounds
}

//=============================================================
// Create mob
// - load animation sheet
//=============================================================
func (m *mob) create() {
	m.canvas = make([]*pixelgl.Canvas, 0)
	m.frames = make(map[int][]uint32)

	// Load animation
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	imgfile, err := os.Open(m.sheetFile)
	if err != nil {
		Error(fmt.Sprintf("Failed to open sheet file %v", m.sheetFile))
		return
	}

	defer imgfile.Close()

	imgCfg, _, err := image.DecodeConfig(imgfile)
	if err != nil {
		Error(fmt.Sprintf("Failed to decode sheet file %v: %v", m.sheetFile, err))
		return
	}

	imgfile.Seek(0, 0)
	img, _, _ := image.Decode(imgfile)

	m.frameHeight = imgCfg.Height

	// Initiate bounds for qt
	m.bounds = Bounds{
		X:      m.pos.X,
		Y:      m.pos.Y,
		Width:  float64(m.frameWidth),
		Height: float64(m.frameHeight),
		entity: Entity(m),
	}

	f := 0
	for w := 0; w < imgCfg.Width; w += m.frameWidth {
		f++
		for x := w; x < m.frameWidth; x++ {
			for y := 0; y < imgCfg.Height; y++ {
				r, g, b, a := img.At(x, imgCfg.Height-y-1).RGBA()
				m.frames[f] = []uint32{r, g, b, a}
			}
		}
	}

	m.buildFrames()
}

//=============================================================
// Build each frame into a canvas
//=============================================================
func (m *mob) buildFrames() {
	model := imdraw.New(nil)
	for i := 0; i < len(m.frames); i++ {
		for x := 0; x < m.frameWidth; x++ {
			for y := 0; y < m.frameHeight; y++ {
				r := m.frames[i][0]
				g := m.frames[i][1]
				b := m.frames[i][2]
				a := m.frames[i][3]
				model.Color = pixel.RGB(
					float64(r), float64(g), float64(b),
				).Mul(pixel.Alpha(float64(a)))

				model.Push(
					pixel.V(float64(x*wPixelSize), float64(y*wPixelSize)),
					pixel.V(float64(x*wPixelSize+wPixelSize), float64(y*wPixelSize+wPixelSize)),
				)
				model.Rectangle(0)
			}
		}
		m.canvas[i].Clear(pixel.RGBA{0, 0, 0, 0})
		model.Draw(m.canvas[i])
	}
}

//=============================================================
//
//  Function to implement Entity interface
//
//=============================================================
//=============================================================
//
//=============================================================
func (m *mob) hit(x, y float64) bool {
	return true
}

//=============================================================
//
//=============================================================
func (m *mob) explode() {
}

//=============================================================
//
//=============================================================
func (m *mob) getMass() float64 {
	return 1.0
}

//=============================================================
//
//=============================================================
func (m *mob) getType() entityType {
	return entityEnemy
}

//=============================================================
//
//=============================================================
func (m *mob) draw(dt float64) {
	idx := 0
	switch m.currentAnim {
	case animWalk:
	case animJump:
	case animClimb:
	case animShoot:
	default:
		// Idle
	}

	m.canvas[idx].Draw(global.gWin, pixel.IM.Moved(pixel.V(m.bounds.X+m.bounds.Width/2, m.bounds.Y+m.bounds.Height/2)))

}
