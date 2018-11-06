//=============================================================
// controller.go
//-------------------------------------------------------------
// Controllers (input) + AI
//=============================================================
package main

import (
	"github.com/faiface/pixel/pixelgl"
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
	if e.getType() != entityChunk {
		c.entity = e
		// TBD: Tell camera to follow this entity.
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
	// Handle controllers
	if global.gWin.Pressed(pixelgl.KeyS) {
		global.gCamera.zoom -= 0.05
	}
	if global.gWin.Pressed(pixelgl.KeyW) {
		global.gCamera.zoom += 0.05
	}
	if global.gWin.Pressed(pixelgl.KeyLeft) {
		//global.gCamera.pos.X -= 2.1
		c.entity.move(-dt, 0)
	}
	if global.gWin.Pressed(pixelgl.KeyRight) {
		//global.gCamera.pos.X += 2.1
		c.entity.move(dt, 0)
	}
	if global.gWin.Pressed(pixelgl.KeyUp) {
		//	global.gCamera.pos.Y += 2.1
		c.entity.move(0, dt)
	}
	if global.gWin.Pressed(pixelgl.KeyDown) {
		//	global.gCamera.pos.Y -= 2.1
		c.entity.move(0, -dt)
	}
	if global.gWin.Pressed(pixelgl.KeyM) {
		PrintMemoryUsage()
	}
	if global.gWin.Pressed(pixelgl.KeyQ) {
		c.quit = true
	}
	if global.gWin.Pressed(pixelgl.KeyP) {
		global.gCamera.setFollow(nil)
	}

	// Handle mouse
	if global.gWin.Pressed(pixelgl.MouseButtonLeft) {
		mouse := global.gCamera.cam.Unproject(global.gWin.MousePosition())
		global.gWorld.Explode(mouse.X, mouse.Y, 10)
		global.gParticleEngine.effectExplosion(mouse.X, mouse.Y, 10)
	}
}
