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
	name        string
	oType       itemType
	owner       Entity
	reloadTime  float64
	animateIdle bool
	drawFunc    func(dt, elapsed float64)
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
	o.explode()
	return
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
func (o *object) throw() {

}

//=============================================================
//
//=============================================================
func (o *object) pickup() {
}

//=============================================================
//
//=============================================================
func (o *object) action() {
}

//=============================================================
//
//=============================================================
func (o *object) move(dx, dy float64) {
	o.phys.keyMove.X = dx
	o.phys.keyMove.Y = dy
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
func (o *object) getMass() float64 {
	return o.mass
}

//=============================================================
//
//=============================================================
func (o *object) setPosition(x, y float64) {
	o.bounds.X = x
	o.bounds.Y = y
}

//=============================================================
//
//=============================================================
func (o *object) setOwner(e Entity) {
	o.owner = e
}

//=============================================================
//
//=============================================================
func (o *object) removeOwner() {
	o.owner = nil
}

//=============================================================
// Get bounds
//=============================================================
func (o *object) getBounds() *Bounds {
	return o.bounds
}

//=============================================================
//
//=============================================================
func (o *object) draw(dt, elapsed float64) {
	// This should be kept in weapon somehow...
	// But currently weapon has no draw.
	o.reloadTime += dt

	o.canvas.Clear(pixel.RGBA{0, 0, 0, 0})
	o.batches[0].Draw(o.canvas)
	if o.owner == nil {
		o.physics(dt)
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
		owner := o.owner.(*mob)
		offset := 0.0
		switch o.bounds.entity.(type) {
		case *item:
			offset = 10.0
		}
		o.canvas.Draw(global.gWin, pixel.IM.ScaledXY(pixel.ZV, pixel.V(o.scale*owner.dir, o.scale)).
			Moved(pixel.V(owner.bounds.X+owner.bounds.Width/2, offset+owner.bounds.Y+owner.bounds.Height/2-2)).
			Rotated(pixel.Vec{o.bounds.X + o.bounds.Width/2, o.bounds.Y + o.bounds.Height/2}, o.rotation*owner.dir))
		// Update oect positions based on mob
		o.bounds.X = owner.bounds.X
		o.bounds.Y = owner.bounds.Y
	}

	// Call custom draw
	if o.drawFunc != nil {
		o.drawFunc(dt, elapsed)
	}
}
