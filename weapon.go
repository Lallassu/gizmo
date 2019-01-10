//=============================================================
// weapon.go
//-------------------------------------------------------------
// Implements different types of weapons
//=============================================================
package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"image"
)

//=============================================================
//
//=============================================================
type weapon struct {
	object
	shot             ammo
	wType            objectType
	bullets          int
	spread           float64
	reload           float64
	reloadTime       float64
	ammoBar          *pixel.Sprite
	ammoImg          *pixel.PictureData
	ammoCount        int
	currentAmmoCount int
}

//=============================================================
//
//=============================================================
func (w *weapon) newWeapon(x, y float64, wType objectType) {
	w.reloadTime = 0
	w.wType = wType

	ammoBarType := "ammobar.png"

	switch wType {
	case weaponAk47:
		w.ammoCount = 50
		ammoBarType = "ammobar.png"
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
			power: 5 + global.gRand.rand(), // 2
		}
		w.spread = 0.5
		w.bullets = 1
		w.reload = 0.05
	case weaponShotgun:
		w.ammoCount = 20
		ammoBarType = "ammobar.png"
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

	// Create HP Bar
	var img image.Image
	img, _, _, _ = loadTexture(fmt.Sprintf("%v%v", wAssetMixedPath, ammoBarType))
	w.ammoImg = pixel.PictureDataFromImage(img)
	w.ammoBar = pixel.NewSprite(w.ammoImg, pixel.R(0, 0, 40, 5))

	// Must change entity type in bounds for QT lookup
	w.bounds.entity = Entity(w)

	w.currentAmmoCount = w.ammoCount
}

//=============================================================
// Get Type
//=============================================================
func (w *weapon) getType() objectType {
	return w.wType
}

//=============================================================
//
//=============================================================
func (w *weapon) draw(dt, elapsed float64) {
	w.reloadTime += dt
	w.object.draw(dt, elapsed)

	// Draw ammobar
	if w.owner == nil {
		w.ammoBar.Draw(global.gWin, pixel.IM.Scaled(pixel.ZV, 0.3).Moved(pixel.V(w.bounds.X+w.bounds.Width/2, w.bounds.Y+w.bounds.Height+5)))
	} else {
		w.ammoBar.Draw(global.gWin, pixel.IM.Scaled(pixel.ZV, 0.3).Moved(pixel.V(w.owner.bounds.X+w.owner.bounds.Width/2, w.owner.bounds.Y+w.owner.bounds.Height+7)))
	}
}

//=============================================================
//
//=============================================================
func (w *weapon) shoot() {
	// Use mass = 5 and fx/fy = 0.5 for missile
	if w.owner == nil {
		return
	}
	if w.currentAmmoCount <= 0 {
		return
	}

	w.currentAmmoCount--
	w.ammoBar.Set(w.ammoImg, pixel.R(0, 0, 40*(float64(w.currentAmmoCount)/float64(w.ammoCount)), 5))

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

		global.gParticleEngine.effectSmoke(
			w.bounds.X+w.bounds.Width/2+w.owner.dir*3,
			w.bounds.Y+w.bounds.Height,
			1)

		global.gParticleEngine.ammoShell(
			w.bounds.X+w.bounds.Width/2+w.owner.dir*3,
			w.bounds.Y+w.bounds.Height,
			w.owner.dir,
			0.5)
		w.reloadTime = 0
	}
}
