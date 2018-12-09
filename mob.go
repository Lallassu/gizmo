//=============================================================
// mob.go
//-------------------------------------------------------------
// Anything that can move/be destroyed etc.
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"image"
	"math"
	"math/rand"
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
	img         image.Image
	carry       *object

	prevPos      []pixel.Vec
	force        pixel.Vec
	velo         pixel.Vec
	restitution  float64
	climbing     bool
	jumping      bool
	jumpPower    float64
	acceleration float64
	keyMove      pixel.Vec
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

	m.animRate = 0.1
	m.jumpPower = 4
	m.speed = 200
	m.mass = 20
	m.currentAnim = animIdle
	m.dir = 1

	fullWidth := 0

	m.img, fullWidth, m.frameHeight, m.size = loadTexture(m.sheetFile)

	// Initiate bounds for qt
	m.bounds = &Bounds{
		X:      x,
		Y:      y,
		Width:  float64(m.frameWidth),
		Height: float64(m.frameHeight),
		entity: Entity(m),
	}

	f := 0
	for w := 0; w < fullWidth; w += m.frameWidth {
		m.frames[f] = make([]uint32, m.size*m.size)
		for x := 0; x <= m.frameWidth; x++ {
			for y := 0; y <= m.frameHeight; y++ {
				r, g, b, a := m.img.At(w+x, m.frameHeight-y).RGBA()
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
// Build each frame
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
//=============================================================
func (m *mob) hit(x_, y_, vx, vy float64, power int) bool {
	global.gParticleEngine.effectBlood(x_, y_, vx, vy, 1)

	x := int(math.Abs(float64(m.bounds.X - x_)))
	y := int(math.Abs(float64(m.bounds.Y - y_)))

	pow := power * power
	for rx := x - power; rx <= x+power; rx++ {
		xx := (rx - x) * (rx - x)
		for ry := y - power; ry <= y+power; ry++ {
			if ry < 0 {
				continue
			}
			val := (ry-y)*(ry-y) + xx
			if val < pow {
				for i := 0; i < len(m.frames); i++ {
					pos := m.size*rx + ry
					if pos >= 0 && pos < m.size*m.size {
						if m.frames[i][pos] != 0 {
							m.frames[i][pos] = 0
							global.gParticleEngine.newParticle(
								particle{
									x:           float64(x_),
									y:           float64(y_),
									size:        1,
									restitution: -0.1 - global.gRand.randFloat()/4,
									life:        wParticleDefaultLife,
									fx:          10,
									fy:          10,
									vx:          vx, //float64(5 - rand.Intn(10)),
									vy:          float64(5 - global.gRand.rand()),
									mass:        1,
									pType:       particleRegular,
									color:       m.frames[i][pos],
									static:      true,
								})
						}
					}
				}
			}
		}
	}

	m.buildFrames()

	// Check status of mob, if dead => remove from QT
	global.gWorld.qt.Remove(m.bounds)

	return true
}

//=============================================================
// Shoot if weapon attached
//=============================================================
func (m *mob) shoot() {
	if m.carry != nil {
		m.carry.shoot()
		m.currentAnim = animShoot
	}
}

//=============================================================
// Attach object to self
//=============================================================
func (m *mob) attach(o *object) {
	if m.carry == nil {
		m.carry = o
		o.setOwner(m)
	}
}

//=============================================================
// Throw/drop object
//=============================================================
func (m *mob) throw() {
	if m.carry != nil {
		m.carry.owner = nil
		m.carry = nil
	}
}

//=============================================================
//
//=============================================================
func (m *mob) explode() {
	for i := 0; i < len(m.frames); i++ {
		for x := 0; x < m.frameWidth; x++ {
			for y := 0; y < m.frameHeight; y++ {
				pos := m.size*x + y
				if m.frames[i][pos] != 0 {
					// Remove part
					global.gParticleEngine.newParticle(
						particle{
							x:           m.bounds.X + float64(x),
							y:           m.bounds.Y + float64(y),
							size:        1,
							restitution: -0.1 - global.gRand.randFloat()/4,
							life:        wParticleDefaultLife,
							fx:          10,
							fy:          10,
							vx:          float64(5 - global.gRand.rand()),
							vy:          float64(5 - global.gRand.rand()),
							mass:        1,
							pType:       particleRegular,
							color:       m.frames[i][pos],
							static:      true,
						})
				}
				m.frames[i][pos] = 0
			}
		}
	}
	m.buildFrames()
}

//=============================================================
//
//=============================================================
func (m *mob) move(x, y float64) {
	m.keyMove.X = x
	m.keyMove.Y = y

	if x != 0 {
		if x > 0 {
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
// Get bounds
//=============================================================
func (m *mob) getBounds() *Bounds {
	return m.bounds
}

//=============================================================
//
//=============================================================
func (m *mob) getType() entityType {
	return m.mobType
}

//=============================================================
//
//=============================================================
func (m *mob) setPosition(x, y float64) {
	m.bounds.X = x
	m.bounds.Y = y
}

//=============================================================
//
//=============================================================
func (m *mob) hitCeiling(x, y float64) bool {
	for px := 0.0; px < m.bounds.Width; px++ {
		if global.gWorld.IsRegular(x+px, y+m.bounds.Height+1) {
			return true
		}
	}
	return false
}

//=============================================================
//
//=============================================================
func (m *mob) hitFloor(x, y float64) bool {
	for px := 0.0; px < m.bounds.Width; px++ {
		if global.gWorld.IsRegular(x+px, y+1) {
			return true
		}
	}
	return false
}

//=============================================================
//
//=============================================================
func (m *mob) hitWallRight(x, y float64) bool {
	for py := m.bounds.Height / 3; py < m.bounds.Height; py++ {
		if global.gWorld.IsRegular(x-2, y+py) {
			return true
		}
	}
	return false
}

//=============================================================
//
//=============================================================
func (m *mob) hitWallLeft(x, y float64) bool {
	for py := m.bounds.Height / 3; py < m.bounds.Height; py++ {
		if global.gWorld.IsRegular(x+m.bounds.Width+1, y+py) {
			return true
		}
	}
	return false
}

//=============================================================
// Physics for mob.
// I don't want real physics, better to have a good feeling for
// movement than accurate physic simulation.
//=============================================================
func (m *mob) physics(dt float64) {
	if m.keyMove.X != 0 {
		m.velo.X = dt * m.speed * m.dir
		m.velo.X = m.dir * math.Max(math.Abs(m.velo.X), m.speed/100)
	} else {
		if m.hitFloor(m.bounds.X, m.bounds.Y-5) {
			m.velo.X = math.Max(math.Abs(m.velo.X)-dt*m.speed/10, 0) * m.dir
		} else {
			m.velo.X = math.Max(math.Abs(m.velo.X)-dt*m.speed/100, 0) * m.dir
		}
	}

	m.climbing = false
	m.velo.Y += wGravity * dt
	m.velo.Y = math.Max(m.velo.Y, wGravity)
	if m.keyMove.Y > 0 {
		if m.IsOnLadder() {
			m.velo.Y = m.speed / 2 * dt
			m.climbing = true
			m.velo.X = 0
		} else {
			if !m.jumping {
				m.velo.Y = m.jumpPower
				m.jumping = true
			}
		}
	}

	if m.velo.Y != 0 {
		if m.velo.Y > 0 {
			if !m.hitCeiling(m.bounds.X, m.bounds.Y+m.velo.Y) {
				m.bounds.Y += m.velo.Y
			} else {
				m.velo.Y = 0
			}
		} else {
			if !m.hitFloor(m.bounds.X, m.bounds.Y+m.velo.Y) {
				m.bounds.Y += m.velo.Y
			} else {
				m.velo.Y = 0
				m.jumping = false
			}
		}
	}

	if m.velo.X != 0 {
		if m.velo.X > 0 {
			if !m.hitWallLeft(m.bounds.X+m.velo.X, m.bounds.Y+m.velo.Y) {
				m.bounds.X += m.velo.X
			} else {
				m.velo.X = 0
			}
		} else {
			if !m.hitWallRight(m.bounds.X+m.velo.X, m.bounds.Y+m.velo.Y) {
				m.bounds.X += m.velo.X
			} else {
				m.velo.X = 0
			}
		}
	}

	if m.climbing {
		m.currentAnim = animClimb
	} else if m.jumping {
		m.currentAnim = animJump
	} else if m.velo.X != 0 {
		m.currentAnim = animWalk
	} else {
		m.currentAnim = animIdle
	}

	m.keyMove.X = 0
	m.keyMove.Y = 0
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
//
//=============================================================
func (m *mob) draw(dt float64) {
	shooting := false
	if m.currentAnim == animShoot {
		shooting = true
	}
	// Update physics
	m.physics(dt)

	if shooting {
		m.currentAnim = animShoot
	}

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

	// Draw any object attached.
	if m.carry != nil {
		m.carry.canvas.Draw(global.gWin, pixel.IM.ScaledXY(pixel.ZV, pixel.V(m.carry.scale*m.dir, m.carry.scale)).
			Moved(pixel.V(m.bounds.X+m.bounds.Width/2, m.bounds.Y+m.bounds.Height/2-2)).
			Rotated(pixel.Vec{m.carry.bounds.X + m.carry.bounds.Width/2, m.carry.bounds.Y + m.carry.bounds.Height/2}, m.carry.rotation*m.dir))
		// Update object positions based on mob
		m.carry.bounds.X = m.bounds.X
		m.carry.bounds.Y = m.bounds.Y
	}
}
