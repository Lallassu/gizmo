//=============================================================
// phys.go
//-------------------------------------------------------------
// Physics for MOBs (incl. player) and objects.
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"math"
)

type phys struct {
	bounds       *Bounds
	hitRightWall bool
	hitLeftWall  bool
	keyMove      pixel.Vec
	velo         pixel.Vec
	force        pixel.Vec
	speed        float64
	dir          float64
	climbing     bool
	jumping      bool
	jumpPower    float64
	mass         float64
	falling      bool
	rotation     float64
	restitution  float64
	scale        float64
	offset       float64
	throwable    bool
	moving       bool
	duck         bool
}

//=============================================================
//
//=============================================================
func (p *phys) createPhys(x, y, width, height float64) {
	// Initiate bounds for qt
	if p.scale == 0 {
		p.scale = 1
	}

	p.throwable = false

	p.bounds = &Bounds{
		X:      x,
		Y:      y,
		Width:  float64(width) * p.scale,
		Height: float64(height) * p.scale,
	}

	// Add object to QT
	global.gWorld.AddObject(p.bounds)

	// Create an offset for CD loops
	p.offset = width / 4
	if p.offset < 0 {
		p.offset = 1
	}

	if p.speed == 0 {
		p.speed = 100
	}
}

//=============================================================
//
//=============================================================
func (p *phys) hitCeiling(x, y float64) bool {
	for px := 0.0; px < p.bounds.Width; px += p.offset {
		if global.gWorld.IsRegular(x+px, y+p.bounds.Height+1) {
			return true
		}
	}
	return false
}

//=============================================================
//
//=============================================================
func (p *phys) hitFloor(x, y float64) bool {
	for px := 0.0; px < p.bounds.Width; px += 2 { // use instead of offset.
		if global.gWorld.IsRegular(x+px, y+1) {
			return true
		}
	}
	return false
}

//=============================================================
//
//=============================================================
func (p *phys) hitWallLeft(x, y float64) bool {
	for py := p.bounds.Height / 2; py < p.bounds.Height; py += p.offset {
		if global.gWorld.IsRegular(x-2, y+py) {
			p.hitRightWall = true
			return true
		}
	}
	p.hitRightWall = false
	return false
}

//=============================================================
//
//=============================================================
func (p *phys) hitWallRight(x, y float64) bool {
	for py := p.bounds.Height / 2; py < p.bounds.Height; py += p.offset {
		if global.gWorld.IsRegular(x+p.bounds.Width+1, y+py) {
			p.hitLeftWall = true
			return true
		}
	}
	p.hitLeftWall = false
	return false
}

//=============================================================
// Check if on ladder
//=============================================================
func (p *phys) IsOnLadder() bool {
	for px := p.bounds.Width / 3; px < p.bounds.Width-p.bounds.Width/3; px += p.offset {
		for py := 0.0; py < p.bounds.Height; py += 2 {
			if global.gWorld.IsLadder(p.bounds.X+px, p.bounds.Y+py) {
				return true
			}
		}
	}
	return false
}

//=============================================================
//=============================================================
func (p *phys) physics(dt float64) {
	if p.keyMove.X != 0 {
		p.velo.X = dt * p.speed * p.dir
	} else {
		if p.hitFloor(p.bounds.X, p.bounds.Y-1) {
			p.velo.X = 0
			p.moving = false
		} else {
			if p.throwable {
				p.velo.X += dt * p.speed / 100 * p.dir
				if p.velo.X != 0 {
					p.moving = true
				} else {
					p.moving = false
				}
			} else {
				p.velo.X = math.Min(math.Abs(p.velo.X)-dt*p.speed/100, 0) * p.dir
			}
		}
	}

	p.climbing = false
	p.velo.Y += wGravity * dt
	p.velo.Y = math.Max(p.velo.Y, wGravity)
	p.duck = false
	if p.keyMove.Y > 0 {
		if p.IsOnLadder() {
			p.velo.Y = p.speed / 2 * dt
			p.climbing = true
			p.velo.X /= 5
		} else {
			if !p.jumping {
				p.velo.Y = p.jumpPower * dt
				p.jumping = true
			}
		}
	} else if p.keyMove.Y < 0 {
		p.duck = true
	}

	p.falling = false
	if p.velo.Y != 0 {
		if p.velo.Y > 0 {
			if !p.hitCeiling(p.bounds.X, p.bounds.Y+p.velo.Y) {
				p.bounds.Y += p.velo.Y
			} else {
				p.velo.Y = 0
			}
		} else {
			if !p.hitFloor(p.bounds.X, p.bounds.Y+p.velo.Y) {
				// TBD: Check if hitting an object. Then place above object.
				p.bounds.Y += p.velo.Y
				p.falling = true
			} else {
				p.velo.Y = 0
				p.jumping = false
			}
		}
	}

	if p.velo.X != 0 {
		if p.velo.X > 0 {
			if !p.hitWallRight(p.bounds.X+p.velo.X, p.bounds.Y+p.velo.Y) {
				p.bounds.X += p.velo.X
			} else {
				p.velo.X = 0
			}
		} else {
			if !p.hitWallLeft(p.bounds.X+p.velo.X, p.bounds.Y+p.velo.Y) {
				p.bounds.X += p.velo.X
			} else {
				p.velo.X = 0
			}
		}
	}

	p.keyMove.X = 0
	p.keyMove.Y = 0
	p.unStuck(dt)
}

//=============================================================
// Unstuck the objet if stuck.
//=============================================================
func (p *phys) unStuck(dt float64) {
	bottom := false
	top := false
	offset := 1.0
	// Check bottom pixels
	for x := p.bounds.X; x < p.bounds.X+p.bounds.Width; x += p.offset {
		if global.gWorld.IsRegular(x, p.bounds.Y+offset) {
			bottom = true
			break
		}
	}

	if !bottom {
		//Check top pixels
		for x := p.bounds.X; x < p.bounds.X+p.bounds.Width; x += p.offset {
			if global.gWorld.IsRegular(x, p.bounds.Y+p.bounds.Height-offset) {
				top = true
				break
			}
		}
	}

	if bottom {
		p.bounds.Y += 10 * p.mass * dt
	} else if top {
		p.bounds.Y -= 10 * p.mass * dt
	}
}
