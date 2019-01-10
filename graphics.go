//=============================================================
// graphics.go
//-------------------------------------------------------------
// Graphics for items and mobs
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image"
)

type graphics struct {
	sheetFile   string
	walkFrames  []int
	jumpFrames  []int
	climbFrames []int
	shootFrames []int
	idleFrames  []int
	frameWidth  float64
	frameHeight float64
	size        float64
	currentAnim animationType
	animCounter float64
	animRate    float64
	batches     map[int]*pixel.Batch
	triangles   map[int]*pixel.TrianglesData
	frames      map[int][]uint32
	img         image.Image
	canvas      *pixelgl.Canvas
	animated    bool
	scalexy     float64
}

//=============================================================
// Create graphics
//=============================================================
func (gfx *graphics) createGfx(x, y float64) {
	gfx.frames = make(map[int][]uint32)
	gfx.batches = make(map[int]*pixel.Batch)
	gfx.triangles = make(map[int]*pixel.TrianglesData)

	if gfx.scalexy == 0 {
		gfx.scalexy = 1.0
	}
	gfx.animRate = 0.1
	gfx.currentAnim = animIdle
	fullWidth := 0.0

	gfx.img, fullWidth, gfx.frameHeight, gfx.size = loadTexture(gfx.sheetFile)

	if !gfx.animated {
		gfx.frameWidth = fullWidth
		if gfx.frameWidth == gfx.frameHeight {
			gfx.size += 1
		}
	}

	f := 0
	for w := 0.0; w < fullWidth; w += gfx.frameWidth {
		gfx.frames[f] = make([]uint32, int(gfx.size)*int(gfx.size))
		for x := 0.0; x <= gfx.frameWidth; x++ {
			for y := 0.0; y <= gfx.frameHeight; y++ {
				r, g, b, a := gfx.img.At(int(w+x), int(gfx.frameHeight-y)).RGBA()
				if r == 0 && g == 0 && b == 0 && a == 0 {
					continue
				}
				gfx.frames[f][int(x*gfx.size+y-1)] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | a&0xFF
			}
		}
		gfx.triangles[f] = pixel.MakeTrianglesData(100)
		gfx.batches[f] = pixel.NewBatch(gfx.triangles[f], nil)
		f++
	}

	gfx.canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(gfx.frameWidth), float64(gfx.frameHeight)))

	// Build all frames
	gfx.buildFrames()

}

//=============================================================
// Build each frame
//=============================================================
func (gfx *graphics) buildFrames() {
	v := 0
	rc := uint32(0)
	gc := uint32(0)
	bc := uint32(0)
	p2 := uint32(0)
	r1 := uint32(0)
	g1 := uint32(0)
	b1 := uint32(0)
	draw := 0
	same_x := 1.0
	same_y := 1.0
	pos := 0

	// Build batch for each frame.
	for i := 0; i < len(gfx.frames); i++ {
		for x := 0.0; x < float64(gfx.frameWidth); x++ {
			for y := 0.0; y < float64(gfx.frameHeight); y++ {
				p := gfx.frames[i][int(x*gfx.size+y)]
				if p == 0 || p&0xFF>>7 == 0 {
					continue
				}
				rc = p >> 24 & 0xFF
				gc = p >> 16 & 0xFF
				bc = p >> 8 & 0xFF
				same_x = 1.0
				same_y = 1.0

				for l := x + 1; l < gfx.frameWidth; l++ {
					// Check color
					pos = int(l*gfx.size + y)
					p2 = gfx.frames[i][pos]
					if p2 == 0 {
						break
					}
					r1 = p2 >> 24 & 0xFF
					g1 = p2 >> 16 & 0xFF
					b1 = p2 >> 8 & 0xFF

					if r1 == rc && g1 == gc && b1 == bc && ((p2&0xFF)>>7) == 1 {
						// Same color and not yet visited!
						gfx.frames[i][pos] &= 0xFFFFFF7F
						same_x++
						new_y := 1.0
						for k := y; k < gfx.frameHeight; k++ {
							pos = int(l*gfx.size + k)
							p2 = gfx.frames[i][pos]
							if p2 == 0 {
								break
							}
							r1 = p2 >> 24 & 0xFF
							g1 = p2 >> 16 & 0xFF
							b1 = p2 >> 8 & 0xFF

							if r1 == rc && g1 == gc && b1 == bc && ((p2&0xFF)>>7) == 1 {
								gfx.frames[i][pos] &= 0xFFFFFF7F
								new_y++
							} else {
								break
							}
						}
						if new_y < same_y {
							break
						} else {
							same_y = new_y
						}
					} else {
						break
					}
				}

				draw++

				// Convert to decimal
				r := float64(p>>24&0xFF) / 255.0
				g := float64(p>>16&0xFF) / 255.0
				b := float64(p>>8&0xFF) / 255.0
				a := float64(p&0xFF) / 255.0

				// Increase length of triangles if we need to draw more than we had before.
				if draw*6 >= len(*gfx.triangles[i]) {
					gfx.triangles[i].SetLen(draw*6 + 10)
				}

				// Size of triangle is given by how large the greedy algorithm found out.
				(*gfx.triangles[i])[v].Position = pixel.Vec{x, y}
				(*gfx.triangles[i])[v+1].Position = pixel.Vec{x + same_x, y}
				(*gfx.triangles[i])[v+2].Position = pixel.Vec{x + same_x, y + same_y}
				(*gfx.triangles[i])[v+3].Position = pixel.Vec{x, y}
				(*gfx.triangles[i])[v+4].Position = pixel.Vec{x, y + same_y}
				(*gfx.triangles[i])[v+5].Position = pixel.Vec{x + same_x, y + same_y}

				//(*gfx.triangles[i])[v].Picture = pixel.Vec{x, y}
				//(*gfx.triangles[i])[v+1].Picture = pixel.Vec{x + same_x, y}
				//(*gfx.triangles[i])[v+2].Picture = pixel.Vec{x + same_x, y + same_y}
				//(*gfx.triangles[i])[v+3].Picture = pixel.Vec{x, y}
				//(*gfx.triangles[i])[v+4].Picture = pixel.Vec{x, y + same_y}
				//(*gfx.triangles[i])[v+5].Picture = pixel.Vec{x + same_x, y + same_y}
				//(*gfx.triangles[i])[v].Intensity = 0.5
				//(*gfx.triangles[i])[v+1].Intensity = 0.5
				//(*gfx.triangles[i])[v+2].Intensity = 0.5
				//(*gfx.triangles[i])[v+3].Intensity = 0.5
				//(*gfx.triangles[i])[v+4].Intensity = 0.5
				//(*gfx.triangles[i])[v+5].Intensity = 0.5
				for n := 0; n < 6; n++ {
					(*gfx.triangles[i])[v+n].Color = pixel.RGBA{r, g, b, a}
				}

				v += 6

			}
		}
		// Reset the greedy bit

		for x := 0.0; x < gfx.frameWidth; x++ {
			for y := 0.0; y < gfx.frameHeight; y++ {
				pos = int(x*gfx.size + y)
				if gfx.frames[i][pos] != 0 {
					gfx.frames[i][pos] |= 0x00000080
				}
			}
		}
		gfx.triangles[i].SetLen(draw * 6)
		gfx.batches[i].Dirty()
	}
}

//=============================================================
//
//=============================================================
func (gfx *graphics) hitGfx(lx, ly int, gx, gy, vx, vy float64, power int, blood bool) {
	if global.gRand.rand() < 1 {
		global.gParticleEngine.effectBlood(gx, gy, vx, vy, 1)
	}

	pow := power * power
	for rx := lx - power; rx <= lx+power; rx++ {
		xx := (rx - lx) * (rx - lx)
		for ry := ly - power; ry <= ly+power; ry++ {
			if ry < 0 {
				continue
			}
			val := (ry-ly)*(ry-ly) + xx
			if val < pow {
				for i := 0; i < len(gfx.frames); i++ {
					pos := int(gfx.size)*rx + ry
					if pos >= 0 && pos < int(gfx.size*gfx.size) {
						if gfx.frames[i][pos] != 0 {
							if blood {
								// Don't color eyes, assume white.
								p := gfx.frames[i][pos]
								if !(p>>24&0xFF == 0xFF && p>>16&0xFF == 0xFF && p>>8&0xFF == 0xFF) {
									// Blood
									r := 175 + global.gRand.rand()*5
									g := 10 + global.gRand.rand()*2
									b := 10 + global.gRand.rand()*2
									a := global.gRand.rand() * 255
									gfx.frames[i][pos] = uint32(r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | a&0xFF)
								}
							}
						}
					}
				}
			}
		}
	}

	gfx.buildFrames()
}

//=============================================================
//
//=============================================================
func (gfx *graphics) explodeGfx(gx, gy float64, blood bool) {
	size := gfx.scalexy
	if size < 0.5 {
		size = 0.5
	}

	for i := 0; i < len(gfx.frames); i++ {
		for x := 0.0; x < gfx.frameWidth; x++ {
			for y := 0.0; y < gfx.frameHeight; y++ {
				pos := int(gfx.size*x + y)
				if gfx.frames[i][pos] != 0 {
					// Remove parts (Don't create every particle)
					if global.gRand.randFloat() < 0.2 {

						if blood {
							global.gParticleEngine.effectBlood(gx+float64(x), gy+float64(y), float64(5-global.gRand.rand()), float64(5-global.gRand.rand()), global.gRand.rand()/10)
						}
						global.gParticleEngine.newParticle(
							particle{
								x:           gx + float64(x)*gfx.scalexy,
								y:           gy + float64(y)*gfx.scalexy,
								size:        size,
								restitution: -0.1 - global.gRand.randFloat()/4,
								life:        wParticleDefaultLife,
								fx:          float64(15 - global.gRand.rand()),
								fy:          float64(15 - global.gRand.rand()),
								vx:          float64(5 - global.gRand.rand()),
								vy:          float64(5 - global.gRand.rand()),
								mass:        1,
								pType:       particleRegular,
								color:       gfx.frames[i][pos],
								static:      true,
							})
					}
					gfx.frames[i][pos] = 0
				}
			}
		}
	}
}
