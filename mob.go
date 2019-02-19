package main

import (
	"fmt"
	"image"
	"math"
	"math/rand"

	"github.com/faiface/pixel"
)

type mob struct {
	phys
	graphics
	life     float64
	carry    interface{}
	ai       *ai
	drawFunc func(dt, elapsed float64)
	hpBar    *pixel.Sprite
	hpImg    *pixel.PictureData
	maxLife  float64
}

// Create mob
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
	m.createGfx(x, y, false)
	m.createPhys(x, y, m.frameWidth, m.frameHeight)

	m.graphics.scalexy = m.phys.scale

	// Set entity type for bounds.
	m.bounds.entity = m

}

func (m *mob) hit(posX, posY, vx, vy float64, power int) {
	// Create some text above the mob.
	m.graphics.hitTexts = append(m.graphics.hitTexts, &hitText{global.gFont.write(fmt.Sprintf("-%v", power*2)), 3.0})

	pow := float64(power)
	// If distance is close, explode, otherwise push away.
	dist := distance(pixel.Vec{X: posX + pow/2, Y: posY + pow/2}, pixel.Vec{X: m.bounds.X + m.bounds.Width/2, Y: m.bounds.Y + m.bounds.Height/2})
	if dist < float64(power*2) {
		// If carry somerhing, hit that first!
		if m.carry != nil {
			switch item := m.carry.(type) {
			case *weapon:
				item.hit(posX, posY, vx, vy, power)
				return
			case *object:
				item.hit(posX, posY, vx, vy, power)
				return
			}
		}

		x := int(math.Abs(float64(m.bounds.X - posX)))
		y := int(math.Abs(float64(m.bounds.Y - posY)))

		// Gfx update
		m.hitGfx(x, y, posX, posY, vx, vy, power, true)

		// Blood effect
		global.gParticleEngine.effectBlood(posX, posY, vx, vy, 1)

		m.setLife(-float64(power * 2))
	} else {
		if vx == 0 && vy == 0 {
			if posX < m.bounds.X {
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

// Shoot if weapon attached
func (m *mob) shoot() {
	if m.carry != nil {
		switch item := m.carry.(type) {
		case *weapon:
			item.shoot()
			m.currentAnim = animShoot
		}
	}
}

func (m *mob) pickup() {
	// Check if anything to pickup?
	var obj objectInterface
	for _, v := range global.gWorld.qt.RetrieveIntersections(m.bounds) {
		switch i := v.entity.(type) {
		case *object:
			obj = i
		case *item:
			obj = i
		case *explosive:
			obj = i
		case *weapon:
			obj = i
		}
	}
	if obj != nil {
		if obj.isFree() && m.carry == nil {
			if obj.getType() != itemPowerupHealth {
				m.carry = obj
			}
			obj.setOwner(m)
		}
	}
}

func (m *mob) action() {
	for _, v := range global.gWorld.qt.RetrieveIntersections(m.bounds) {
		switch i := v.entity.(type) {
		case *item:
			i.action(m)
			return
		}
	}

	// Check if close to doors
	for _, p := range global.gWorld.doors {
		d := distance(p, pixel.Vec{X: m.bounds.X, Y: m.bounds.Y})
		if d < 10 {
			// Get a random new door position for player.
			pos := global.gWorld.doors[rand.Intn(len(global.gWorld.doors)-1)]
			m.bounds.X = pos.X
			m.bounds.Y = pos.Y
			break
		}
	}
}

// Throw/drop object
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

// Die
func (m *mob) die() {
	// Drop weapon
	m.throw()
	m.life = 0
	m.explodeGfx(m.bounds.X, m.bounds.Y, true)
	global.gWorld.qt.Remove(m.bounds)
}

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

func (m *mob) getPosition() pixel.Vec {
	return pixel.Vec{X: m.bounds.X, Y: m.bounds.Y}
}

func (m *mob) draw(dt, elapsed float64) {
	shooting := false
	if m.currentAnim == animShoot {
		shooting = true
	}

	for _, c := range m.graphics.hitTexts {
		if c.ttl >= 0 {
			c.canvas.Draw(global.gWin, pixel.IM.Scaled(pixel.ZV, 0.1*c.ttl).Moved(pixel.V(m.bounds.X+10, m.bounds.Y+m.bounds.Height+10)))
			c.ttl -= dt * 2
		}
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
	m.canvas.Clear(pixel.RGBA{R: 0, G: 0, B: 0, A: 0})
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
