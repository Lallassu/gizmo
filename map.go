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
	width := 0.0
	height := 0.0

	switch level {
	case 1:
		width = 2048.0
		height = 2048.0
		img = global.gWorld.NewMap(int(width), int(height), 100, 10, mapEasy)
	case 2:
	}

	crates := []pixel.Vec{}
	weapons := []pixel.Vec{}
	enemies := []pixel.Vec{}
	player := pixel.Vec{-1, -1}
	for x := 0.0; x <= width; x++ {
		for y := 0.0; y <= height; y++ {
			r, g, b, _ := img.At(int(x), int(height-y)).RGBA()

			// White = crate
			if r == 0xFFFF && g == 0xFFFF && b == 0xFFFF {
				tooClose := false
				// This way we can add crates based on many pixels in the same place w/o floodfilling.
				for _, c := range crates {
					if distance(c, pixel.Vec{x, y}) < 50 {
						tooClose = true
						break
					}
				}
				if !tooClose {
					crates = append(crates, pixel.Vec{x, y})
				}
				// Blue = enemies
			} else if r == 0x0 && g == 0x0 && b == 0xFFFF {
				tooClose := false
				// This way we can add crates based on many pixels in the same place w/o floodfilling.
				for _, c := range enemies {
					if distance(c, pixel.Vec{x, y}) < 50 {
						tooClose = true
						break
					}
				}
				if !tooClose {
					enemies = append(enemies, pixel.Vec{x, y})
				}
				// Yellow = weapons
			} else if r == 0xFFFF && g == 0xFFFF && b == 0x0 {
				tooClose := false
				// This way we can add crates based on many pixels in the same place w/o floodfilling.
				for _, c := range weapons {
					if distance(c, pixel.Vec{x, y}) < 50 {
						tooClose = true
						break
					}
				}
				if !tooClose {
					weapons = append(weapons, pixel.Vec{x, y})
				}
				// Green = Player
			} else if r == 0x0 && g == 0xFFFF && b == 0x0 {
				if player.X == -1 {
					player = pixel.Vec{x, y}
				}
			}
		}
	}

	// NOTE: Z-order is applied depending on QT adding order.
	// Hence we add player + enemies first.

	// place player
	global.gPlayer.create(float64(player.X), float64(player.Y))

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

	for _, e := range enemies {
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
		test.create(e.X, e.Y)
	}

	// Place weapons
	for _, e := range weapons {
		w := &weapon{}
		w.newWeapon(e.X, e.Y, ak47)
	}

	// Place crates
	for _, e := range crates {
		w := &item{}
		w.newItem(e.X, e.Y, itemCrate)
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
