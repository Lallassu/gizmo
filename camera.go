package main

import (
	"math"

	"github.com/faiface/pixel"
)

// camera controls the camera
type camera struct {
	zoom       float64
	pos        pixel.Vec
	scale      pixel.Vec
	follow     entity
	cam        pixel.Matrix
	shakeDt    float64
	shakePower float64
}

func (c *camera) create() {
	//c.wScalePos = pixel.Vec{X: float64(global.gWorld.width / 2), Y: float64(global.gWorld.height / 2)}
	c.setPosition(0.0, 0.0)
	c.zoom = 1
}

func (c *camera) setFollow(e entity) {
	c.follow = e
}

func (c *camera) setPosition(x, y float64) {
	c.pos = pixel.Vec{X: x, Y: y}
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
		pos.X -= float64(global.gVariableConfig.WindowWidth) / 2 / c.zoom
		pos.Y -= float64(global.gVariableConfig.WindowHeight) / 2 / c.zoom
	}

	if c.shakeDt > 0 {
		pos.Y += c.shakePower/2 - global.gRand.randFloat()*c.shakePower
		pos.X += c.shakePower/2 - global.gRand.randFloat()*c.shakePower
		c.shakeDt--
	} else {
		c.shakeDt = 0
	}

	pos = pixel.Lerp(c.pos, pos, 1-math.Pow(1.0/128, dt))
	c.cam = pixel.IM.Moved(pos.Scaled(-1 / c.zoom))
	c.cam = c.cam.Scaled(pos, c.zoom)
	global.gWin.SetMatrix(c.cam)
	c.pos = pos
}
