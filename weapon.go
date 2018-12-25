//=============================================================
// weapon.go
//-------------------------------------------------------------
// Implements different types of weapons
//=============================================================
package main

import (
	"fmt"
)

//=============================================================
//
//=============================================================
type weapon struct {
	object
	shot    ammo
	wType   weaponType
	bullets int
	spread  float64
	reload  float64
}

//=============================================================
//
//=============================================================
func (w *weapon) newWeapon(x, y float64, wType weaponType) {
	w.reloadTime = 0
	w.wType = wType
	switch wType {
	case ak47:
		w.sheetFile = fmt.Sprintf("%v%v", wAssetObjectsPath, "ak47_weapon.png")
		w.name = "ak47"
		w.animated = false
		w.rotation = 0.1
		w.scale = 0.15
		w.shot = ammo{
			color: 0xFFFF33FF,
			size:  0.5,
			life:  3.0,
			fx:    10.0,
			fy:    10.0,
			power: 1 + global.gRand.rand(), // 2
		}
		w.spread = 0.5
		w.bullets = 1
		w.reload = 0.05
	case shotgun:
		w.sheetFile = fmt.Sprintf("%v%v", wAssetObjectsPath, "shotgun_weapon.png")
		w.animated = false
		w.rotation = 0.1
		w.name = "Shotgun"
		w.scale = 0.15
		w.shot = ammo{
			color: 0xFFFF33FF,
			size:  0.2,
			life:  0.2,
			fx:    10.0,
			fy:    10.0,
			power: 1,
		}
		w.bullets = 10
		w.spread = 5
		w.reload = 0.5
	}
	w.create(x, y)

	// Animate up/down when idle
	w.animateIdle = true

	// Must change entity type in bounds for QT lookup
	w.bounds.entity = Entity(w)
}

//=============================================================
//
//=============================================================
func (w *weapon) shoot() {
	// Use mass = 5 and fx/fy = 0.5 for missile
	if w.owner == nil {
		return
	}
	if w.reloadTime > w.reload {
		//global.gSounds.play("shot")
		w.rotation = 0.1
		for i := 0; i < w.bullets; i++ {
			w.shot.x = w.bounds.X + w.bounds.Width/2 + w.owner.dir*3
			w.shot.y = w.bounds.Y + w.bounds.Height
			w.shot.vx = 10.0 * w.owner.dir
			w.shot.vy = 10.0*w.rotation + (w.spread - global.gRand.randFloat()*w.spread*2)
			w.shot.mass = 6 + global.gRand.randFloat()*4
			w.shot.owner = w.owner
			global.gAmmoEngine.newAmmo(w.shot)
		}

		global.gParticleEngine.ammoShell(
			w.bounds.X+w.bounds.Width/2+w.owner.dir*3,
			w.bounds.Y+w.bounds.Height,
			w.owner.dir,
			0.5)
		w.reloadTime = 0
	}
}
