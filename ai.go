//=============================================================
// ai.go
//-------------------------------------------------------------
// Stear an entity of type mob
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"math/rand"
)

//=============================================================
//
//=============================================================
type AI struct {
	entity     Entity
	dir_x      float64
	dir_y      float64
	weaponList []pixel.Vec
}

//=============================================================
//
//=============================================================
func (a *AI) create(e Entity) {
	a.entity = e
	a.dir_x = 0.01
}

//=============================================================
// Update the information where weapons exists in the world.
//=============================================================
func (a *AI) updateWeaponList(list *pixel.Vec) {

}

//=============================================================
// Try to find a weapon and go towards it.
//=============================================================
func (a *AI) findWeapon() {

	// If at weapon position. Call m.pickup()
}

//=============================================================
// Update AI
//=============================================================
func (a *AI) update(dt, time float64) {
	// TBD: assumes mob, handle with reflection
	m := a.entity.(*mob)

	if m.carry == nil {
		a.findWeapon()
	}

	if m.hitLeftWall {
		a.dir_x = -dt
	} else if m.hitRightWall {
		a.dir_x = dt
	} else if m.IsOnLadder() {
		a.dir_y = dt
	} else {
		a.dir_y = -dt
	}

	// Only wander around if have weapon otherwise we search for one.
	if m.carry != nil {
		if rand.Float64() < 0.01 {
			a.dir_x *= -1
		}
	}

	// Jump randomly
	if rand.Float64() < 0.01 {
		a.dir_y = dt
	}

	m.move(a.dir_x, a.dir_y)
}
