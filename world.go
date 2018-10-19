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
)

//=============================================================
// World Structure
//=============================================================
type world struct {
	width      int
	height     int
	coloring   *mapColor
	qt         *Quadtree
	pixels     []uint32
	currentMap mapType
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

	w.pixels = make([]uint32, w.width*w.height)

	Debug("Generating world", w.width, w.height)
	g := generator{}
	pixels := g.NewWorld(w.width, w.height)

	// Ad dall pixels as red before coloring
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
	for _, v := range w.qt.RetrieveIntersections(Bounds{X: 0, Y: 0, Width: float64(w.width), Height: float64(w.height)}) {
		v.entity.draw(0)
	}
	Debug("Tree Size:", w.qt.Total)
	Debug("World generation complete.")
}

//=============================================================
// Add object to world (QT)
//=============================================================
func (w *world) AddObject(x, y float64, obj Entity) {

}

//=============================================================
// Remove object from world (QT)
//=============================================================
func (w *world) RemoveObject(obj Entity) {

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
	return 1
}

//=============================================================
// Draw
//=============================================================
func (w *world) Draw(dt float64) {
	// Draw those around camera position only.
	for _, v := range w.qt.RetrieveIntersections(Bounds{X: global.gCamera.pos.X, Y: global.gCamera.pos.Y, Width: wViewMax, Height: wViewMax}) {
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
		//	         w.MarkChunkDirty(x, y)
	}
}

//=============================================================
// Remove a pixel from the world map
//=============================================================
func (w *world) RemovePixel(x, y int) {
	pos := w.width*x + y
	if pos < w.width*w.height && pos >= 0 {
		if w.pixels[pos]&0xFF == wStaticColor8 ||
			w.pixels[pos]&0xFF == wBackground8 || w.pixels[pos]&0xFF == wShadow8 {
			return
		}

		// Remove shadow
		for i := 0; i < 25; i++ {
			pos2 := (x+i)*w.width + y - i
			if pos2 < w.width*w.height && pos2 >= 0 {
				if w.pixels[pos2]&0xFF == wShadow8 {
					//w.RemoveShadow(x+i, y-i)
				}
			}
		}

		// Add shadow to visible pixel
		for i := 0; i < 5; i++ {
			pos2 := (x+i)*w.width + y - i
			if pos2 < w.width*w.height && pos2 >= 0 {
				if w.pixels[pos2]&0xFF == wBackground8 {
					//w.AddShadow(x+i, y-i)
				}
			}
		}

		// Particle
		// if w.Exists(float64(x), float64(y)) {
		// 	p := Particle{
		// 		X:     float64(x),
		// 		Y:     float64(y),
		// 		Size:  1,
		// 		Color: w.pixels[pos],
		// 	}
		// 	w.pts.NewParticle(p)
		// }

		// Set bg pixel.
		if w.pixels[pos] != 0 {
			v := w.coloring.getBackground()
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
	ff := make([]pixel.Vec, 10)
	for rx := x - power; rx <= x+power; rx++ {
		vx := (rx - x) * (rx - x)
		for ry := y - power; ry <= y+power; ry++ {
			if ry < 0 {
				continue
			}
			val := (ry-y)*(ry-y) + vx
			if val <= pow {
				w.RemovePixel(rx, ry)
				//w.ObjectHit(float64(rx), float64(ry))
			} else {
				ff = append(ff, pixel.Vec{X: float64(rx), Y: float64(ry)})
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
// Mark chunk as dirty to rebuild it
//=============================================================
func (w *world) markChunkDirty(x, y int) {
	// Get all chunks in this area.
	for _, v := range w.qt.RetrieveIntersections(Bounds{X: float64(x), Y: float64(y), Width: 3, Height: 3}) {
		if v.entity.getType() == entityChunk {
			v.entity.(*chunk).dirty = true
		}
	}
}

//=============================================================
// paint generated map
//=============================================================
func (w *world) paintMap() {
	color := w.coloring.getBackground()
	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			p := w.pixels[x*w.width+y]
			// Sides
			if x+1 < w.width {
				p2 := w.pixels[(x+1)*w.width+y]
				if p == 0 && p2 != 0 {
					for i := 0; i < wBorderSize; i++ {
						if i < wStaticBorderSize {
							w.pixels[(x+i)*w.width+y] = color & wStaticColor8
						} else {
							w.pixels[(x+i)*w.width+y] = color
						}
					}
				}
				if p != 0 && p2 == 0 {
					for i := 0; i < wBorderSize; i++ {
						if x-i > 0 && x-i < w.width {
							if i < wStaticBorderSize {
								w.pixels[(x-i)*w.width+y] = color & wStaticColor8
							} else {
								w.pixels[(x-i)*w.width+y] = color
							}
						}
					}
				}
			}
			// Top/Bottom
			if y+1 < w.height {
				p2 := w.pixels[x*w.width+y+1]
				if p == 0 && p2 != 0 {
					for i := 0; i < wBorderSize; i++ {
						if i < wStaticBorderSize {
							w.pixels[x*w.width+y+i] = color & wStaticColor8
						} else {
							w.pixels[x*w.width+y+i] = color
						}
					}
				}
				if p != 0 && p2 == 0 {
					for i := 0; i < wBorderSize; i++ {
						if y-i > 0 && y-i < w.height {
							if i < wStaticBorderSize {
								w.pixels[x*w.width+y-i] = color & wStaticColor8
							} else {
								w.pixels[x*w.width+y-i] = color
							}
						}
					}
				}
			}
		}
	}

	// Corners
	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			p := w.pixels[x*w.width+y]
			// Corners
			if y+1 < w.height && x+1 < w.width && x > 0 && y > 0 {
				p2 := w.pixels[(x-1)*w.width+y+1]
				p3 := w.pixels[(x-1)*w.width+y]
				p4 := w.pixels[x*w.width+y+1]
				p5 := w.pixels[(x-1)*w.width+y-1]
				p6 := w.pixels[x*w.width+y-1]
				p7 := w.pixels[(x+1)*w.width+y-1]
				p8 := w.pixels[(x+1)*w.width+y+1]

				//   x--
				//   |
				if p != 0 && p2 == 0 && p3 != 0 && p4 != 0 {
					for i := 0; i < wBorderSize; i++ {
						for j := 0; j < wBorderSize; j++ {
							if x+j < w.height && y-i > 0 {
								w.pixels[(x+j)*w.width+y-i+j] = color
							}
						}
					}
				}
				//   |
				//   x--
				if p != 0 && p2 != 0 && p5 == 0 && p6 != 0 && p3 == 0 {
					for i := 0; i <= wBorderSize; i++ {
						for j := 0; j < wBorderSize; j++ {
							if x+j < w.height && y+i < w.width {
								w.pixels[(x+j)*w.width+y+i-j] = color
							}
						}
					}
				}
				//   |
				// --x
				if p != 0 && p2 != 0 && p8 != 0 && p5 != 0 && p7 == 0 {
					for i := 0; i < wBorderSize; i++ {
						for j := 0; j < wBorderSize; j++ {
							if x-j > 0 && y+i < w.width {
								w.pixels[(x-j)*w.width+y+i-j] = color
							}
						}
					}
				}
				// --x
				//   |
				if p != 0 && p8 == 0 && p2 != 0 && p7 != 0 {
					for i := 0; i < wBorderSize; i++ {
						for j := 0; j < wBorderSize; j++ {
							if x-j > 0 && y-i > 0 {
								w.pixels[(x-j)*w.width+y-i+j] = color
							}
						}
					}
				}
			}
		}
	}

	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			p := w.pixels[x*w.width+y]
			r := p >> 24 & 0xFF
			g := p >> 16 & 0xFF
			b := p >> 8 & 0xFF
			if r == 0xFF && g == 0x00 && b == 0x00 {
				// Keep alpha (shadows)
				v := w.coloring.getBackground()
				// add some alpha to background
				v &= wBackground32
				w.pixels[x*w.width+y] = v
			} else if r == 0x00 && g == 0x00 && b == 0x00 {
				//v := uint32(0x2c5557FF)
				v := w.coloring.getBackground()
				// add some alpha to background
				w.pixels[x*w.width+y] = v
			}
		}
	}

	// Ladders
	color = w.coloring.getLadder()
	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			if y+1 < w.height && x+1 < w.width && x > 0 && y > 0 {
				before := w.pixels[(x-1)*w.width+y] & 0xFF
				point := w.pixels[x*w.width+y] & 0xFF
				after := w.pixels[(x+1)*w.width+y] & 0xFF
				above := w.pixels[x*w.width+y+1] & 0xFF
				long := uint32(0)
				if x+23 < w.width {
					long = w.pixels[(x+23)*w.width+y] & 0xFF
				}
				if above == wBackground8 && point == 0xFF && before == 0xFF && after == wBackground8 && long == 0xFF {
					for i := 0; i < 18; i++ {
						if i == 5 || i == 17 {
							for n := 0; n < 500000; n++ {
								if y-n > 0 {
									if w.pixels[(x+i)*w.width+y-n]&0xFF == wBackground8 && w.pixels[(x+i)*w.width+y-n]&0xFF != wLadder8 {
										w.pixels[(x+i)*w.width+y-n] = (color & wLadder32)
										// Shadows
										if w.pixels[(x+i+1)*w.width+y-n-1]&0xFF != 0xFF {
											w.pixels[(x+i+1)*w.width+y-n-1] &= (color & wLadder32)
										}
									} else {
										break
									}
								}
							}

						}
						for n := 0; ; n += 5 {
							if y-n > 0 {
								if w.pixels[(x+i)*w.width+y-n]&0xFF == wBackground8 {
									w.pixels[(x+i)*w.width+y-n] = (color & wLadder32)
									// Dont shadow above walls
									if w.pixels[(x+i+1)*w.width+y-n-1]&0xFF != 0xFF {
										w.pixels[(x+i+1)*w.width+y-n-1] &= (color & wLadder32)
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
						for i := 1; i < 5; i++ {
							p := w.pixels[(x+i)*w.width+y-i]
							if p&0xFF != wShadow8 && p&0xFF != 0xFF {
								r := uint32(float64((p >> 24 & 0xFF)) / 1.5)
								g := uint32(float64((p >> 16 & 0xFF)) / 1.5)
								b := uint32(float64((p >> 8 & 0xFF)) / 1.5)
								w.pixels[(x+i)*w.width+y-i] = (r & 0xFF << 24) | (g & 0xFF << 16) | (b & 0xFF << 8) | wShadow8&0xFF
							}
						}
					}
					if (right&0xFF == wShadow8 || right&0xFF == wBackground8) && point&0xFF == 0xFF {
						for i := 0; i < 5; i++ {
							p := w.pixels[(x+i)*w.width+y-i]
							if p&0xFF != wShadow8 && p&0xFF != 0xFF {
								r := uint32(float64((p >> 24 & 0xFF)) / 1.5)
								g := uint32(float64((p >> 16 & 0xFF)) / 1.5)
								b := uint32(float64((p >> 8 & 0xFF)) / 1.5)
								w.pixels[(x+i)*w.width+y-i] = (r & 0xFF << 24) | (g & 0xFF << 16) | (b & 0xFF << 8) | wShadow8&0xFF
							}
							p = w.pixels[(x+i)*w.width+y-i-1]
							if p&0xFF != wShadow8 && p&0xFF != 0xFF && i < 4 {
								r := uint32(float64((p >> 24 & 0xFF)) / 1.5)
								g := uint32(float64((p >> 16 & 0xFF)) / 1.5)
								b := uint32(float64((p >> 8 & 0xFF)) / 1.5)
								w.pixels[(x+i)*w.width+y-i-1] = (r & 0xFF << 24) | (g & 0xFF << 16) | (b & 0xFF << 8) | wShadow8&0xFF
							}
						}
					}
				}
			}
		}
	}
}
