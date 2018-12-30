//=============================================================
// map.go
//-------------------------------------------------------------
// Generates world with items etc.
//=============================================================
package main

import (
	"fmt"
	_ "github.com/faiface/pixel"
)

type Map struct {
}

func (m *Map) newMap(level int) {
	weapons := 0
	enemies := 0
	items := 0
	weapons = 10 + global.gRand.rand()
	enemies = 1 + global.gRand.rand()/2
	items = 10 + global.gRand.rand()/2

	switch level {
	case 1:
		global.gWorld.NewMap(1024, 1024, 100, 10, mapEasy)
	case 2:
	}

	// NOTE: Z-order is applied depending on QT adding order.
	// Hence we add player + enemies first.

	// place player
	player := 1
	for player != 0 {
		if p, fit := global.gWorld.fitInWorld(50); fit {
			global.gPlayer.create(float64(p.X), float64(p.Y))
			player--
		}
	}

	for enemies != 0 {
		if p, fit := global.gWorld.fitInWorld(50); fit {
			test := mob{
				graphics: graphics{
					sheetFile:   fmt.Sprintf("%v%v", wAssetMobsPath, "enemy1.png"),
					animated:    true,
					walkFrames:  []int{8, 9, 10, 11, 12, 13, 14},
					idleFrames:  []int{0, 2, 3, 4, 5, 6},
					shootFrames: []int{26},
					jumpFrames:  []int{15, 16, 17, 18, 19, 20},
					climbFrames: []int{1, 7},
					frameWidth:  12.0,
				},
				life: 100.0,
				phys: phys{speed: 100},
				ai:   &AI{},
			}
			test.create(p.X, p.Y)
			enemies--
		}
	}

	// Place weapons
	for weapons != 0 {
		if p, fit := global.gWorld.fitInWorld(50); fit {
			w := &weapon{}
			w.newWeapon(p.X, p.Y, ak47)
			weapons--
		}
	}

	// Place items
	for items != 0 {
		if p, fit := global.gWorld.fitInWorld(50); fit {
			w := &item{}
			w.createItem(p.X, p.Y, itemCrate)
			items--
		}
	}
}
