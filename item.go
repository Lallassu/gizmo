package main

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

type item struct {
	object
	iType objectType
}

func (i *item) newItem(x, y float64, iType objectType) {
	i.iType = iType
	switch iType {
	case itemPortal:
		i.sheetFile = fmt.Sprintf("%v%v", wAssetObjectsPath, "portal.png")
		i.animated = false
		i.animateIdle = false
		i.name = "Portal"
		i.scale = 0.25
	case itemCrate:
		i.sheetFile = fmt.Sprintf("%v%v", wAssetObjectsPath, "crate2.png")
		i.animated = false
		i.animateIdle = false
		i.name = "Crate"
		i.scale = 1
	case itemDoor:
		i.sheetFile = fmt.Sprintf("%v%v", wAssetObjectsPath, "door.png")
		i.static = true
		i.name = "Door"
		i.scale = 1
	case itemPowerupHealth:
		i.sheetFile = fmt.Sprintf("%v%v", wAssetObjectsPath, "poweruphp3.png")
		i.animated = false
		//i.idleFrames = []int{0, 1, 2, 3, 4, 5, 6}
		//	i.frameWidth = 32
		i.animateIdle = true
		i.name = "Powerup HP"
		i.scale = 0.5
	}
	i.create(x, y)

	// Test fragment shader (Must be set after gfx is created)
	uPos := mgl32.Vec2{float32((i.bounds.Width / 2) * (1 / i.scale)), float32((i.bounds.Height / 2) * (1 / i.scale))}
	if iType == itemPortal {
		i.graphics.canvas.SetUniform("uTime", &global.uTime)
		i.graphics.canvas.SetUniform("uPos", &uPos)
		i.graphics.canvas.SetFragmentShader(fragmentShaderPortal)
	} else if iType == itemDoor {
		i.graphics.canvas.SetUniform("uTime", &global.uTime)
		i.graphics.canvas.SetUniform("uPos", &uPos)
		i.graphics.canvas.SetFragmentShader(fragmentShaderDoor)
	}

	// Must set this after create
	i.bounds.entity = entity(i)
}

// Get Type
func (i *item) getType() objectType {
	return i.iType
}

// Action (activate)
func (i *item) action(m *mob) {
	switch i.iType {
	case itemPortal:
		// 1. Check if portal is close to a wall
		dir := 0.0
		for x := 0.0; x < 50; x++ {
			if global.gWorld.IsRegular(i.bounds.X+x, i.bounds.Y+i.bounds.Height/2) {
				dir = 1
				break
			} else if global.gWorld.IsRegular(i.bounds.X-x, i.bounds.Y+i.bounds.Height/2) {
				dir = -1
				break
			}
		}
		if dir != 0 {
			// 2. If so, check if there is a position on the other side where the mob fits.
			for x := 50.0; x < 150; x++ {
				if global.gWorld.IsRegular(i.bounds.X+(dir*x), i.bounds.Y+i.bounds.Height/2) {
					// 3. Move player.
					m.bounds.X = i.bounds.X + (dir * (x + 10))
					// 4. Destroy portal (explode)
					i.explode()
					break
				}
			}
		}
	}
}

// Attach
func (i *item) setOwner(m *mob) {
	switch i.iType {
	case itemPowerupHealth:
		i.evaporate(i.bounds.X, i.bounds.Y+i.animateOffset)
		m.setLife(50) // TBD
		m.graphics.hitTexts = append(m.graphics.hitTexts, &hitText{global.gFont.write(fmt.Sprintf("+%v", 50)), 3.0})
		global.gWorld.qt.Remove(i.bounds)
		return
	}
	i.object.setOwner(m)
}

// Custom draw function
func (i *item) draw(dt, elapsed float64) {
	// Set uniform used for shaders
	i.object.draw(dt, elapsed)
}

// custom explode function called after object.explode
func (i *item) explode() {
	i.object.explode()
}

// custom hit function called after object.hit
func (i *item) hit(x, y, vx, vy float64, power int) {
	i.object.hit(x, y, vx, vy, power)
}
