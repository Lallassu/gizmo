// Different objects that are subject to physics and flood fill
// for destuction into pieces. Much like mob, but special adds
// like FF and different physics and actions.
package main

import (
	"math"

	"github.com/faiface/pixel"
)

type object struct {
	phys
	graphics
	light         *light
	lightOffsetX  float64
	lightOffsetY  float64
	name          string
	oType         objectType
	owner         *mob
	prevOwner     *mob
	animateIdle   bool
	animateOffset float64
	static        bool
}

func (o *object) create(x, y float64) {
	o.mass = 5
	//o.active = true
	o.createGfx(x, y, o.static)
	o.createPhys(x, y, o.frameWidth, o.frameHeight)
	o.graphics.scalexy = o.phys.scale
}

func (o *object) getType() objectType {
	return o.oType
}

func (o *object) hit(posX, posY, vx, vy float64, pow int) {
	if o.static {
		return
	}

	power := float64(pow)
	// If distance is close, explode, otherwise push away.
	dist := distance(pixel.Vec{X: posX + power/2, Y: posY + power/2}, pixel.Vec{X: o.bounds.X + o.bounds.Width/2, Y: o.bounds.Y + o.bounds.Height/2})
	if dist < float64(power*2) {
		o.explode()
	} else {
		if vx == 0 && vy == 0 {
			if posX < o.bounds.X {
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

// Add light to object
func (o *object) AddLight(x, y float64, l *light) {
	o.light = l
	o.lightOffsetX = x
	o.lightOffsetY = y
}

func (o *object) isFree() bool {
	if o.owner == nil {
		return true
	}
	return false
}

func (o *object) explode() {
	global.gWorld.qt.Remove(o.bounds)
	o.explodeGfx(o.bounds.X, o.bounds.Y, false)
	if o.owner != nil {
		o.owner.throw()
	}
}

func (o *object) getPosition() pixel.Vec {
	return pixel.Vec{X: o.bounds.X, Y: o.bounds.Y}
}

func (o *object) setOwner(m *mob) {
	o.owner = m
}

// Action (activate)
func (o *object) action(m *mob) {
}

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

func (o *object) draw(dt, elapsed float64) {

	if o.light != nil {
		o.light.bounds.X = o.bounds.X + o.lightOffsetX
		o.light.bounds.Y = o.bounds.Y + o.lightOffsetY
	}

	if !o.static {
		o.canvas.Clear(pixel.RGBA{R: 0, G: 0, B: 0, A: 0})
		idx := 0
		if o.animated {
			o.animCounter += dt
			idx = int(math.Floor(o.animCounter / 0.01))
			idx = o.idleFrames[idx/30%len(o.idleFrames)]
		}

		o.batches[idx].Draw(o.canvas)
	}

	if o.owner == nil && !o.static {
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
		o.animateOffset = 0.0
		if !(o.falling || !o.animateIdle) {
			// Animate up/down
			o.animateOffset = 5 + math.Sin(elapsed)*4

		}
		o.canvas.Draw(global.gWin, pixel.IM.ScaledXY(pixel.ZV, pixel.V(o.scale, o.scale)).Moved(pixel.V(o.bounds.X+o.bounds.Width/2, o.animateOffset+o.bounds.Y+o.bounds.Height/2)).Rotated(pixel.V(o.bounds.X+o.bounds.Width/2, o.bounds.Y+o.bounds.Height/2), o.rotation))
	} else {
		if !o.static {
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
				Rotated(pixel.Vec{X: o.bounds.X + o.bounds.Width/2, Y: o.bounds.Y + o.bounds.Height/2}, o.rotation*o.owner.dir))
			// Update oect positions based on mob
			o.bounds.X = o.owner.bounds.X
			o.bounds.Y = o.owner.bounds.Y
		} else {
			o.graphics.sprite.Draw(o.graphics.canvas, pixel.IM.Moved(pixel.V(o.bounds.Width/2, o.bounds.Height/2)))
			o.canvas.Draw(global.gWin, pixel.IM.ScaledXY(pixel.ZV, pixel.V(o.scale, o.scale)).Moved(pixel.V(o.bounds.X+o.bounds.Width/2, o.bounds.Y+o.bounds.Height/2)))
		}
	}

}
