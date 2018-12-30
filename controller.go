//=============================================================
// controller.go
//-------------------------------------------------------------
// Controllers (input) + AI
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type controller struct {
	quit         bool
	entity       Entity
	moveLeftKey  pixelgl.Button
	moveRightKey pixelgl.Button
	moveClimbKey pixelgl.Button
	moveJumpKey  pixelgl.Button
	lightDt      float64
}

//=============================================================
// Set which entity to control
//=============================================================
func (c *controller) setActiveEntity(e Entity) {
	switch item := e.(type) {
	case *mob:
		c.entity = item
	}
}

//=============================================================
// Initialize controls
//=============================================================
func (c *controller) create() {
	c.quit = false
}

//=============================================================
// Handle input for both mouse and keyboard
//=============================================================
func (c *controller) update(dt float64) {
	// Global not bound to entity
	if global.gWin.Pressed(pixelgl.KeyM) {
		PrintMemoryUsage()
	}
	if global.gWin.Pressed(pixelgl.KeyQ) {
		c.quit = true
	}
	if global.gWin.Pressed(pixelgl.KeyP) {
		global.gCamera.setFollow(nil)
	}
	if global.gWin.Pressed(pixelgl.KeyG) {
		global.gWorld.gravity += 0.1
	}
	if global.gWin.Pressed(pixelgl.KeyH) {
		global.gWorld.gravity -= 0.1
	}

	// Controllers for entity
	if c.entity == nil {
		return
	}

	move := pixel.Vec{0, 0}

	// TEST
	if global.gWin.Pressed(pixelgl.KeyK) {
		//global.gWorld.Explode(global.gPlayer.bounds.X, global.gPlayer.bounds.Y+20, 20)
	}
	c.lightDt += dt
	if global.gWin.Pressed(pixelgl.KeyL) {
		if c.lightDt > 1 {
			l := &light{}
			pos := global.gPlayer.getPosition()
			l.create(pos.X, pos.Y, 300, 360, 100, pixel.RGBA{0.8, 0.6, 0, 0.3}, true, 1)
			c.lightDt = 0
		}
	}

	// Go into door
	if global.gWin.Pressed(pixelgl.KeyA) {
		c.entity.(*mob).action()
	}

	// Test pickup
	if global.gWin.Pressed(pixelgl.KeyB) {
		c.entity.(*mob).pickup()
	}

	// Throw object
	if global.gWin.Pressed(pixelgl.KeyV) {
		c.entity.(*mob).throw()
	}

	if global.gWin.Pressed(pixelgl.KeyS) {
		global.gCamera.zoom -= 0.05
	}
	if global.gWin.Pressed(pixelgl.KeyW) {
		global.gCamera.zoom += 0.05
	}
	if global.gWin.Pressed(pixelgl.KeyLeft) {
		move.X = -dt
	}
	if global.gWin.Pressed(pixelgl.KeyRight) {
		move.X = dt
	}
	if global.gWin.Pressed(pixelgl.KeyUp) {
		move.Y = dt
	}
	if global.gWin.Pressed(pixelgl.KeyDown) {
		move.Y = -dt
	}

	c.entity.(*mob).move(move.X, move.Y)

	// Handle mouse
	if global.gWin.Pressed(pixelgl.MouseButtonLeft) || global.gWin.Pressed(pixelgl.KeyLeftShift) {
		c.entity.(*mob).shoot()
	}
}
