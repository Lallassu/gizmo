package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// controller handle input from devices
type controller struct {
	quit         bool
	entity       entity
	moveLeftKey  pixelgl.Button
	moveRightKey pixelgl.Button
	moveClimbKey pixelgl.Button
	moveJumpKey  pixelgl.Button
	lightDt      float64
}

// Set which entity to control
func (c *controller) setActiveEntity(e entity) {
	switch item := e.(type) {
	case *mob:
		c.entity = item
	}
}

// Initialize controls
func (c *controller) create() {
	c.quit = false
}

// Handle input for both mouse and keyboard
func (c *controller) update(dt float64) {
	if global.gWin.Pressed(pixelgl.KeyQ) {
		c.quit = true
	}

	// If main menu is visible, just stear the menu and not char.
	if global.gActiveMenu != nil {
		if global.gWin.JustPressed(pixelgl.KeyUp) {
			global.gActiveMenu.moveUp()
		}
		if global.gWin.JustPressed(pixelgl.KeyDown) {
			global.gActiveMenu.moveDown()
		}
		if global.gWin.JustPressed(pixelgl.KeyEnter) {
			global.gActiveMenu.selectItem()
		}
		if global.gWin.JustPressed(pixelgl.KeyEscape) {
			global.gActiveMenu = nil
		}
		return
	}

	// Global not bound to entity
	if global.gWin.Pressed(pixelgl.KeyM) {
		//	PrintMemoryUsage()
		global.gActiveMenu = global.gMainMenu
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

	move := pixel.Vec{X: 0, Y: 0}

	// TEST
	if global.gWin.Pressed(pixelgl.KeyK) {
		//global.gWorld.Explode(global.gPlayer.bounds.X, global.gPlayer.bounds.Y+20, 20)
	}
	c.lightDt += dt
	if global.gWin.Pressed(pixelgl.KeyL) {
		if c.lightDt > 1 {
			l := &light{}
			pos := global.gPlayer.getPosition()
			l.create(pos.X, pos.Y, 300, 360, 100, pixel.RGBA{R: 0.8, G: 0.6, B: 0, A: 0.3}, true, 1)
			c.lightDt = 0
		}
	}

	// Go into door
	if global.gWin.Pressed(pixelgl.KeyA) {
		if m, ok := c.entity.(*mob); ok {
			m.action()
		}
	}

	// Test pickup
	if global.gWin.Pressed(pixelgl.KeyB) {
		if m, ok := c.entity.(*mob); ok {
			m.pickup()
		}
	}

	// Throw object
	if global.gWin.Pressed(pixelgl.KeyV) {
		if m, ok := c.entity.(*mob); ok {
			m.throw()
		}
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

	if m, ok := c.entity.(*mob); ok {
		m.move(move.X, move.Y)
	}

	// Handle mouse
	if global.gWin.Pressed(pixelgl.MouseButtonLeft) || global.gWin.Pressed(pixelgl.KeyLeftShift) {
		if m, ok := c.entity.(*mob); ok {
			m.shoot()
		}
	}
}
