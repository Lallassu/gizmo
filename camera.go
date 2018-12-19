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
	zoom       float64
	pos        pixel.Vec
	scale      pixel.Vec
	follow     Entity
	cam        pixel.Matrix
	shakeDt    float64
	shakePower float64
}

func (c *camera) create() {
	//c.wScalePos = pixel.Vec{X: float64(global.gWorld.width / 2), Y: float64(global.gWorld.height / 2)}
	c.setPosition(0.0, 0.0)
	c.zoom = 1
}

func (c *camera) setFollow(e Entity) {
	c.follow = e
}

func (c *camera) setPosition(x, y float64) {
	c.pos = pixel.Vec{x, y}
}

func (c *camera) shake(pos pixel.Vec, power int) {
	// Distance to player from explosion
	dist := distance(global.gPlayer.getPosition(), pos)
	c.shakeDt = (float64(power) / dist) * 100
	c.shakePower = c.shakeDt * 10
	if c.shakePower < 5 {
		c.shakeDt = 0
		c.shakePower = 0
	} else if c.shakePower > 40 {
		c.shakeDt = 20
		c.shakePower = 20
	}
}

func (c *camera) update(dt float64) {
	pos := c.pos
	if c.follow != nil {
		pos = c.follow.getPosition()
		pos.X -= float64(global.gWindowWidth) / 2 / c.zoom
		pos.Y -= float64(global.gWindowHeight) / 2 / c.zoom
	}

	if c.shakeDt > 0 {
		pos.Y += c.shakePower/2 - global.gRand.randFloat()*c.shakePower
		pos.X += c.shakePower/2 - global.gRand.randFloat()*c.shakePower
		c.shakeDt -= 1
	} else {
		c.shakeDt = 0
	}

	pos = pixel.Lerp(c.pos, pos, 1-math.Pow(1.0/128, dt))
	c.cam = pixel.IM.Moved(pos.Scaled(-1 / c.zoom))
	c.cam = c.cam.Scaled(pos, c.zoom)
	global.gWin.SetMatrix(c.cam)
	c.pos = pos
}
