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
	oType          itemType
	owner          *mob
	prevOwner      *mob
	reloadTime     float64
	animateIdle    bool
	drawFunc       func(dt, elapsed float64)
	explodeFunc    func()
	hitFunc        func(x, y, vx, vy float64, power int)
}

//=============================================================
//
//=============================================================
func (o *object) create(x, y float64) {
	o.mass = 5
	//o.active = true
	o.animateIdle = false

	o.createGfx(x, y)
	o.createPhys(x, y, o.frameWidth, o.frameHeight)
	o.graphics.scalexy = o.phys.scale
}

//=============================================================
//
//  Function to implement Entity interface
//
//=============================================================
//=============================================================
//
//=============================================================
func (o *object) hit(x_, y_, vx, vy float64, power int) {
	// TBD: Remove explode and hit instead.
	o.explode()
	if o.hitFunc != nil {
		o.hitFunc(x_, y_, vx, vy, power)
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
	if o.explodeFunc != nil {
		o.explodeFunc()
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
	o.owner = nil
}

//=============================================================
//
//=============================================================
func (o *object) draw(dt, elapsed float64) {
	// This should be kept in weapon somehow...
	// But currently weapon has no draw.
	o.reloadTime += dt

	if o.light != nil {
		o.light.bounds.X = o.bounds.X + o.light_offset_x
		o.light.bounds.Y = o.bounds.Y + o.light_offset_y
	}

	o.canvas.Clear(pixel.RGBA{0, 0, 0, 0})
	o.batches[0].Draw(o.canvas)
	if o.owner == nil {
		o.physics(dt)

		// Check if thrown.
		if o.moving {
			for _, v := range global.gWorld.qt.RetrieveIntersections(o.bounds) {
				switch item := v.entity.(type) {
				case *mob:
					if item != o.prevOwner {
						item.hit(item.bounds.X, item.bounds.Y, o.velo.X, o.velo.Y, 20)
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
			offset = 5 + math.Sin(o.reloadTime)*3
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

	// Call custom draw
	if o.drawFunc != nil {
		o.drawFunc(dt, elapsed)
	}
}
