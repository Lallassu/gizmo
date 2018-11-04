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
	"math"
	"math/rand"
	"os"
)

type mob struct {
	sheetFile   string
	life        float64
	walkFrames  []int
	jumpFrames  []int
	climbFrames []int
	shootFrames []int
	idleFrames  []int
	frameWidth  int
	frameHeight int
	currentAnim animationType
	canvas      map[int]*pixelgl.Canvas
	frames      map[int][]uint32
	bounds      *Bounds
	mobType     entityType
	animCounter float64
	animRate    float64
	speed       float64
	dir         float64
}

//=============================================================
// Create mob
// - load animation sheet
//=============================================================
func (m *mob) create(x, y float64) {
	m.canvas = make(map[int]*pixelgl.Canvas)
	m.frames = make(map[int][]uint32)
	m.currentAnim = animIdle

	m.animRate = 0.1
	m.speed = 10

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
	m.bounds = &Bounds{
		X:      x,
		Y:      y,
		Width:  float64(m.frameWidth),
		Height: float64(m.frameHeight),
		entity: Entity(m),
	}

	f := 0
	size := m.frameWidth
	if m.frameWidth < m.frameHeight {
		size = m.frameHeight
	}
	for w := 0; w < imgCfg.Width; w += m.frameWidth {
		m.frames[f] = make([]uint32, size*size)
		for x := 0; x < m.frameWidth; x++ {
			for y := 0; y < imgCfg.Height; y++ {
				r, g, b, a := img.At(w+x, imgCfg.Height-y).RGBA()
				m.frames[f][x*m.frameWidth+y] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | a&0xFF
			}
		}
		m.canvas[f] = pixelgl.NewCanvas(pixel.R(0, 0, float64(m.frameWidth), float64(m.frameHeight)))
		f++
	}

	// Build all frames (canvases)
	m.buildFrames()

	// Add object to QT
	global.gWorld.AddObject(m.bounds)
}

//=============================================================
// Build each frame into a canvas
//=============================================================
func (m *mob) buildFrames() {
	model := imdraw.New(nil)
	for i := 0; i < len(m.frames); i++ {
		for x := 0; x < m.frameWidth; x++ {
			for y := 0; y < m.frameHeight; y++ {
				p := m.frames[i][m.frameWidth*x+y]
				if p == 0 {
					continue
				}

				model.Color = pixel.RGB(
					float64(p>>24&0xFF)/255.0,
					float64(p>>16&0xFF)/255.0,
					float64(p>>8&0xFF)/255.0,
				).Mul(pixel.Alpha(float64(p&0xFF) / 255.0))
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
func (m *mob) hit(x_, y_ float64) bool {
	x := int(math.Abs(float64(m.bounds.X - x_)))
	y := int(math.Abs(float64(m.bounds.Y - y_)))
	for i := 0; i < len(m.frames); i++ {
		pos := m.frameWidth*x + y
		if pos >= 0 && pos < m.frameWidth*m.frameWidth {
			if m.frames[i][pos] != 0 {
				// Make some blood (TBD: ADD TO PARTICLE ENGINE)
				for i := 0; i < 10; i++ {
					r := 175 + rand.Intn(50)
					g := 10 + rand.Intn(20)
					b := 10 + rand.Intn(20)
					a := 20 + rand.Intn(220)
					m.frames[i][pos] = 0
					global.gParticleEngine.newParticle(
						particle{
							x:           float64(x_),
							y:           float64(y_),
							size:        rand.Float64() * 3,
							restitution: -0.1 - rand.Float64()/4,
							life:        wParticleDefaultLife,
							fx:          rand.Float64() * 5,
							fy:          rand.Float64() * 5,
							vx:          float64(5 - rand.Intn(10)),
							vy:          float64(5 - rand.Intn(10)),
							mass:        2,
							pType:       particleRegular,
							color:       uint32((r & 0xFF << 24) | (g & 0xFF << 16) | (b & 0xFF << 8) | (a & 0xFF)),
							static:      true,
						})
				}
				// Remove part
				m.frames[i][pos] = 0
				global.gParticleEngine.newParticle(
					particle{
						x:           float64(x_),
						y:           float64(y_),
						size:        1,
						restitution: -0.1 - rand.Float64()/4,
						life:        wParticleDefaultLife,
						fx:          10,
						fy:          10,
						vx:          float64(5 - rand.Intn(10)),
						vy:          float64(5 - rand.Intn(10)),
						mass:        1,
						pType:       particleRegular,
						color:       m.frames[i][pos],
						static:      true,
					})
			}
		}
	}
	m.buildFrames()
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
func (m *mob) move(x, y float64) {
	if math.Abs(x) > 0 {
		m.currentAnim = animWalk
	} else {
		m.currentAnim = animIdle
	}
	if math.Abs(y) > 0 {
		m.currentAnim = animClimb
	}

	if x > 0 {
		m.dir = 1
	} else {
		m.dir = -1
	}
	m.bounds.X += x
	m.bounds.Y += y
}

//=============================================================
//
//=============================================================
func (m *mob) getPosition() pixel.Vec {
	return pixel.Vec{m.bounds.X, m.bounds.Y}
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

func (m *mob) setPos(x, y float64) {
	m.bounds.X = x
	m.bounds.Y = y
}

//=============================================================
//
//=============================================================
func (m *mob) draw(dt float64) {
	m.animCounter += dt
	idx := int(math.Floor(m.animCounter / m.animRate))

	switch m.currentAnim {
	case animWalk:
		idx = m.walkFrames[idx%len(m.walkFrames)]
	case animJump:
		idx = m.jumpFrames[idx%len(m.jumpFrames)]
	case animClimb:
		idx = m.climbFrames[idx%len(m.climbFrames)]
	case animShoot:
		idx = m.shootFrames[idx%len(m.shootFrames)]
	case animIdle:
		idx = m.idleFrames[idx%len(m.idleFrames)]
	default:
		idx = m.idleFrames[idx%len(m.idleFrames)]
	}
	m.canvas[idx].Draw(global.gWin, pixel.IM.ScaledXY(pixel.ZV, pixel.V(-m.dir, 1)).Moved(pixel.V(m.bounds.X+m.bounds.Width/2, m.bounds.Y+m.bounds.Height/2)))

}
