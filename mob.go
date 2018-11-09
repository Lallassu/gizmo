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
	"time"
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
	size        int
	currentAnim animationType
	models      map[int]*imdraw.IMDraw
	canvas      *pixelgl.Canvas
	frames      map[int][]uint32
	bounds      *Bounds
	mobType     entityType
	animCounter float64
	animRate    float64
	speed       float64
	dir         float64
	mass        float64
	cdPixels    [][2]uint32

	prevPos     []pixel.Vec
	force       pixel.Vec
	moveVec     pixel.Vec
	restitution float64
	climbing    bool
	jumping     float64
	jumpPower   float64
}

//=============================================================
// Create mob
// - load animation sheet
//=============================================================
func (m *mob) create(x, y float64) {
	m.models = make(map[int]*imdraw.IMDraw)
	m.frames = make(map[int][]uint32)
	m.cdPixels = make([][2]uint32, 10)
	m.prevPos = make([]pixel.Vec, 100)
	m.currentAnim = animIdle

	m.animRate = 0.1
	m.jumpPower = 35.0
	m.speed = 200
	m.mass = 50

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
	m.size = m.frameWidth
	if m.frameWidth < m.frameHeight {
		m.size = m.frameHeight
	}
	for w := 0; w < imgCfg.Width; w += m.frameWidth {
		m.frames[f] = make([]uint32, m.size*m.size)
		for x := 0; x <= m.frameWidth; x++ {
			for y := 0; y <= imgCfg.Height; y++ {
				r, g, b, a := img.At(w+x, imgCfg.Height-y).RGBA()
				m.frames[f][x*m.size+y] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | a&0xFF
			}
		}
		f++
	}

	// Generate some CD pixel for faster CD check.
	rand.Seed(time.Now().UTC().UnixNano())
	for x := 0; x < 20; x++ {
		m.cdPixels = append(m.cdPixels, [2]uint32{uint32(rand.Intn(m.frameWidth)), uint32(rand.Intn(m.frameHeight))})
	}

	m.canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(m.frameWidth), float64(m.frameHeight)))

	// Build all frames
	m.buildFrames()

	// Add object to QT
	global.gWorld.AddObject(m.bounds)
}

//=============================================================
// Build each frame into a canvas
//=============================================================
func (m *mob) buildFrames() {
	for i := 0; i < len(m.frames); i++ {
		m.models[i] = imdraw.New(nil)
		for x := 0; x < m.frameWidth; x++ {
			for y := 0; y < m.frameHeight; y++ {
				p := m.frames[i][x*m.size+y]
				if p == 0 {
					continue
				}

				m.models[i].Color = pixel.RGB(
					float64(p>>24&0xFF)/255.0,
					float64(p>>16&0xFF)/255.0,
					float64(p>>8&0xFF)/255.0,
				).Mul(pixel.Alpha(float64(p&0xFF) / 255.0))
				m.models[i].Push(
					pixel.V(float64(x*wPixelSize), float64(y*wPixelSize)),
					pixel.V(float64(x*wPixelSize+wPixelSize), float64(y*wPixelSize+wPixelSize)),
				)
				m.models[i].Rectangle(0)
			}
		}
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
		pos := m.size*x + y
		if pos >= 0 && pos < m.frameWidth*m.frameWidth {
			if m.frames[i][pos] != 0 {
				// Add some blood
				global.gParticleEngine.effectBlood(x_, y_, 1)
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
// Check if on ground
//=============================================================
func (m *mob) IsOnGround() bool {
	for x := m.bounds.X; x < m.bounds.X+m.bounds.Width; x += 2 {
		if global.gWorld.IsRegular(x, m.bounds.Y) {
			return true
		}
	}
	return false
}

//=============================================================
// Check if on ladder
//=============================================================
func (m *mob) IsOnLadder() bool {
	for _, p := range m.cdPixels {
		if global.gWorld.IsLadder(m.bounds.X+float64(p[0]), m.bounds.Y+float64(p[1])) {
			return true
		}
	}
	return false
}

//=============================================================
// Check if on wall
//=============================================================
func (m *mob) IsOnWall() bool {
	offset := 2.0
	for _, p := range m.cdPixels {
		if global.gWorld.IsRegular(m.bounds.X+float64(p[0])+m.force.X, m.bounds.Y+float64(p[1])+m.force.Y+offset) {
			return true
		}
	}
	return false
}

//=============================================================
//
//=============================================================
func (m *mob) move(dx, dy float64) {
	// Add the force, movenment is handled in the physics function
	m.force.X += dx * m.speed
	m.force.Y += dy * m.speed

	// Rotation of animation
	if dx != 0 {
		if dx > 0 {
			m.dir = 1
		} else {
			m.dir = -1
		}
	}
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
	return m.mass
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
func (m *mob) setPosition(x, y float64) {
	m.bounds.X = x
	m.bounds.Y = y
}

//=============================================================
//=============================================================
func (m *mob) saveMove() {
	m.prevPos = append(m.prevPos, pixel.Vec{m.bounds.X, m.bounds.Y})
	// TBD: Only remove every second or something
	if len(m.prevPos) > 100 {
		m.prevPos = m.prevPos[:100]
	}
}

//=============================================================
// Physics for mob.
// I don't want real physics, better to have a good feeling for
// movement than accurate physic simulation.
//=============================================================
func (m *mob) physics(dt float64) {
	m.saveMove()

	// Only move if no wall collision
	if !m.IsOnWall() {
		m.bounds.X += m.force.X
		m.currentAnim = animWalk
		if m.force.X == 0 {
			m.currentAnim = animIdle
		}
	}

	if m.jumping > 0 {
		// First half, up, then decrease force.
		if m.jumping > m.jumpPower/2 {
			m.force.Y = m.speed * dt
		} else {
			m.force.Y = -m.speed * dt
		}
		if !m.IsOnWall() {
			m.bounds.Y += m.force.Y
			m.bounds.X += m.force.X / 2
			m.currentAnim = animJump
		}
		m.jumping--
	} else {
		if m.IsOnLadder() && m.force.Y > 0 && m.force.X != 0 && m.jumping <= 0 {
			// Jump from ladder
			m.jumping = m.jumpPower
		} else if m.IsOnLadder() {
			// Climb
			m.bounds.Y += m.force.Y / 2
			m.currentAnim = animClimb
		} else if m.IsOnGround() {
			// Jump
			if m.force.Y > 0 && m.jumping <= 0 && !m.IsOnLadder() {
				m.jumping = m.jumpPower
			}
		} else {
			m.bounds.Y += global.gWorld.gravity * dt * m.mass
			m.currentAnim = animJump
		}
	}

	// Check if stuck!
	//  for m.IsStuck() {

	//  }

	m.force.X = 0
	m.force.Y = 0

}

//=============================================================
//
//=============================================================
func (m *mob) draw(dt float64) {
	// Update physics
	m.physics(dt)

	// Update animation
	m.animCounter += dt
	idx := int(math.Floor(m.animCounter / m.animRate))

	switch m.currentAnim {
	case animWalk:
		idx = m.walkFrames[idx%len(m.walkFrames)]
	case animJump:
		idx = m.jumpFrames[idx%len(m.jumpFrames)]
	case animClimb:
		idx = m.climbFrames[idx/2%len(m.climbFrames)]
	case animShoot:
		idx = m.shootFrames[idx%len(m.shootFrames)]
	case animIdle:
		idx = m.idleFrames[idx/30%len(m.idleFrames)]
	default:
		idx = m.idleFrames[idx/30%len(m.idleFrames)]
	}

	// reset animation
	if m.currentAnim != animClimb && m.currentAnim != animJump {
		m.currentAnim = animIdle
	}

	m.canvas.Clear(pixel.RGBA{0, 0, 0, 0})
	m.models[idx].Draw(m.canvas)
	m.canvas.Draw(global.gWin, pixel.IM.ScaledXY(pixel.ZV, pixel.V(-m.dir, 1)).Moved(pixel.V(m.bounds.X+m.bounds.Width/2, m.bounds.Y+m.bounds.Height/2)))
}
