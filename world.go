//=============================================================
// world.go
//-------------------------------------------------------------
// Keep control of map (all pixels)
// Destuction of map
// Additions to map
// Generation of map
// Map flood fill
// Quadtree for entities
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"math/rand"
)

//=============================================================
// World Structure
//=============================================================
type world struct {
	width      int
	height     int
	size       int
	coloring   *mapColor
	qt         *Quadtree
	pixels     []uint32
	currentMap mapType
	gravity    float64
}

//=============================================================
//=============================================================
// World Public Functions
//=============================================================
//=============================================================

//=============================================================
// Initialize world first time.
//=============================================================
func (w *world) Init() {
	w.qt = &Quadtree{
		Bounds:     Bounds{X: 0, Y: 0, Width: float64(w.width), Height: float64(w.height)},
		MaxObjects: 4,
		MaxLevels:  8,
		Level:      0,
	}
	w.gravity = wGravity
}

//=============================================================
// New Map
//=============================================================
func (w *world) NewMap(maptype mapType) {
	// Generate map based on maptype
	w.currentMap = maptype
	w.qt.Clear()

	switch maptype {
	case mapEasy:
		w.width = 1024
		w.height = 1024
	case mapNormal:
		w.width = 1535
		w.height = 1536
	case mapHard:
		w.width = 2048
		w.height = 2048
	case mapHell:
	case mapWtf:
	}

	w.coloring = GenerateMapColor(maptype)

	w.size = w.width * w.height

	w.pixels = make([]uint32, w.size)

	Debug("Generating world", w.width, w.height)
	g := generator{}
	pixels := g.NewWorld(w.width, w.height)

	// Add all pixels as red before coloring
	for i := 0; i < len(pixels); i += 2 {
		w.AddPixel(int(pixels[i]), int(pixels[i+1]), uint32(0xFF0000FF))
	}

	// Paint the map with colors
	w.paintMap()

	// Initialize pixel pointers in the chunks
	for x := 0; x < w.width; x += wPixelsPerChunk {
		for y := 0; y < w.height; y += wPixelsPerChunk {
			c := &chunk{}
			c.create(float64(x), float64(y))
			w.qt.Insert(c.bounds)
		}
	}

	// Build all chunks first time.
	for _, v := range w.qt.RetrieveIntersections(&Bounds{X: 0, Y: 0, Width: float64(w.width), Height: float64(w.height)}) {
		v.entity.draw(-1)
	}

	Debug("Tree Size:", w.qt.Total)
	Debug("World generation complete.")
}

//=============================================================
// Add object to world (QT)
//=============================================================
func (w *world) AddObject(obj *Bounds) {
	w.qt.Insert(obj)
}

//=============================================================
// Remove object from world (QT)
//=============================================================
func (w *world) RemoveObject(obj Entity) {

}

//=============================================================
// Check if pixel is a background
//=============================================================
func (w *world) IsBackground(x_, y_ float64) bool {
	x := int(x_)
	y := int(y_)
	pos := w.width*x + y
	if pos < w.size && pos >= 0 {
		if w.pixels[pos]&0xFF == wBackground8 {
			return true
		}
	}
	return false
}

//=============================================================
// Check if pixel is a shadow
//=============================================================
func (w *world) IsShadow(x_, y_ float64) bool {
	x := int(x_)
	y := int(y_)
	pos := w.width*x + y
	if pos < w.width*w.height && pos >= 0 {
		if w.pixels[pos]&0xFF == wShadow8 {
			return true
		}
	}
	return false
}

//=============================================================
// Check if pixel is regular
//=============================================================
func (w *world) IsRegular(x_, y_ float64) bool {
	x := int(x_)
	y := int(y_)
	pos := w.width*x + y
	if pos < w.width*w.height && pos >= 0 {
		if w.pixels[pos]&0xFF == 0xFF {
			return true
		}
	}
	return false
}

//=============================================================
// Check if it's a wall
//=============================================================
func (w *world) IsWall(x_, y_ float64) bool {
	x := int(x_)
	y := int(y_)
	pos := w.width*x + y
	if pos < w.width*w.height && pos >= 0 {
		if w.pixels[pos] != 0 && w.pixels[pos]&0xFF != wBackground8 && w.pixels[pos]&0xFF != wShadow8 && w.pixels[pos]&0xFF != wLadder8 {
			return true
		}
	}
	return false
}

//=============================================================
// Check if it's a ladder
//=============================================================
func (w *world) IsLadder(x_, y_ float64) bool {
	x := int(x_)
	y := int(y_)
	pos := w.width*x + y
	if pos < w.width*w.height && pos >= 0 {
		if w.pixels[pos]&0xFF == wLadder8 {
			return true
		}
	}
	return false

}

//=============================================================
// Check if pixel exists
//=============================================================
func (w *world) PixelExists(x, y float64) bool {
	return true
}

//=============================================================
// Get color of the specified pixel
// Return -1 if not exist
//=============================================================
func (w *world) PixelColor(x, y float64) int32 {
	if x < 0 || y < 0 || x >= float64(w.width) || y >= float64(w.height) {
		return -1
	}
	return int32(w.pixels[uint32(int(x)*w.width+int(y))])
}

//=============================================================
// Draw
//=============================================================
func (w *world) Draw(dt float64) {
	// Draw those around camera position only.
	pos := pixel.Vec{0, 0}
	if global.gCamera.follow != nil {
		pos = global.gCamera.follow.getPosition()
	}
	for _, v := range w.qt.RetrieveIntersections(&Bounds{X: pos.X - wViewMax/2, Y: pos.Y - wViewMax/2, Width: wViewMax, Height: wViewMax}) {
		v.entity.draw(dt)
	}
}

//=============================================================
// Add pixel with color (replace if already exists)
//=============================================================
func (w *world) AddPixel(x, y int, color uint32) {
	pos := w.width*x + y
	if pos < w.width*w.height && pos >= 0 {
		w.pixels[w.width*x+y] = color
		w.markChunkDirty(x, y)
	}
}

//=============================================================
// Remove a pixel from the world map
//=============================================================
func (w *world) RemovePixel(x, y int) {
	pos := w.width*x + y
	if pos < w.width*w.height && pos >= 0 {
		if w.pixels[pos]&0xFF == wStaticColor8 ||
			w.pixels[pos]&0xFF == wBackground8 {
			return
		}

		//		if w.pixels[pos]&0xFF == wShadow8 {

		// Remove shadow
		for i := 0; i < 5; i++ {
			pos2 := (x+i)*w.width + y - i
			if pos2 < w.width*w.height && pos2 >= 0 {
				w.removeShadow(x+i, y-i)
			}
		}

		// Particle
		if w.IsRegular(float64(x), float64(y)) {
			global.gParticleEngine.newParticle(
				particle{
					x:           float64(x),
					y:           float64(y),
					size:        1,
					restitution: -0.1 - rand.Float64()/4,
					life:        wParticleDefaultLife,
					fx:          10,
					fy:          10,
					vx:          float64(5 - rand.Intn(10)),
					vy:          float64(5 - rand.Intn(10)),
					mass:        1,
					pType:       particleRegular,
					color:       w.pixels[pos],
					static:      true,
				})
		}

		// Set bg pixel.
		if w.pixels[pos] != 0 {
			v := w.coloring.getBackgroundSoft()
			v &= wBackground32
			w.pixels[pos] = v
		}
		w.markChunkDirty(x, y)
	}
}

//=============================================================
// Explode in world
// Also hits objects in the world.
//=============================================================
func (w *world) Explode(x_, y_ float64, power int) {
	x := int(x_)
	y := int(y_)
	pow := power * power
	ff := make([]pixel.Vec, 50)
	for rx := x - power; rx <= x+power; rx++ {
		vx := (rx - x) * (rx - x)
		for ry := y - power; ry <= y+power; ry++ {
			if ry < 0 {
				continue
			}
			val := (ry-y)*(ry-y) + vx
			if val < pow {
				w.RemovePixel(rx, ry)
				//w.ObjectHit(float64(rx), float64(ry))
				//for _, v := range w.qt.RetrieveIntersections(&Bounds{X: float64(rx), Y: float64(ry), Width: 1, Height: 1}) {
				//	v.entity.hit(x_, y_)
				//}
			} else {
				ff = append(ff, pixel.Vec{X: float64(rx), Y: float64(ry)})
			}
		}
	}

	for _, v := range w.qt.RetrieveIntersections(&Bounds{X: x_ - float64(power), Y: y_ - float64(power), Width: float64(power + power), Height: float64(power + power)}) {
		v.entity.hit(x_, y_, 0, 0, power)
	}

	// Add shadows
	for n := 0; n < len(ff); n++ {
		ffx := int(ff[n].X)
		ffy := int(ff[n].Y)
		pp := ffx*w.width + ffy
		if pp >= 0 && pp < w.width*w.height {
			if w.pixels[pp]&0xFF == 0xFF {
				for i := 0; i < 5; i++ {
					pos2 := (ffx+i)*w.width + ffy - i
					if pos2 < w.width*w.height && pos2 >= 0 {
						if w.pixels[pos2]&0xFF == wBackground8 {
							w.addShadow(ffx+i, ffy-i)
						}
					}
				}
			}
		}
	}

	// Floodfill
	// pixels := make([]Vec, 0)
	// for i := 0; i < len(ff); i++ {
	// 	pixels = append(pixels, w.FloodFill(ff[i].X, ff[i].Y)...)
	// }

	// for i := 0; i < len(pixels); i++ {
	// 	w.UnMarkPixelVisited(pixels[i].X, pixels[i].Y)
	// }}
}

//=============================================================
//=============================================================
// World Internal Functions
//=============================================================
//=============================================================

//=============================================================
// Flood fill in map
//=============================================================
func (w *world) floodFill(x, y int) {

}

//=============================================================
// Remove shadows from map on given position
//=============================================================
func (w *world) removeShadow(x, y int) {
	pos := w.width*x + y
	if pos < w.size && pos >= 0 {
		if w.pixels[pos]&0xFF == wShadow8 {
			r := uint32(float64(w.pixels[pos]>>24&0xFF) * 1.5)
			g := uint32(float64(w.pixels[pos]>>16&0xFF) * 1.5)
			b := uint32(float64(w.pixels[pos]>>8&0xFF) * 1.5)
			w.pixels[pos] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | wBackground8
			w.markChunkDirty(x, y)
		}
	}
}

//=============================================================
// Add shadows to map on given position
//=============================================================
func (w *world) addShadow(x, y int) {
	pos := w.width*x + y
	if pos < w.width*w.height && pos >= 0 {
		if w.pixels[pos]&0xFF != wShadow8 && w.pixels[pos]&0xFF != 0xFF {
			r := uint32(float64(w.pixels[pos]>>24&0xFF) / 1.5)
			g := uint32(float64(w.pixels[pos]>>16&0xFF) / 1.5)
			b := uint32(float64(w.pixels[pos]>>8&0xFF) / 1.5)
			w.pixels[pos] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | wShadow8
			w.markChunkDirty(x, y)
		}
	}
}

//=============================================================
// Mark chunk as dirty to rebuild it
//=============================================================
func (w *world) markChunkDirty(x, y int) {
	// Get all chunks in this area.
	for _, v := range w.qt.RetrieveIntersections(&Bounds{X: float64(x), Y: float64(y), Width: 3, Height: 3}) {
		if v.entity.getType() == entityChunk {
			v.entity.(*chunk).dirty = true
		}
	}
}

//=============================================================
// paint generated map
// everything has to be performed in a specific order.
//=============================================================
func (w *world) paintMap() {
	pcgGen := &pcg{}

	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			p := w.pixels[x*w.width+y]
			r := p >> 24 & 0xFF
			g := p >> 16 & 0xFF
			b := p >> 8 & 0xFF
			if r == 0xFF && g == 0x00 && b == 0x00 {
				v := w.coloring.getBackgroundSoft()
				// add some alpha to background
				v &= wBackground32
				w.pixels[x*w.width+y] = v
			} else if r == 0x00 && g == 0x00 && b == 0x00 {
				v := w.coloring.getBackground()
				w.pixels[x*w.width+y] = v
			}
		}
	}

	// Ladders
	color := w.coloring.getLadder()
	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			if y+1 < w.height && x+60 < w.width && x > 0 && y > 0 {
				before := w.pixels[(x-1)*w.width+y] & 0xFF
				point := w.pixels[x*w.width+y] & 0xFF
				after := w.pixels[(x+1)*w.width+y] & 0xFF
				above := w.pixels[x*w.width+y+1] & 0xFF
				long := w.pixels[(x+60)*w.width+y] & 0xFF

				if (above == wBackground8 || above == wShadow8) && point == 0xFF && before == 0xFF && (after == wBackground8 || after == wShadow8) && long == 0xFF {
					for i := 5; i < 18; i++ {
						if i == 5 || i == 17 {
							for n := 0; n < 500000; n++ {
								if y-n > 0 {
									if (w.pixels[(x+i)*w.width+y-n]&0xFF == wBackground8 || w.pixels[(x+i)*w.width+y-n]&0xFF == wShadow8) && w.pixels[(x+i)*w.width+y-n]&0xFF != wLadder8 {
										w.pixels[(x+i)*w.width+y-n] = color & wLadder32
										// Shadows
										if w.pixels[(x+i+1)*w.width+y-n-1]&0xFF != 0xFF {
											w.pixels[(x+i+1)*w.width+y-n-1] &= 0x555555FF & wLadder32
										}
									} else {
										break
									}
								}
							}

						}
						for n := 0; ; n += 5 {
							if y-n > 0 {
								if w.pixels[(x+i)*w.width+y-n]&0xFF == wBackground8 || w.pixels[(x+i)*w.width+y-n]&0xFF == wShadow8 {
									w.pixels[(x+i)*w.width+y-n] = color & wLadder32
									// Dont shadow above walls
									if w.pixels[(x+i+1)*w.width+y-n-1]&0xFF != 0xFF {
										w.pixels[(x+i+1)*w.width+y-n-1] &= 0x555555FF & wLadder32
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
				if y-5 > 0 && x+5 < w.width {
					below := w.pixels[x*w.width+y-1]
					right := w.pixels[(x+1)*w.width+y]
					if (below&0xFF == wShadow8 || below&0xFF == wBackground8) && point&0xFF == 0xFF {
						for i := 1; i < wShadowLength; i++ {
							p := w.pixels[(x+i)*w.width+y-i]
							if p&0xFF != wShadow8 && p&0xFF != 0xFF {
								r := uint32(float64(p>>24&0xFF) / 1.5)
								g := uint32(float64(p>>16&0xFF) / 1.5)
								b := uint32(float64(p>>8&0xFF) / 1.5)
								w.pixels[(x+i)*w.width+y-i] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | wShadow8&0xFF
							}
						}
					}
					if (right&0xFF == wShadow8 || right&0xFF == wBackground8) && point&0xFF == 0xFF {
						for i := 0; i < wShadowLength; i++ {
							p := w.pixels[(x+i)*w.width+y-i]
							if p&0xFF != wShadow8 && p&0xFF != 0xFF {
								r := uint32(float64(p>>24&0xFF) / 1.5)
								g := uint32(float64(p>>16&0xFF) / 1.5)
								b := uint32(float64(p>>8&0xFF) / 1.5)
								w.pixels[(x+i)*w.width+y-i] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | wShadow8&0xFF
							}
							p = w.pixels[(x+i)*w.width+y-i-1]
							if p&0xFF != wShadow8 && p&0xFF != 0xFF && i < 4 {
								r := uint32(float64(p>>24&0xFF) / 1.5)
								g := uint32(float64(p>>16&0xFF) / 1.5)
								b := uint32(float64(p>>8&0xFF) / 1.5)
								w.pixels[(x+i)*w.width+y-i-1] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | wShadow8&0xFF
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
	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			if y+1 < w.height && x+1 < w.width && x > 0 && y > 0 {
				before := w.pixels[(x-1)*w.width+y] & 0xFF
				point := w.pixels[x*w.width+y] & 0xFF
				after := w.pixels[(x+1)*w.width+y] & 0xFF
				above := w.pixels[x*w.width+y+1] & 0xFF
				below := w.pixels[x*w.width+y-1] & 0xFF

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
	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			if y+1 < w.height && x+1 < w.width && x > 0 && y > 0 {
				before := w.pixels[(x-1)*w.width+y] & 0xFF
				point := w.pixels[x*w.width+y] & 0xFF
				after := w.pixels[(x+1)*w.width+y] & 0xFF
				above := w.pixels[x*w.width+y+1] & 0xFF
				below := w.pixels[x*w.width+y-1] & 0xFF
				cornerRight := w.pixels[(x+1)*w.width+y-1] & 0xFF
				cornerLeft := w.pixels[(x-1)*w.width+y+1] & 0xFF
				cornerRight2 := w.pixels[(x+1)*w.width+y+1] & 0xFF
				cornerLeft2 := w.pixels[(x-1)*w.width+y-1] & 0xFF

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
	// for x := 0; x < w.width; x++ {
	// 	for y := 0; y < w.height; y++ {
	// 		p := w.pixels[x*w.width+y] & 0xFF

	// 		if p == wShadow8 || p == wBackground8 {
	// 			if rand.Float64() < 0.0001 {
	// 				pcgGen.GenerateBricks(x, y)
	// 			}
	// 		}
	// 	}
	// }

	// Background gfx
	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			if y+30 < w.height && x+1 < w.width && x > 0 && y > 0 {
				p := w.pixels[x*w.width+y] & 0xFF
				pp := w.pixels[x*w.width+y+1] & 0xFF
				up := w.pixels[x*w.width+y+30] & 0xFF
				upLow := w.pixels[x*w.width+y+5] & 0xFF
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
	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			if y+30 < w.height && x+1 < w.width && x > 0 && y > 0 {
				p := w.pixels[x*w.width+y] & 0xFF
				pp := w.pixels[x*w.width+y+1] & 0xFF
				upLow := w.pixels[x*w.width+y+5] & 0xFF
				if p == 0xFF && (upLow == wBackground8 || upLow == wShadow8) && (pp == wShadow8 || pp == wBackground8) {
					pcgGen.GenerateBottomAirIntake(x, y)
				}
			}
		}
	}
	// Background gfx
	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			if y+wDoorHeight < w.height && x+wDoorLen < w.width && x > 0 && y > 0 {
				p := w.pixels[x*w.width+y] & 0xFF
				pafter := w.pixels[(x+wDoorLen)*w.width+y+1] & 0xFF
				pbelow := w.pixels[(x+wDoorLen)*w.width+y-1] & 0xFF
				pbelowLong := w.pixels[x*w.width+y-wDoorHeight-55] & 0xFF
				pp := w.pixels[x*w.width+y+1] & 0xFF
				up := w.pixels[x*w.width+y+wDoorHeight] & 0xFF
				down := w.pixels[x*w.width+y-1] & 0xFF

				if p == 0xFF && pbelowLong == 0xFF && (down == wShadow8 || down == wBackground8) {
					pcgGen.GenerateLamp(x, y-1)
				}

				if pbelow != wBackground8 && pafter == wBackground8 && p == 0xFF && up == wBackground8 && pp == wBackground8 {
					// Check there is no ladder
					skip := false
					for i := 0; i < wDoorLen; i++ {
						if w.pixels[(x+i)*w.width+y+1]&0xFF == wLadder8 {
							skip = true
							break
						}
					}
					if !skip {
						pcgGen.GenerateDoor(x, y+1)
					}
				}
			}
		}
	}

}
