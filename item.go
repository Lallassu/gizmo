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
	iType itemType
}

//=============================================================
//
//=============================================================
func (i *item) newItem(x, y float64, iType itemType) {
	i.iType = iType
	switch iType {
	case itemPlant:
	case itemBucket:
	case itemCrate:
		i.sheetFile = fmt.Sprintf("%v%v", wAssetObjectsPath, "crate2.png")
		i.animated = false
		i.name = "crate"
		i.scale = 1
		//i.drawFunc = i.drawItem
		i.drawFunc = func(dt, e float64) {
		}
	case itemBarrel:
	}

	i.create(x, y)
	i.bounds.entity = Entity(i)
}

//=============================================================
// Custom draw function called after object.draw
//=============================================================
func (i *item) drawItem(dt, elapsed float64) {
}

//=============================================================
// custom explode function called after object.explode
//=============================================================
func (i *item) explodeFunc() {
}

//=============================================================
// custom hit function called after object.hit
//=============================================================
func (i *item) hitFunc(x, y, vx, vy float64, power int) {
}
