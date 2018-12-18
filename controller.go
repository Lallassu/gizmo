//=============================================================
// controller.go
//-------------------------------------------------------------
// Controllers (input) + AI
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math/rand"
)

type controller struct {
	quit         bool
	entity       Entity
	moveLeftKey  pixelgl.Button
	moveRightKey pixelgl.Button
	moveClimbKey pixelgl.Button
	moveJumpKey  pixelgl.Button
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
		//global.gCamera.pos.X -= 2.1
		//c.entity.move(-dt, 0)
		move.X = -dt
	}
	if global.gWin.Pressed(pixelgl.KeyRight) {
		//global.gCamera.pos.X += 2.1
		//c.entity.move(dt, 0)
		move.X = dt
	}
	if global.gWin.Pressed(pixelgl.KeyUp) {
		//	global.gCamera.pos.Y += 2.1
		//	c.entity.move(0, dt)
		//	global.gSounds.play("jump")
		move.Y = dt
	}
	if global.gWin.Pressed(pixelgl.KeyDown) {
		//	global.gCamera.pos.Y -= 2.1
		//c.entity.move(0, -dt)
		move.Y = -dt
	}
	if global.gWin.Pressed(pixelgl.KeyL) {
		c.entity.setPosition(float64(rand.Intn(global.gWorld.width)), float64(rand.Intn(global.gWorld.height)))
	}

	c.entity.move(move.X, move.Y)

	// Handle mouse
	if global.gWin.Pressed(pixelgl.MouseButtonLeft) || global.gWin.Pressed(pixelgl.KeyLeftControl) {
		//mouse := global.gCamera.cam.Unproject(global.gWin.MousePosition())
		// global.gWorld.Explode(mouse.X, mouse.Y, 10)
		// global.gParticleEngine.effectExplosion(mouse.X, mouse.Y, 10)
		c.entity.(*mob).shoot()
		//global.gSounds.play("shot")
	}
}
