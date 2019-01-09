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
	"math"
)

type Map struct {
}

var colorDefinitions map[string]pixel.RGBA

func (m *Map) newMap(level int) {
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
	m.paintMap()

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
			maxLife: 100.0,
			phys:    phys{speed: 100},
			ai:      &AI{},
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

	for _, p := range items[itemPowerupHealth] {
		w := &item{}
		w.newItem(p.X, p.Y, itemPowerupHealth)
	}

	// Place random lamps
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

//=============================================================
// paint generated map
// everything has to be performed in a specific order.
//=============================================================
func (m *Map) paintMap() {
	pcgGen := &pcg{}

	for x := 0; x < global.gWorld.width; x++ {
		for y := 0; y < global.gWorld.height; y++ {
			p := global.gWorld.pixels[x*global.gWorld.width+y]
			r := p >> 24 & 0xFF
			g := p >> 16 & 0xFF
			b := p >> 8 & 0xFF
			if r == 0xFF && g == 0x00 && b == 0x00 {
				v := global.gMapColor.backgroundSoft
				// add some alpha to background
				v &= wBackground32
				global.gWorld.pixels[x*global.gWorld.width+y] = v
			} else if r == 0x00 && g == 0x00 && b == 0x00 {
				v := global.gMapColor.background
				global.gWorld.pixels[x*global.gWorld.width+y] = v
			}
		}
	}

	// Ladders
	color := global.gMapColor.ladders
	for x := 0; x < global.gWorld.width; x++ {
		for y := 0; y < global.gWorld.height; y++ {
			if y+1 < global.gWorld.height && x+60 < global.gWorld.width && x > 0 && y > 0 {
				before := global.gWorld.pixels[(x-1)*global.gWorld.width+y] & 0xFF
				point := global.gWorld.pixels[x*global.gWorld.width+y] & 0xFF
				after := global.gWorld.pixels[(x+1)*global.gWorld.width+y] & 0xFF
				above := global.gWorld.pixels[x*global.gWorld.width+y+1] & 0xFF
				long := global.gWorld.pixels[(x+60)*global.gWorld.width+y] & 0xFF

				long2_above := global.gWorld.pixels[(x+22)*global.gWorld.width+y+1] & 0xFF
				long2_before := global.gWorld.pixels[(x+21)*global.gWorld.width+y] & 0xFF
				long2 := global.gWorld.pixels[(x+22)*global.gWorld.width+y] & 0xFF

				if ((above == wBackground8 || above == wShadow8) && point == 0xFF && before == 0xFF && (after == wBackground8 || after == wShadow8) && long == 0xFF) ||
					(above == wBackground8 && point == wBackground8 && after == wBackground8 && before == wBackground8 && long2 == 0xFF && long2_above == wBackground8 && long2_before == wBackground8) {
					for i := 5; i < 18; i++ {
						if i == 5 || i == 17 {
							for n := 0; n < 500000; n++ {
								if y-n > 0 {
									if (global.gWorld.pixels[(x+i)*global.gWorld.width+y-n]&0xFF == wBackground8 || global.gWorld.pixels[(x+i)*global.gWorld.width+y-n]&0xFF == wShadow8) && global.gWorld.pixels[(x+i)*global.gWorld.width+y-n]&0xFF != wLadder8 {
										global.gWorld.pixels[(x+i)*global.gWorld.width+y-n] = color & wLadder32
										// Shadows
										if global.gWorld.pixels[(x+i+1)*global.gWorld.width+y-n-1]&0xFF != 0xFF {
											global.gWorld.pixels[(x+i+1)*global.gWorld.width+y-n-1] &= 0x555555FF & wLadder32
										}
									} else {
										break
									}
								}
							}

						}
						for n := 0; ; n += 5 {
							if y-n > 0 {
								if global.gWorld.pixels[(x+i)*global.gWorld.width+y-n]&0xFF == wBackground8 || global.gWorld.pixels[(x+i)*global.gWorld.width+y-n]&0xFF == wShadow8 {
									global.gWorld.pixels[(x+i)*global.gWorld.width+y-n] = color & wLadder32
									// Dont shadow above walls
									if global.gWorld.pixels[(x+i+1)*global.gWorld.width+y-n-1]&0xFF != 0xFF {
										global.gWorld.pixels[(x+i+1)*global.gWorld.width+y-n-1] &= 0x555555FF & wLadder32
									}

								} else {
									break
								}
							} else {
								break
							}
						}
					}
				}
				// Make shadows
				if y-5 > 0 && x+5 < global.gWorld.width {
					below := global.gWorld.pixels[x*global.gWorld.width+y-1]
					right := global.gWorld.pixels[(x+1)*global.gWorld.width+y]
					if (below&0xFF == wShadow8 || below&0xFF == wBackground8) && point&0xFF == 0xFF {
						for i := 1; i < wShadowLength; i++ {
							p := global.gWorld.pixels[(x+i)*global.gWorld.width+y-i]
							if p&0xFF != wShadow8 && p&0xFF != 0xFF {
								r := uint32(math.Ceil(float64(p>>24&0xFF) / wShadowDepth))
								g := uint32(math.Ceil(float64(p>>16&0xFF) / wShadowDepth))
								b := uint32(math.Ceil(float64(p>>8&0xFF) / wShadowDepth))
								global.gWorld.pixels[(x+i)*global.gWorld.width+y-i] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | wShadow8&0xFF
							}
						}
					}
					if (right&0xFF == wShadow8 || right&0xFF == wBackground8) && point&0xFF == 0xFF {
						for i := 0; i < wShadowLength; i++ {
							p := global.gWorld.pixels[(x+i)*global.gWorld.width+y-i]
							if p&0xFF != wShadow8 && p&0xFF != 0xFF {
								r := uint32(math.Ceil(float64(p>>24&0xFF) / wShadowDepth))
								g := uint32(math.Ceil(float64(p>>16&0xFF) / wShadowDepth))
								b := uint32(math.Ceil(float64(p>>8&0xFF) / wShadowDepth))
								global.gWorld.pixels[(x+i)*global.gWorld.width+y-i] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | wShadow8&0xFF
							}
							p = global.gWorld.pixels[(x+i)*global.gWorld.width+y-i-1]
							if p&0xFF != wShadow8 && p&0xFF != 0xFF && i < 4 {
								r := uint32(math.Ceil(float64(p>>24&0xFF) / wShadowDepth))
								g := uint32(math.Ceil(float64(p>>16&0xFF) / wShadowDepth))
								b := uint32(math.Ceil(float64(p>>8&0xFF) / wShadowDepth))
								global.gWorld.pixels[(x+i)*global.gWorld.width+y-i-1] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | wShadow8&0xFF
							}
						}
					}
				}
			}

		}
	}

	// Create objects/materials AFTER ladders otherwise we must take objects into account
	// when generating ladders. Overhead with multiple loops, but easier.
	// First walls
	for x := 0; x < global.gWorld.width; x++ {
		for y := 0; y < global.gWorld.height; y++ {
			if y+1 < global.gWorld.height && x+1 < global.gWorld.width && x > 0 && y > 0 {
				before := global.gWorld.pixels[(x-1)*global.gWorld.width+y] & 0xFF
				point := global.gWorld.pixels[x*global.gWorld.width+y] & 0xFF
				after := global.gWorld.pixels[(x+1)*global.gWorld.width+y] & 0xFF
				above := global.gWorld.pixels[x*global.gWorld.width+y+1] & 0xFF
				below := global.gWorld.pixels[x*global.gWorld.width+y-1] & 0xFF

				// Roof
				if point == 0xFF && (below == wBackground8 || below == wShadow8) {
					pcgGen.MetalFlat(x, y, false)
				}
				// Walls
				if point == 0xFF && (after == wBackground8 || after == wShadow8) {
					pcgGen.MetalWall(x, y, false)
				}
				if point == 0xFF && (before == wBackground8 || before == wShadow8) {
					pcgGen.MetalWall(x, y, true)
				}
				// Floor
				if point == 0xFF && (above == wBackground8 || above == wShadow8) {
					pcgGen.MetalFlat(x, y, true)
				}
				// Colored floor
				if point == 0xFF && (above == wBackground8 || above == wShadow8) {
					pcgGen.MetalFloor(x, y)
				}
			}
		}
	}

	// Corners
	for x := 0; x < global.gWorld.width; x++ {
		for y := 0; y < global.gWorld.height; y++ {
			if y+1 < global.gWorld.height && x+1 < global.gWorld.width && x > 0 && y > 0 {
				before := global.gWorld.pixels[(x-1)*global.gWorld.width+y] & 0xFF
				point := global.gWorld.pixels[x*global.gWorld.width+y] & 0xFF
				after := global.gWorld.pixels[(x+1)*global.gWorld.width+y] & 0xFF
				above := global.gWorld.pixels[x*global.gWorld.width+y+1] & 0xFF
				below := global.gWorld.pixels[x*global.gWorld.width+y-1] & 0xFF
				cornerRight := global.gWorld.pixels[(x+1)*global.gWorld.width+y-1] & 0xFF
				cornerLeft := global.gWorld.pixels[(x-1)*global.gWorld.width+y+1] & 0xFF
				cornerRight2 := global.gWorld.pixels[(x+1)*global.gWorld.width+y+1] & 0xFF
				cornerLeft2 := global.gWorld.pixels[(x-1)*global.gWorld.width+y-1] & 0xFF

				// corner to the left downwards
				if point == 0xFF && (below == wBackground8 || below == wShadow8) && (before == wBackground8 || before == wShadow8) {
					pcgGen.MetalCornerDown(x, y, true)
				}
				if point == 0xFF && (below == wBackground8 || below == wShadow8) && (after == wBackground8 || after == wShadow8) {
					pcgGen.MetalCornerDown(x, y, false)
				}
				if point == 0xFF && (above == wBackground8 || above == wShadow8) && (after == wBackground8 || after == wShadow8) {
					pcgGen.MetalCornerUp(x, y, true)
				}
				if point == 0xFF && (above == wBackground8 || above == wShadow8) && (before == wBackground8 || before == wShadow8) {
					pcgGen.MetalCornerUp(x, y, false)
				}
				if point == 0xFF && after == 0xFF && (cornerRight == wShadow8 || cornerRight == wBackground8) && below == 0xFF && cornerLeft == 0xFF && above == 0xFF {
					pcgGen.MetalCornerRight(x, y, false)
				}
				if point == 0xFF && before == 0xFF && (cornerLeft2 == wShadow8 || cornerLeft2 == wBackground8) && below == 0xFF && cornerRight2 == 0xFF && above == 0xFF {
					pcgGen.MetalCornerRight(x, y, true)
				}
				if point == 0xFF && after == 0xFF && (cornerRight2 == wShadow8 || cornerRight2 == wBackground8) && above == 0xFF && cornerLeft2 == 0xFF && below == 0xFF {
					pcgGen.MetalCornerLeft(x, y, false)
				}
				if point == 0xFF && before == 0xFF && after == 0xFF && (cornerLeft == wShadow8 || cornerLeft == wBackground8) && above == 0xFF && cornerRight == 0xFF && below == 0xFF {
					pcgGen.MetalCornerLeft(x, y, true)
				}
			}
		}
	}
	// Cracks in the wall?
	// for x := 0; x < global.gWorld.width; x++ {
	// 	for y := 0; y < global.gWorld.height; y++ {
	// 		p := global.gWorld.pixels[x*global.gWorld.width+y] & 0xFF

	// 		if p == wShadow8 || p == wBackground8 {
	// 			if rand.Float64() < 0.0001 {
	// 				pcgGen.GenerateBricks(x, y)
	// 			}
	// 		}
	// 	}
	// }

	// Background gfx
	for x := 0; x < global.gWorld.width; x++ {
		for y := 0; y < global.gWorld.height; y++ {
			if y+30 < global.gWorld.height && x+1 < global.gWorld.width && x > 0 && y > 0 {
				p := global.gWorld.pixels[x*global.gWorld.width+y] & 0xFF
				pp := global.gWorld.pixels[x*global.gWorld.width+y+1] & 0xFF
				up := global.gWorld.pixels[x*global.gWorld.width+y+30] & 0xFF
				upLow := global.gWorld.pixels[x*global.gWorld.width+y+5] & 0xFF
				if p == 0xFF && (up == wBackground8 || up == wShadow8) && (pp == wShadow8 || pp == wBackground8) {
					pcgGen.GenerateLine(x, y+30)
				}
				if p == 0xFF && (upLow == wBackground8 || upLow == wShadow8) && (pp == wShadow8 || pp == wBackground8) {
					pcgGen.GenerateBottomLine(x, y+3)
				}
			}
		}
	}
	// Air intake on floor
	for x := 0; x < global.gWorld.width; x++ {
		for y := 0; y < global.gWorld.height; y++ {
			if y+30 < global.gWorld.height && x+1 < global.gWorld.width && x > 0 && y > 0 {
				p := global.gWorld.pixels[x*global.gWorld.width+y] & 0xFF
				pp := global.gWorld.pixels[x*global.gWorld.width+y+1] & 0xFF
				upLow := global.gWorld.pixels[x*global.gWorld.width+y+5] & 0xFF
				if p == 0xFF && (upLow == wBackground8 || upLow == wShadow8) && (pp == wShadow8 || pp == wBackground8) {
					pcgGen.GenerateBottomAirIntake(x, y)
				}
			}
		}
	}
	// Background gfx
	for x := 0; x < global.gWorld.width; x++ {
		for y := 0; y < global.gWorld.height; y++ {
			if y+wDoorHeight < global.gWorld.height && x+wDoorLen < global.gWorld.width && x > 0 && y > 0 {
				p := global.gWorld.pixels[x*global.gWorld.width+y] & 0xFF
				pafter := global.gWorld.pixels[(x+wDoorLen)*global.gWorld.width+y+1] & 0xFF
				pbelow := global.gWorld.pixels[(x+wDoorLen)*global.gWorld.width+y-1] & 0xFF
				//pbelowLong := global.gWorld.pixels[x*global.gWorld.width+y-wDoorHeight-55] & 0xFF
				pp := global.gWorld.pixels[x*global.gWorld.width+y+1] & 0xFF
				up := global.gWorld.pixels[x*global.gWorld.width+y+wDoorHeight] & 0xFF
				//down := global.gWorld.pixels[x*global.gWorld.width+y-1] & 0xFF

				// if p == 0xFF && pbelowLong == 0xFF && (down == wShadow8 || down == wBackground8) {
				// 	pcgGen.GenerateLamp(x, y-1)
				// }

				if pbelow != wBackground8 && pafter == wBackground8 && p == 0xFF && up == wBackground8 && pp == wBackground8 {
					// Check there is no ladder
					skip := false
					for i := 0; i < wDoorLen; i++ {
						if global.gWorld.pixels[(x+i)*global.gWorld.width+y+1]&0xFF == wLadder8 {
							skip = true
							break
						}
					}
					if !skip {
						if pcgGen.GenerateDoor(x, y+1) {
							// save door position
							global.gWorld.doors = append(global.gWorld.doors, pixel.Vec{float64(x + wDoorLen/2), float64(y + 1)})
						}
					}
				}
			}
		}
	}

}
