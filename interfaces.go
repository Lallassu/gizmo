//=============================================================
// interfaces.go
//-------------------------------------------------------------
// Interfaces
//=============================================================
package main

import (
	"github.com/faiface/pixel"
)

type entity interface {
	hit(x, y, vx, vy float64, power int)
	draw(dt, elapsed float64)
	getPosition() pixel.Vec
}

type objectInterface interface {
	setOwner(m *mob)
	isFree() bool
	getType() objectType
	action(m *mob)
}
