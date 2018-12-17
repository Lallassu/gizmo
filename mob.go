//=============================================================
// mob.go
//-------------------------------------------------------------
// Anything that can move/be destroyed etc.
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image"
	"math"
)

type mob struct {
	sheetFile   string
	life        float64
	walkFrames  []int
	jumpFrames  []int
	climbFrames []int
	shootFrames []int
	idleFrames  []int
	frameWidth  float64
	frameHeight float64
	size        float64
	currentAnim animationType
	mobType     entityType
	animCounter float64
	animRate    float64
	speed       float64
	dir         float64
	mass        float64
	carry       *object
	velo        pixel.Vec
	climbing    bool
	jumping     bool
	jumpPower   float64
	keyMove     pixel.Vec

	batches      map[int]*pixel.Batch
	triangles    map[int]*pixel.TrianglesData
	frames       map[int][]uint32
	bounds       *Bounds
	img          image.Image
	canvas       *pixelgl.Canvas
	ai           *AI
	hitLeftWall  bool
	hitRightWall bool
}

//=============================================================
// Create mob
// - load animation sheet
//=============================================================
func (m *mob) create(x, y float64) {
	m.frames = make(map[int][]uint32)
	m.batches = make(map[int]*pixel.Batch)
	m.triangles = make(map[int]*pixel.TrianglesData)

	if m.ai != nil {
		m.ai.create(m)
	}

	m.animRate = 0.1
	m.jumpPower = 200
	if m.speed == 0 {
		m.speed = 200
	}
	m.mass = 20
	m.currentAnim = animIdle
	m.dir = 1

	fullWidth := 0.0

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
	for w := 0.0; w < fullWidth; w += m.frameWidth {
		m.frames[f] = make([]uint32, int(m.size)*int(m.size))
		for x := 0.0; x <= m.frameWidth; x++ {
			for y := 0.0; y <= m.frameHeight; y++ {
				r, g, b, a := m.img.At(int(w+x), int(m.frameHeight-y)).RGBA()
				m.frames[f][int(x*m.size+y)] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | a&0xFF
			}
		}
		m.triangles[f] = pixel.MakeTrianglesData(100)
		m.batches[f] = pixel.NewBatch(m.triangles[f], nil)
		f++
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
	v := 0
	rc := uint32(0)
	gc := uint32(0)
	bc := uint32(0)
	p2 := uint32(0)
	r1 := uint32(0)
	g1 := uint32(0)
	b1 := uint32(0)
	draw := 0
	same_x := 1.0
	same_y := 1.0
	pos := 0

	// Build batch for each frame.
	for i := 0; i < len(m.frames); i++ {
		for x := 0.0; x < float64(m.frameWidth); x++ {
			for y := 0.0; y < float64(m.frameHeight); y++ {
				p := m.frames[i][int(x*m.size+y)]
				if p == 0 || p&0xFF>>7 == 0 {
					continue
				}
				rc = p >> 24 & 0xFF
				gc = p >> 16 & 0xFF
				bc = p >> 8 & 0xFF
				same_x = 1.0
				same_y = 1.0

				for l := x + 1; l < m.bounds.Width; l++ {
					// Check color
					pos = int(l*m.size + y)
					p2 = m.frames[i][pos]
					if p2 == 0 {
						break
					}
					r1 = p2 >> 24 & 0xFF
					g1 = p2 >> 16 & 0xFF
					b1 = p2 >> 8 & 0xFF

					if r1 == rc && g1 == gc && b1 == bc && ((p2&0xFF)>>7) == 1 {
						// Same color and not yet visited!
						m.frames[i][pos] &= 0xFFFFFF7F
						same_x++
						new_y := 1.0
						for k := y; k < m.bounds.Height; k++ {
							pos = int(l*m.size + k)
							p2 = m.frames[i][pos]
							if p2 == 0 {
								break
							}
							r1 = p2 >> 24 & 0xFF
							g1 = p2 >> 16 & 0xFF
							b1 = p2 >> 8 & 0xFF

							if r1 == rc && g1 == gc && b1 == bc && ((p2&0xFF)>>7) == 1 {
								m.frames[i][pos] &= 0xFFFFFF7F
								new_y++
							} else {
								break
							}
						}
						if new_y < same_y {
							break
						} else {
							same_y = new_y
						}
					} else {
						break
					}
				}

				draw++

				// Convert to decimal
				r := float64(p>>24&0xFF) / 255.0
				g := float64(p>>16&0xFF) / 255.0
				b := float64(p>>8&0xFF) / 255.0
				a := float64(p&0xFF) / 255.0

				// Increase length of triangles if we need to draw more than we had before.
				if draw*6 >= len(*m.triangles[i]) {
					m.triangles[i].SetLen(draw*6 + 10)
				}

				// Size of triangle is given by how large the greedy algorithm found out.
				(*m.triangles[i])[v].Position = pixel.Vec{x, y}
				(*m.triangles[i])[v+1].Position = pixel.Vec{x + same_x, y}
				(*m.triangles[i])[v+2].Position = pixel.Vec{x + same_x, y + same_y}
				(*m.triangles[i])[v+3].Position = pixel.Vec{x, y}
				(*m.triangles[i])[v+4].Position = pixel.Vec{x, y + same_y}
				(*m.triangles[i])[v+5].Position = pixel.Vec{x + same_x, y + same_y}
				for n := 0; n < 6; n++ {
					(*m.triangles[i])[v+n].Color = pixel.RGBA{r, g, b, a}
				}

				v += 6

			}
		}
		// Reset the greedy bit

		for x := 0.0; x < m.bounds.Width; x++ {
			for y := 0.0; y < m.bounds.Height; y++ {
				pos = int(x*m.size + y)
				if m.frames[i][pos] != 0 {
					m.frames[i][pos] |= 0x00000080
				}
			}
		}
		m.triangles[i].SetLen(draw * 6)
		m.batches[i].Dirty()
	}
}

//=============================================================
//
//=============================================================
func (m *mob) hit(x_, y_, vx, vy float64, power int) bool {
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
					pos := int(m.size)*rx + ry
					if pos >= 0 && pos < int(m.size*m.size) {
						if m.frames[i][pos] != 0 {
							m.life -= 0.1
							if m.life <= 0 {
								m.explode()
								return true
							}
							if global.gRand.rand() < 1 {
								global.gParticleEngine.effectBlood(x_, y_, vx, vy, 1)
							}
							// Blood
							r := 175 + global.gRand.rand()*5
							g := 10 + global.gRand.rand()*2
							b := 10 + global.gRand.rand()*2
							a := global.gRand.rand() * 255
							m.frames[i][pos] = uint32(r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | a&0xFF)

							//global.gParticleEngine.newParticle(
							//	particle{
							//		x:           float64(x_) + global.gRand.randFloat(),
							//		y:           float64(y_) + global.gRand.randFloat(),
							//		size:        global.gRand.randFloat(),
							//		restitution: -0.1 - global.gRand.randFloat()/4,
							//		life:        wParticleDefaultLife,
							//		fx:          10,
							//		fy:          10,
							//		vx:          vx + global.gRand.randFloat(), //float64(5 - rand.Intn(10)),
							//		vy:          float64(5 - global.gRand.rand()),
							//		mass:        1,
							//		pType:       particleRegular,
							//		color:       m.frames[i][pos],
							//		static:      true,
							//	})
						}
					}
				}
			}
		}
	}

	m.buildFrames()

	// Check status of mob, if dead => remove from QT
	//	global.gWorld.qt.Remove(m.bounds)

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
//
//=============================================================
func (m *mob) pickup() {
	// Check if anything to pickup?
	for _, v := range global.gWorld.qt.RetrieveIntersections(m.bounds) {
		if v.entity.getType() == entityObject {
			if v.entity.(*object).isFree() {
				m.attach(v.entity.(*object))
				break
			}
		}
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
// Die
//=============================================================
func (m *mob) die() {
	// Drop weapon
	m.throw()
	m.life = 0
	global.gWorld.qt.Remove(m.bounds)
}

//=============================================================
//
//=============================================================
func (m *mob) explode() {
	for i := 0; i < len(m.frames); i++ {
		for x := 0.0; x < m.frameWidth; x++ {
			for y := 0.0; y < m.frameHeight; y++ {
				pos := int(m.size*x + y)
				if m.frames[i][pos] != 0 {
					// Remove part
					if global.gRand.rand() < 1 {
						global.gParticleEngine.effectBlood(m.bounds.X+float64(x), m.bounds.Y+float64(y), float64(5-global.gRand.rand()), float64(5-global.gRand.rand()), global.gRand.rand()/10)
						global.gParticleEngine.newParticle(
							particle{
								x:           m.bounds.X + float64(x),
								y:           m.bounds.Y + float64(y),
								size:        1,
								restitution: -0.1 - global.gRand.randFloat()/4,
								life:        wParticleDefaultLife,
								fx:          float64(15 - global.gRand.rand()),
								fy:          float64(15 - global.gRand.rand()),
								vx:          float64(5 - global.gRand.rand()),
								vy:          float64(5 - global.gRand.rand()),
								mass:        1,
								pType:       particleRegular,
								color:       m.frames[i][pos],
								static:      true,
							})
					}
				}
				m.frames[i][pos] = 0
			}
		}
	}
	m.die()
	//	m.buildFrames()
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
func (m *mob) hitWallLeft(x, y float64) bool {
	for py := m.bounds.Height / 2; py < m.bounds.Height; py++ {
		if global.gWorld.IsRegular(x-2, y+py) {
			m.hitRightWall = true
			return true
		}
	}
	m.hitRightWall = false
	return false
}

//=============================================================
//
//=============================================================
func (m *mob) hitWallRight(x, y float64) bool {
	for py := m.bounds.Height / 2; py < m.bounds.Height; py++ {
		if global.gWorld.IsRegular(x+m.bounds.Width+1, y+py) {
			m.hitLeftWall = true
			return true
		}
	}
	m.hitLeftWall = false
	return false
}

//=============================================================
// Check if on ladder
//=============================================================
func (m *mob) IsOnLadder() bool {
	for px := m.bounds.Width / 3; px < m.bounds.Width-m.bounds.Width/3; px += 2 {
		for py := 0.0; py < m.bounds.Height; py += 2 {
			if global.gWorld.IsLadder(m.bounds.X+px, m.bounds.Y+py) {
				return true
			}
		}
	}
	return false
}

//=============================================================
// Unstuck the objet if stuck.
//=============================================================
func (m *mob) unStuck(dt float64) {
	bottom := false
	top := false
	offset := 1.0
	// Check bottom pixels
	for x := m.bounds.X; x < m.bounds.X+m.bounds.Width; x += 2 {
		if global.gWorld.IsRegular(x, m.bounds.Y+offset) {
			bottom = true
			break
		}
	}

	//Check top pixels
	for x := m.bounds.X; x < m.bounds.X+m.bounds.Width; x += 2 {
		if global.gWorld.IsRegular(x, m.bounds.Y+m.bounds.Height-offset) {
			top = true
			break
		}
	}

	if bottom {
		m.bounds.Y += 10 * m.mass * dt
	} else if top {
		m.bounds.Y -= 10 * m.mass * dt
	}
}

//=============================================================
// Physics for mob.
// I don't want real physics, better to have a good feeling for
// movement than accurate physic simulation.
//=============================================================
func (m *mob) physics(dt float64) {
	if m.keyMove.X != 0 {
		m.velo.X = dt * m.speed * m.dir
	} else {
		if m.hitFloor(m.bounds.X, m.bounds.Y-5) {
			//m.velo.X = math.Max(math.Abs(m.velo.X)-dt*m.speed/10, 0) * m.dir
			m.velo.X = 0
		} else {
			m.velo.X = math.Min(math.Abs(m.velo.X)-dt*m.speed/100, 0) * m.dir
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
				m.velo.Y = m.jumpPower * dt
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
				if m.velo.Y < -6 {
					m.explode()
				}
				m.velo.Y = 0
				m.jumping = false
			}
		}
	}

	if m.velo.X != 0 {
		if m.velo.X > 0 {
			if !m.hitWallRight(m.bounds.X+m.velo.X, m.bounds.Y+m.velo.Y) {
				m.bounds.X += m.velo.X
			} else {
				m.velo.X = 0
			}
		} else {
			if !m.hitWallLeft(m.bounds.X+m.velo.X, m.bounds.Y+m.velo.Y) {
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
	go m.unStuck(dt)
}

//=============================================================
//
//=============================================================
func (m *mob) draw(dt, elapsed float64) {
	shooting := false
	if m.currentAnim == animShoot {
		shooting = true
	}

	// Update physics & AI
	m.physics(dt)
	if m.ai != nil {
		go func() {
			m.ai.update(dt, elapsed)
			m.hitRightWall = false
			m.hitLeftWall = false
		}()
	}

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

	//m.batches[idx].SetMatrix(pixel.IM.ScaledXY(pixel.ZV, pixel.V(-m.dir, 1)).Moved(pixel.V(m.bounds.X+m.bounds.Width/2, m.bounds.Y+m.bounds.Height/2)))
	m.canvas.Clear(pixel.RGBA{0, 0, 0, 0})
	//m.canvas.SetComposeMethod(pixel.ComposeOver)
	//m.canvas.SetColorMask(pixel.Alpha(1))
	m.batches[idx].Draw(m.canvas)
	m.canvas.Draw(global.gWin, (pixel.IM.ScaledXY(pixel.ZV, pixel.V(-m.dir, 1)).Moved(pixel.V(m.bounds.X+m.bounds.Width/2, m.bounds.Y+m.bounds.Height/2))))

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
