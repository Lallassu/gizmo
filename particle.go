//=============================================================
// particle.go
//-------------------------------------------------------------
// Individual particles
//=============================================================
package main

import (
	"math"
)

type particle struct {
	x      float64
	y      float64
	size   float64
	color  uint32
	life   float64
	mass   float64
	pType  particleType
	static bool
	active bool
	fx     float64
	fy     float64

	restitution float64
	bounces     int
	vx          float64
	vy          float64
	prevX       float64
	prevY       float64
	mdt         float64
}

//=============================================================
// Update particle
//=============================================================
func (p *particle) update(dt float64) {
	if p.life <= 0 {
		p.active = false
		return
	}
	p.life -= dt
	ax := p.fx * dt * p.vx * p.mass
	ay := p.fy * dt * p.vy * p.mass

	p.prevX = p.x
	p.prevY = p.y

	// Take the largest
	// divide the other with it.
	// Lerp pixels (to get almost pixel perfect collisions)
	if p.pType != particleSmoke {
		lx := 1.0
		ly := 1.0
		loops := 1
		if math.Abs(ax) > 1 || math.Abs(ay) > 1 {
			if math.Abs(ax) > math.Abs(ay) {
				if ax < 0 {
					lx *= -1
				}
				ly = ay / math.Abs(ax)
				loops = int(math.Abs(ax))
			} else {
				if ay < 0 {
					ly *= -1
				}
				lx = ax / math.Abs(ay)
				loops = int(math.Abs(ay))
			}
		} else {
			lx = ax
			ly = ay
		}

		for n := 0; n < loops; n++ {
			if global.gWorld.IsWall(p.x+lx, p.y+ly) { // || global.gWorld.IsObject(p.x+lx, p.y+ly) {
				p.bounces++
				if p.vy < 0 {
					p.vy *= p.restitution
				} else {
					if p.vx > 0 {
						p.vx *= -p.restitution
						p.vy *= -p.restitution
					} else if p.vx < 0 {
						p.vx *= -p.restitution
						p.vy *= -p.restitution
					}
				}
				break
			} else {
				p.x += lx
				p.y += ly
			}
		}
	} else {
		p.x += ax
		p.y += ay
	}

	switch p.pType {
	case particleSmoke:
		p.size -= 1 * dt
	case particleFire:
		p.size -= 2 * dt
	}
	if p.size < 0 {
		p.stop()
		return
	}

	p.vy -= dt * p.fy

	if p.fx > 0 {
		p.fx -= dt * global.gWorld.gravity * p.mass
	} else {
		p.fx = 0
	}
	if p.fy > 0 && p.mass > 0 {
		p.fy -= dt * global.gWorld.gravity * p.mass
	}

	if p.prevX-p.x == 0 && p.prevY-p.y == 0 {
		p.mdt += dt
	} else {
		p.mdt = 0
	}

	if p.life <= 0 || p.mdt > 1 {
		p.stop()
	}
}

//=============================================================
// Stop particle
//=============================================================
func (p *particle) stop() {
	p.life = 0
	if p.static {
		if global.gWorld.IsRegular(p.x, p.y-1) {
			global.gWorld.AddPixel(int(p.x), int(p.y), p.color)
			global.gWorld.addShadow(int(p.x+1), int(p.y-1))
			// We could add shadow for full length, but more expensive and looks strange when adding it.
			// for i := 0.0; i < wShadowLength; i++ {
			// 	global.gWorld.addShadow(int(p.x+i), int(p.y-i))
			// }
		}
	}
}
