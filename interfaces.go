//=============================================================
// interfaces.go
//-------------------------------------------------------------
// Interfaces
//=============================================================
package main

import ()

type EntityIf interface {
	hit(x, y float64) bool
	explode()
	getPosition() pixel.Vec
	setPosition(x, y float64)
	getMass() float64
	intersects(EntityIf) bool
	getType() entityType
	// Pickup() bool
	// Drop() bool
	// Throw() bool
}
