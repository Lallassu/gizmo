//=============================================================
// item.go
//-------------------------------------------------------------
// Implements different types of items
//=============================================================
package main

import (
	"fmt"
)

//=============================================================
//
//=============================================================
type item struct {
	object
	iType objectType
}

//=============================================================
//
//=============================================================
func (i *item) newItem(x, y float64, iType objectType) {
	i.iType = iType
	switch iType {
	case itemCrate:
		i.sheetFile = fmt.Sprintf("%v%v", wAssetObjectsPath, "crate2.png")
		i.animated = false
		i.animateIdle = false
		i.name = "crate"
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

	// Must set this after create
	i.bounds.entity = Entity(i)
}

//=============================================================
// Get Type
//=============================================================
func (i *item) getType() objectType {
	return i.iType
}

//=============================================================
// Attach
//=============================================================
func (i *item) setOwner(m *mob) {
	switch i.iType {
	case itemPowerupHealth:
		// TBD: Powerup effect
		// TBD: Text how much power?
		// Remove object
		m.setLife(50) // TBD
		global.gWorld.qt.Remove(i.bounds)
		return
	}
	i.object.setOwner(m)
}

//=============================================================
// Custom draw function
//=============================================================
func (i *item) draw(dt, elapsed float64) {
	i.object.draw(dt, elapsed)
}

//=============================================================
// custom explode function called after object.explode
//=============================================================
func (i *item) explode() {
	i.object.explode()
}

//=============================================================
// custom hit function called after object.hit
//=============================================================
func (i *item) hit(x, y, vx, vy float64, power int) {
	i.object.hit(x, y, vx, vy, power)
}
