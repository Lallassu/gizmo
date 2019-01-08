//=============================================================
// explosive.go
//-------------------------------------------------------------
// Implements different types of explosives
//=============================================================
package main

import (
	"fmt"
	"github.com/faiface/pixel"
)

//=============================================================
//
//=============================================================
type explosive struct {
	object
	eType     objectType
	power     int
	countDown bool
	delayTime float64
	light     *light
}

func (e *explosive) newExplosive(x, y float64, eType objectType) {
	e.eType = eType
	switch eType {
	case explosiveRegularMine:
		e.sheetFile = fmt.Sprintf("%v%v", wAssetObjectsPath, "regularMine.png")
		e.name = "Regular Mine"
		e.animated = false
		e.rotation = 0
		e.power = 20
		e.delayTime = 2
		e.scale = 0.5
		e.light = &light{}
		e.light.create(x, y, 360, 360, 10, pixel.RGBA{0.8, 0, 0, 0.3}, true, 0)
		e.light.unlimitedLife = true
		e.light.ownerBounds = e.bounds
		e.light.blinkFrequency = 1
		e.AddLight(5, 6, e.light)
	case explosiveClusterMine:
		e.sheetFile = fmt.Sprintf("%v%v", wAssetObjectsPath, "regularMine.png")
		e.name = "Cluster Mine"
		e.animated = false
		e.rotation = 0
		e.power = 5
		e.delayTime = 2
		e.scale = 0.5
		e.light = &light{}
		e.light.create(x, y, 360, 360, 10, pixel.RGBA{0.9, 0.3, 0, 0.2}, true, 0)
		e.light.unlimitedLife = true
		e.light.ownerBounds = e.bounds
		e.light.blinkFrequency = 1
		e.AddLight(5, 6, e.light)
	}
	e.countDown = false
	e.light.objectCD = false

	e.create(x, y)

	// Animate up/down when idle
	e.animateIdle = false

	// Must change entity type in bounds for QT lookup
	e.bounds.entity = Entity(e)
}

//=============================================================
//
//=============================================================
func (e *explosive) draw(dt, elapsed float64) {
	e.object.draw(dt, elapsed)
	if distance(global.gPlayer.getPosition(), pixel.Vec{e.bounds.X, e.bounds.Y}) < 10 {
		e.countDown = true
		e.light.blinkFrequency = 0.1
	}

	if e.countDown {
		e.delayTime -= dt
		if e.delayTime <= 0 {
			e.light.destroy()
			global.gWorld.qt.Remove(e.bounds)
			global.gWorld.Explode(e.bounds.X, e.bounds.Y, e.power)

			switch e.eType {
			case explosiveClusterMine:
				for i := 0; i < 10; i++ {
					shot := ammo{
						color: 0xFF0000FF,
						size:  0.5,
						life:  1,
						fx:    10.0,
						fy:    10.0,
						power: 10,
					}
					shot.x = e.bounds.X + e.bounds.Width/2
					shot.y = e.bounds.Y + e.bounds.Height
					shot.vx = 10.0 * (0.5 - global.gRand.randFloat())
					shot.vy = 10.0 * global.gRand.randFloat()
					shot.mass = global.gRand.randFloat() * 4
					shot.owner = e
					global.gAmmoEngine.newAmmo(shot)
				}

			}
		}
	}
}

//=============================================================
//
//=============================================================
func (e *explosive) explode() {

}
