//=============================================================
// interfaces.go
//-------------------------------------------------------------
// Interfaces
//=============================================================
package main

import (
	_ "github.com/faiface/pixel"
)

type Entity interface {
	hit(x, y float64) bool
	explode()
	getMass() float64
	getType() entityType
	draw(dt float64)
	move(x, y float64)
}
