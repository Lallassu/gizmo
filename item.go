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

	switch iType {
	case itemPlant:
	case itemBucket:
	case itemCrate:
		i.sheetFile = fmt.Sprintf("%v%v", wAssetObjectsPath, "crate2.png")
		i.animated = false
		i.name = "crate"
		i.phys.scale = 1
	case itemBarrel:
	}
	i.create(x, y)
	i.bounds.entity = Entity(i)
}
