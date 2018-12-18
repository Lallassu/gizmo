package main

type weapon struct {
	object
	shot    ammo
	wType   weaponType
	bullets int
	spread  float64
	reload  float64
}

func (w *weapon) newWeapon(x, y float64, wType weaponType) {
	w.reloadTime = 0
	w.wType = wType
	switch wType {
	case ak47:
		w.textureFile = "assets/objects/ak47_weapon.png"
		w.name = "ak47"
		w.static = false
		w.scale = 0.15
		w.shot = ammo{
			color: 0xFFFF33FF,
			size:  0.5,
			life:  3.0,
			fx:    10.0,
			fy:    10.0,
			power: 2,
		}
		w.spread = 1
		w.bullets = 1
		w.reload = 0.05
	case shotgun:
		w.textureFile = "assets/objects/shotgun_weapon.png"
		w.name = "Shotgun"
		w.static = false
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
		w.spread = 10
		w.reload = 0.5
	}
	w.create(x, y)

	// Must change entity type in bounds for QT lookup
	w.bounds.entity = Entity(w)
}

func (w *weapon) shoot() {
	if w.reloadTime > w.reload {
		// Use mass = 5 and fx/fy = 0.5 for missile
		for i := 0; i < w.bullets; i++ {
			w.shot.x = w.bounds.X + w.bounds.Width/2 + w.owner.(*mob).dir*3
			w.shot.y = w.bounds.Y + w.bounds.Height + (w.spread - global.gRand.randFloat()*w.spread*2)
			w.shot.vx = 10.0 * w.owner.(*mob).dir
			w.shot.vy = 10.0 * w.rotation
			w.shot.mass = 6 + global.gRand.randFloat()*4
			w.shot.owner = w.owner
			global.gAmmoEngine.newAmmo(w.shot)
		}

		global.gParticleEngine.ammoShell(
			w.bounds.X+w.bounds.Width/2+w.owner.(*mob).dir*3,
			w.bounds.Y+w.bounds.Height,
			w.owner.(*mob).dir,
			0.5)
		w.reloadTime = 0
	}
}
