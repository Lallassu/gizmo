//=============================================================
// mob.go
//-------------------------------------------------------------
// Anything that can move/be destroyed etc.
//=============================================================
package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"image"
	"math"
	"math/rand"
)

type mob struct {
	phys
	graphics
	life     float64
	carry    interface{}
	ai       *AI
	drawFunc func(dt, elapsed float64)
	hpBar    *pixel.Sprite
	hpImg    *pixel.PictureData
	maxLife  float64
}

//=============================================================
// Create mob
//=============================================================
func (m *mob) create(x, y float64) {
	if m.ai != nil {
		m.ai.create(m)
	}

	m.jumpPower = 200
	if m.speed == 0 {
		m.speed = 200
	}
	m.mass = 20
	m.dir = 1

	m.life = m.maxLife

	// Create HP Bar
	var img image.Image
	img, _, _, _ = loadTexture(fmt.Sprintf("%v%v", wAssetMixedPath, "hpbar.png"))
	m.hpImg = pixel.PictureDataFromImage(img)
	m.hpBar = pixel.NewSprite(m.hpImg, pixel.R(0, 0, 40, 5))

	// Initiate the graphics for the mob
	m.createGfx(x, y)
	m.createPhys(x, y, m.frameWidth, m.frameHeight)

	m.graphics.scalexy = m.phys.scale

	// Set entity type for bounds.
	m.bounds.entity = m

}

//=============================================================
//
//=============================================================
func (m *mob) hit(x_, y_, vx, vy float64, power int) {
	pow := float64(power)
	// If distance is close, explode, otherwise push away.
	dist := distance(pixel.Vec{x_ + pow/2, y_ + pow/2}, pixel.Vec{m.bounds.X + m.bounds.Width/2, m.bounds.Y + m.bounds.Height/2})
	if dist < float64(power*2) {
		// If carry somerhing, hit that first!
		if m.carry != nil {
			switch item := m.carry.(type) {
			case *weapon:
				item.hit(x_, y_, vx, vy, power)
				return
			case *object:
				item.hit(x_, y_, vx, vy, power)
				return
			}
		}

		x := int(math.Abs(float64(m.bounds.X - x_)))
		y := int(math.Abs(float64(m.bounds.Y - y_)))

		// Gfx update
		m.hitGfx(x, y, x_, y_, vx, vy, power, true)

		// Blood effect
		global.gParticleEngine.effectBlood(x_, y_, vx, vy, 1)

		m.setLife(-float64(power * 2))
	} else {
		if vx == 0 && vy == 0 {
			if x_ < m.bounds.X {
				m.dir = 1
				vx = pow
			} else {
				vx = -pow
				m.dir = -1
			}
		}
		// Temprorary throwable (in order for shockwave effect)
		// Resets in move function
		m.phys.throwable = true
		m.phys.velo.Y += math.Abs(pow * 5 / dist)
		m.phys.velo.X += m.dir * pow * 5 / dist
	}
}

//=============================================================
// Set Health
//=============================================================
func (m *mob) setLife(change float64) {
	m.life += change
	if m.life > m.maxLife {
		m.life = m.maxLife
	}

	m.hpBar.Set(m.hpImg, pixel.R(0, 0, 40*(m.life/100), 5))
	if m.life <= 0 {
		m.die()
	}
}

//=============================================================
// Shoot if weapon attached
//=============================================================
func (m *mob) shoot() {
	if m.carry != nil {
		switch item := m.carry.(type) {
		case *weapon:
			item.shoot()
			m.currentAnim = animShoot
		}
	}
}

//=============================================================
// Attach object to self
//=============================================================
//func (m *mob) attach(o *object) {
func (m *mob) attach(o interface{}) {
	switch item := o.(type) {
	case *weapon:
		item.setOwner(m)
	case *object:
		item.setOwner(m)
	case *item:
		item.setOwner(m)
		// TBD: Handle this better.
		if item.iType == itemPowerupHealth {
			return
		}
	case *explosive:
		item.setOwner(m)
	}

	if m.carry == nil {
		m.carry = o
	}
}

//=============================================================
//
//=============================================================
func (m *mob) action() {
	// Check if close to doors
	for _, p := range global.gWorld.doors {
		d := distance(p, pixel.Vec{m.bounds.X, m.bounds.Y})
		if d < 10 {
			// Get a random new door position for player.
			pos := global.gWorld.doors[rand.Intn(len(global.gWorld.doors)-1)]
			m.bounds.X = pos.X
			m.bounds.Y = pos.Y
			break
		}
	}
}

//=============================================================
//
//=============================================================
func (m *mob) pickup() {
	// TBD: Handle this better...
	// Check if anything to pickup?
	for _, v := range global.gWorld.qt.RetrieveIntersections(m.bounds) {
		switch item := v.entity.(type) {
		case *object:
			if item.isFree() && m.carry == nil {
				m.attach(item)
				return
			}
		case *item:
			if item.isFree() && (m.carry == nil || item.iType == itemPowerupHealth) {
				m.attach(item)
				return
			}
		case *explosive:
			if item.isFree() && m.carry == nil {
				m.attach(item)
				return
			}
		case *weapon:
			if item.isFree() && m.carry == nil {
				m.attach(item)
				return
			}
		}
	}
}

//=============================================================
// Throw/drop object
//=============================================================
func (m *mob) throw() {
	if m.carry != nil {
		switch item := m.carry.(type) {
		case *object:
			item.removeOwner()
		case *weapon:
			item.removeOwner()
		case *item:
			item.removeOwner()
		case *explosive:
			item.removeOwner()
		}
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
	m.explodeGfx(m.bounds.X, m.bounds.Y, true)
	global.gWorld.qt.Remove(m.bounds)
}

//=============================================================
//
//=============================================================
func (m *mob) move(x, y float64) {
	m.phys.keyMove.X = x
	m.phys.keyMove.Y = y

	if x != 0 {
		m.phys.throwable = false
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
func (m *mob) draw(dt, elapsed float64) {
	shooting := false
	if m.currentAnim == animShoot {
		shooting = true
	}

	if m.velo.Y < -6 {
		// TBD: Fall to death, not explode
		// Or power?
	}

	// Update physics & AI
	m.physics(dt)

	if m.climbing {
		m.currentAnim = animClimb
	} else if m.jumping {
		m.currentAnim = animJump
	} else if m.velo.X != 0 {
		m.currentAnim = animWalk
	} else {
		m.currentAnim = animIdle
	}

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

	// Draw hpbar
	m.hpBar.Draw(global.gWin, pixel.IM.Scaled(pixel.ZV, 0.3).Moved(pixel.V(m.bounds.X+m.bounds.Width/2, m.bounds.Y+m.bounds.Height+5)))

	m.canvas.Draw(global.gWin, (pixel.IM.ScaledXY(pixel.ZV, pixel.V(-m.dir, 1)).Moved(pixel.V(m.bounds.X+m.bounds.Width/2, m.bounds.Y+m.bounds.Height/2))))

	// Call custom draw
	if m.drawFunc != nil {
		m.drawFunc(dt, elapsed)
	}

}
