package main

import (
	"fmt"
	"image"

	"github.com/faiface/pixel"
)

type gameMap struct {
}

var colorDefinitions map[string]pixel.RGBA

func (m *gameMap) newMapFromImg() {
	var img image.Image
	var width float64
	var height float64
	img, width, height, _ = loadTexture(fmt.Sprintf("%v%v", wAssetMapsPath, "map1_fg.png"))

	size := width * height
	global.gWorld.NewMap(width, height, size)
	// Add "red" for each world piece.
	ww := int(width)
	hh := int(height)
	tmp := 0
	for x := 0; x <= ww; x++ {
		for y := 0; y <= hh; y++ {
			r, g, b, a := img.At(x, hh-y).RGBA()
			if a == 0 {
				global.gWorld.SetPixel(x, y, (r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | wBackground8&0xFF))
				continue
			}
			tmp++
			global.gWorld.SetPixel(x, y, (r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | 0xFF))
		}
	}
	Debug("TMP: ", tmp)

	global.gWorld.buildAllChunks()
	global.gPlayer.create(float64(100), float64(100))

	for i := 0.0; i < 10.0; i++ {
		w := &weapon{}
		w.newWeapon(100+(i*10), 100.0, weaponAk47)
	}
}

func (m *gameMap) newMap(level int) {
	mapFile := ""
	switch level {
	case 1:
		global.gMapColor.generateColors(mapEasy)
		mapFile = "map1.png"
	case 2:
	}

	var img image.Image
	var width float64
	var height float64
	img, width, height, _ = loadTexture(fmt.Sprintf("%v%v", wAssetMapsPath, mapFile))

	size := width * height

	// Initialize world.
	global.gWorld.NewMap(width, height, size)

	// Add "red" for each world piece.
	ww := int(width)
	hh := int(height)
	for x := 0; x <= ww; x++ {
		for y := 0; y <= hh; y++ {
			r, g, b, _ := img.At(x, hh-y).RGBA()
			add := false
			if r == 0xFFFF && g == 0 && b == 0 {
				add = true
			} else {
				for _, c := range global.gMapColor.entityCodes {
					if c.r == r && c.g == g && c.b == b {
						add = true
						break
					}
				}
			}
			if add {
				global.gWorld.SetPixel(x, y, uint32(0xFF0000FF))
			}
		}
	}

	// Paint map

	// Build all chunks
	global.gWorld.buildAllChunks()

	items := make(map[objectType][]pixel.Vec)
	for x := 0.0; x <= float64(global.gWorld.width); x++ {
		for y := 0.0; y <= float64(global.gWorld.height); y++ {
			r, g, b, _ := img.At(int(x), global.gWorld.height-int(y)).RGBA()

			for t, c := range global.gMapColor.entityCodes {
				if c.r == r && c.g == g && c.b == b {
					tooClose := false
					// This way we can add crates based on many pixels in the same place w/o floodfilling.
					for _, c := range items[t] {
						if distance(c, pixel.Vec{X: x, Y: y}) < 20 {
							tooClose = true
							break
						}
					}
					if !tooClose {
						items[t] = append(items[t], pixel.Vec{X: x, Y: y})
					}
				}
			}
		}
	}

	// NOTE: Z-order is applied depending on QT adding order.
	// Hence we add player + enemies first.

	for _, p := range items[itemDoor] {
		w := &item{}
		// Find floor.
		for {
			if global.gWorld.IsRegular(p.X, p.Y) {
				w.newItem(p.X, p.Y+1, itemDoor)
				break
			}
			p.Y--
		}
	}

	//pcgGen := pcg{}
	for _, p := range items[lampRegular] {
		// Find floor.
		for {
			if global.gWorld.IsRegular(p.X, p.Y) {
				l := &light{}
				//	pcgGen.GenerateLamp(int(p.X), int(p.Y-10))
				l.create(p.X, p.Y-1, -90, 100, 50, pixel.RGBA{R: 0.8, G: 0.6, B: 0, A: 0.3}, false, 0)
				break
			}
			p.Y++
		}
	}

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
			maxLife: 100.0,
			phys:    phys{speed: 100},
			ai:      &ai{},
		}
		test.create(p.X, p.Y)
	}

	// Place weapons
	for _, p := range items[weaponAk47] {
		w := &weapon{}
		w.newWeapon(p.X, p.Y, weaponAk47)
	}

	// Place crates
	for _, p := range items[itemCrate] {
		w := &item{}
		w.newItem(p.X, p.Y, itemCrate)
	}

	// Place portals
	for _, p := range items[itemPortal] {
		w := &item{}
		w.newItem(p.X, p.Y, itemPortal)
	}

	// Cluster mines
	for _, p := range items[explosiveClusterMine] {
		w := &explosive{}
		w.newExplosive(p.X, p.Y, explosiveClusterMine)
	}

	// Health Power ups
	for _, p := range items[itemPowerupHealth] {
		w := &item{}
		w.newItem(p.X, p.Y, itemPowerupHealth)
	}

	// Place random lamps
	//lamps := 20
	// pcgGen := pcg{}
	// for lamps != 0 {
	// 	radius := float64(100 + global.gRand.rand()*5)
	// 	if p, fit := global.gWorld.fitInWorld(10); fit {
	// 		for !global.gWorld.IsRegular(p.X, p.Y) {
	// 			p.Y++
	// 		}

	// 		// Check if lamp is close already.
	// 		skip := false
	// 		for _, b := range global.gWorld.qt.RetrieveIntersections(&Bounds{X: p.X - radius/2, Y: p.Y - radius/2, Width: radius * 2, Height: radius * 2}) {
	// 			switch b.entity.(type) {
	// 			case *light:
	// 				skip = true
	// 				break
	// 			}
	// 		}
	// 		if !skip {
	// 			l := &light{}
	// 			pcgGen.GenerateLamp(int(p.X), int(p.Y))
	// 			l.create(p.X, p.Y-5, -90, 100, radius, pixel.RGBA{0.8, 0.6, 0, 0.3}, false, 0)
	// 			lamps--
	// 		}
	// 	}
	// }
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

// paint generated map
// everything has to be performed in a specific order.
