//=============================================================
// map.go
//-------------------------------------------------------------
// Generates world with items etc.
//=============================================================
package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"image"
)

type Map struct {
}

var colorDefinitions map[string]pixel.RGBA

func (m *Map) newMap(level int) {
	pcgGen := &pcg{}
	//weapons := 10 + global.gRand.rand()
	//enemies := 5 + global.gRand.rand()/2
	//items := 5 + global.gRand.rand()/2
	lamps := 5 + global.gRand.rand()/2
	//regularMines := 30

	var img image.Image
	switch level {
	case 1:
		global.gMapColor.generateColors(mapEasy)
		img = global.gWorld.NewMap(fmt.Sprintf("%v%v", wAssetMapsPath, "map1.png"))
	case 2:
	}

	items := make(map[objectType][]pixel.Vec)
	for x := 0.0; x <= float64(global.gWorld.width); x++ {
		for y := 0.0; y <= float64(global.gWorld.height); y++ {
			r, g, b, _ := img.At(int(x), global.gWorld.height-int(y)).RGBA()

			for t, c := range global.gMapColor.entityCodes {
				if c.r == r && c.g == g && c.b == b {
					tooClose := false
					// This way we can add crates based on many pixels in the same place w/o floodfilling.
					for _, c := range items[t] {
						if distance(c, pixel.Vec{x, y}) < 20 {
							tooClose = true
							break
						}
					}
					if !tooClose {
						items[t] = append(items[t], pixel.Vec{x, y})
					}
				}
			}
		}
	}

	// NOTE: Z-order is applied depending on QT adding order.
	// Hence we add player + enemies first.

	// Should only be one player though.
	for _, p := range items[mobPlayer] {
		global.gPlayer.create(float64(p.X), float64(p.Y))
	}

	for _, p := range items[mobEnemy1] {
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
	}

	// // Place weapons
	for _, p := range items[weaponAk47] {
		w := &weapon{}
		w.newWeapon(p.X, p.Y, weaponAk47)
	}

	// // Place crates
	for _, p := range items[itemCrate] {
		w := &item{}
		w.newItem(p.X, p.Y, itemCrate)
	}

	// Place random lamps
	for lamps != 0 {
		radius := float64(100 + global.gRand.rand()*5)
		if p, fit := global.gWorld.fitInWorld(10); fit {
			for !global.gWorld.IsRegular(p.X, p.Y) {
				p.Y++
			}

			// Check if lamp is close already.
			skip := false
			for _, b := range global.gWorld.qt.RetrieveIntersections(&Bounds{X: p.X - radius/2, Y: p.Y - radius/2, Width: radius * 2, Height: radius * 2}) {
				switch b.entity.(type) {
				case *light:
					skip = true
					break
				}
			}
			if !skip {
				l := &light{}
				pcgGen.GenerateLamp(int(p.X), int(p.Y))
				l.create(p.X, p.Y-5, -90, 100, radius, pixel.RGBA{0.8, 0.6, 0, 0.3}, false, 0)
				lamps--
			}
		}
	}
	// Place items
	// for regularMines != 0 {
	// 	if p, fit := global.gWorld.fitInWorld(10); fit {
	// 		w := &explosive{}
	// 		if global.gRand.rand() > 9 {
	// 			w.newExplosive(p.X, p.Y, explosiveClusterMine)
	// 		} else {
	// 			w.newExplosive(p.X, p.Y, explosiveRegularMine)
	// 		}
	// 		regularMines--
	// 	}
	// }
}
