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
		i.name = "crate"
		i.scale = 1
	}

	i.create(x, y)
	i.bounds.entity = Entity(i)
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
