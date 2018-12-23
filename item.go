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
func (i *item) createItem(x, y float64, iType itemType) {
	i.iType = iType
	switch iType {
	case itemPlant:
	case itemBucket:
	case itemCrate:
		i.sheetFile = fmt.Sprintf("%v%v", wAssetObjectsPath, "crate2.png")
		i.animated = false
		i.name = "crate"
		i.scale = 1
		i.drawFunc = i.drawItem
	case itemBarrel:
	}

	i.create(x, y)
	i.bounds.entity = Entity(i)
}

func (i *item) drawItem(dt, elapsed float64) {
}
