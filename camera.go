//=============================================================
// camera.go
//-------------------------------------------------------------
// Controls the camera
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"math"
)

type camera struct {
	zoom      float64
	pos       pixel.Vec
	fPos      pixel.Vec
	scale     pixel.Vec
	follow    *Entity
	wScalePos pixel.Vec
	cam       pixel.Matrix
}

func (c *camera) create() {
	c.wScalePos = pixel.Vec{X: float64(global.gWorld.width / 2), Y: float64(global.gWorld.height / 2)}
	c.setPosition(0.0, 0.0)
	c.zoom = 1
}

func (c *camera) setFollow(entity *Entity) {
	c.follow = entity
}

func (c *camera) setPosition(x, y float64) {
	c.fPos = pixel.Vec{x, y}
}

func (c *camera) update(dt float64) {
	pos := c.fPos
	if c.follow != nil {
		//	pos = c.follow.getPosition()
	}
	pos = c.pos // TBD
	camPos := pixel.Lerp(c.pos, pos, 1-math.Pow(1.0/128, dt))
	c.cam = pixel.IM.Moved(camPos.Scaled(-1))
	//c.cam = c.cam.Moved(c.wScalePos)
	c.cam = c.cam.Scaled(c.wScalePos, c.zoom)
	global.gWin.SetMatrix(c.cam)
}
