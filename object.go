//=============================================================
// object.go
//-------------------------------------------------------------
// Different objects that are subject to physics and flood fill
// for destuction into pieces. Much like mob, but special adds
// like FF and different physics and actions.
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"math"
)

type object struct {
	phys
	graphics
	light          *light
	light_offset_x float64
	light_offset_y float64
	name           string
	oType          objectType
	owner          *mob
	prevOwner      *mob
	animateIdle    bool
}

//=============================================================
//
//=============================================================
func (o *object) create(x, y float64) {
	o.mass = 5
	//o.active = true
	o.createGfx(x, y)
	o.createPhys(x, y, o.frameWidth, o.frameHeight)
	o.graphics.scalexy = o.phys.scale
}

//=============================================================
//
//=============================================================
func (o *object) getType() objectType {
	return o.oType
}

//=============================================================
//
//  Function to implement Entity interface
//
//=============================================================
//=============================================================
//
//=============================================================
func (o *object) hit(x_, y_, vx, vy float64, pow int) {
	power := float64(pow)
	// If distance is close, explode, otherwise push away.
	dist := distance(pixel.Vec{x_ + power/2, y_ + power/2}, pixel.Vec{o.bounds.X + o.bounds.Width/2, o.bounds.Y + o.bounds.Height/2})
	if dist < float64(power*2) {
		o.explode()
	} else {
		if vx == 0 && vy == 0 {
			if x_ < o.bounds.X {
				o.dir = 1
				vx = power
			} else {
				vx = -power
				o.dir = -1
			}
		}
		o.throwable = true
		o.speed = dist * 2
		o.phys.keyMove.X = o.dir
		o.phys.velo.Y += math.Abs(power * 5 / dist)
	}

	if o.light != nil {
		o.light.destroy()
	}
	return
}

//=============================================================
// Add light to object
// x,y =>
//=============================================================
func (o *object) AddLight(x, y float64, l *light) {
	o.light = l
	o.light_offset_x = x
	o.light_offset_y = y
}

//=============================================================
//
//=============================================================
func (o *object) isFree() bool {
	if o.owner == nil {
		return true
	}
	return false
}

//=============================================================
//
//=============================================================
func (o *object) explode() {
	o.explodeGfx(o.bounds.X, o.bounds.Y, false)
	global.gWorld.qt.Remove(o.bounds)
	if o.owner != nil {
		o.owner.throw()
	}
}

//=============================================================
//
//=============================================================
func (o *object) getPosition() pixel.Vec {
	return pixel.Vec{o.bounds.X, o.bounds.Y}
}

//=============================================================
//
//=============================================================
func (o *object) setOwner(m *mob) {
	o.owner = m
}

//=============================================================
//
//=============================================================
func (o *object) removeOwner() {
	o.throwable = true
	o.speed = 300
	o.dir = o.owner.dir
	o.keyMove.X = o.owner.dir
	o.prevOwner = o.owner
	o.phys.velo.X += o.speed * o.owner.dir
	if math.Abs(o.owner.velo.X) > 0 {
		o.bounds.Y = o.owner.bounds.Y + o.owner.bounds.Height + 5
	}
	o.owner = nil
}

//=============================================================
//
//=============================================================
func (o *object) draw(dt, elapsed float64) {

	if o.light != nil {
		o.light.bounds.X = o.bounds.X + o.light_offset_x
		o.light.bounds.Y = o.bounds.Y + o.light_offset_y
	}

	o.canvas.Clear(pixel.RGBA{0, 0, 0, 0})

	idx := 0
	if o.animated {
		o.animCounter += dt
		idx = int(math.Floor(o.animCounter / 0.01))
		idx = o.idleFrames[idx/30%len(o.idleFrames)]
	}

	o.batches[idx].Draw(o.canvas)
	if o.owner == nil {
		o.physics(dt)

		// Check if thrown.
		if o.moving {
			for _, v := range global.gWorld.qt.RetrieveIntersections(o.bounds) {
				switch item := v.entity.(type) {
				case *mob:
					if item != o.prevOwner {
						item.hit(item.bounds.X, item.bounds.Y, o.velo.X, o.velo.Y, int(math.Abs(o.velo.X*o.velo.Y*2)))
					}
				}
			}
		}

		if o.falling {
			o.rotation += 0.1
		} else {
			o.rotation = 0
		}
		offset := 0.0
		if !(o.falling || !o.animateIdle) {
			// Animate up/down
			offset = 5 + math.Sin(elapsed)*4
		}
		o.canvas.Draw(global.gWin, pixel.IM.ScaledXY(pixel.ZV, pixel.V(o.scale, o.scale)).Moved(pixel.V(o.bounds.X+o.bounds.Width/2, offset+o.bounds.Y+o.bounds.Height/2)).Rotated(pixel.V(o.bounds.X+o.bounds.Width/2, o.bounds.Y+o.bounds.Height/2), o.rotation))
	} else {
		offset := 0.0
		switch o.bounds.entity.(type) {
		case *item:
			if !o.owner.duck {
				offset = 10.0
			} else {
				offset = 5.0
			}
		}
		o.canvas.Draw(global.gWin, pixel.IM.ScaledXY(pixel.ZV, pixel.V(o.scale*o.owner.dir, o.scale)).
			Moved(pixel.V(o.owner.bounds.X+o.owner.bounds.Width/2, offset+o.owner.bounds.Y+o.owner.bounds.Height/2-2)).
			Rotated(pixel.Vec{o.bounds.X + o.bounds.Width/2, o.bounds.Y + o.bounds.Height/2}, o.rotation*o.owner.dir))
		// Update oect positions based on mob
		o.bounds.X = o.owner.bounds.X
		o.bounds.Y = o.owner.bounds.Y
	}

}
