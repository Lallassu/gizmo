//=============================================================
// chunk.go
//-------------------------------------------------------------
// Part of the world. Handles its part of world pixels.
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

//=============================================================
// Chunk Structure
//=============================================================
type chunk struct {
	dirty     bool
	batch     *pixel.Batch
	triangles *pixel.TrianglesData
	bounds    *Bounds
	cType     chunkType
	sprite    *pixel.Sprite
}

//=============================================================
// Impl. the Entity interface
//=============================================================
func (c *chunk) hit(x, y, vx, vy float64, power int) {
}

func (c *chunk) getPosition() pixel.Vec {
	return pixel.Vec{X: c.bounds.X, Y: c.bounds.Y}
}

//=============================================================
// Create a new chunk
//=============================================================
func (c *chunk) create(x, y float64, pixels int) {
	c.dirty = true
	c.triangles = pixel.MakeTrianglesData(400) // Init with some suitable value
	c.batch = pixel.NewBatch(c.triangles, nil)
	c.bounds = &Bounds{
		X:      x,
		Y:      y,
		Width:  float64(pixels),
		Height: float64(pixels),
		entity: entity(c),
	}
}

//=============================================================
// Draw the chunk
//=============================================================
func (c *chunk) draw(dt, elapsed float64) {
	if c.cType == fgChunk {
		if c.dirty {
			c.build()
		}
		c.batch.Draw(global.gWin)
	} else {
		c.sprite.Draw(global.gWin, pixel.IM.Moved(pixel.V(c.bounds.X+c.bounds.Width/2, c.bounds.Y+c.bounds.Height/2)))
	}
}

//=============================================================
// Rebuild/Build the chunk.
//=============================================================
func (c *chunk) build() {
	//	start := time.Now()
	i := 0
	rc := uint32(0)
	gc := uint32(0)
	bc := uint32(0)
	p2 := uint32(0)
	r1 := uint32(0)
	g1 := uint32(0)
	b1 := uint32(0)
	draw := 0
	sameX := 1.0
	sameY := 1.0
	pos := 0
	px := 0.0
	py := 0.0
	xpos := 0.0

	for x := 0.0; x < c.bounds.Width; x++ {
		for y := 0.0; y < c.bounds.Height; y++ {
			p := global.gWorld.pixels[int(float64(global.gWorld.width)*(x+c.bounds.X)+(y+c.bounds.Y))]
			// Skip visisted or empty
			if p == 0 || p&0xFF>>7 == 0 {
				continue
			}
			if p&0xFF != wBackground8 && c.cType == bgChunk {
				continue
			}

			if p&0xFF == wBackground8 && c.cType == fgChunk {
				continue
			}
			rc = p >> 24 & 0xFF
			gc = p >> 16 & 0xFF
			bc = p >> 8 & 0xFF
			sameX = 1.0
			sameY = 1.0

			// Greedy algorithm to check for range of colors.
			// Use first bit in alpha to check for if it has been visited or not.
			// It's not being used anyway. Or at least for now :)

			// First check how far we can go with the same pixel color
			// For each X, walk as long as possible towards Y
			for l := x + 1; l < c.bounds.Width; l++ {
				// Check color
				xpos = float64(global.gWorld.width) * (l + c.bounds.X)
				pos = int(xpos + (y + c.bounds.Y))
				p2 = global.gWorld.pixels[pos]
				if p2 == 0 {
					break
				}
				r1 = p2 >> 24 & 0xFF
				g1 = p2 >> 16 & 0xFF
				b1 = p2 >> 8 & 0xFF

				if r1 == rc && g1 == gc && b1 == bc && ((p2&0xFF)>>7) == 1 {
					if p2&0xFF != wBackground8 && c.cType == bgChunk {
						break
					}
					// Same color and not yet visited!
					global.gWorld.pixels[pos] &= 0xFFFFFF7F
					sameX++
					newY := 1.0
					for k := y; k < c.bounds.Height; k++ {
						pos = int(xpos + (k + c.bounds.Y))
						p2 = global.gWorld.pixels[pos]
						r1 = p2 >> 24 & 0xFF
						g1 = p2 >> 16 & 0xFF
						b1 = p2 >> 8 & 0xFF

						if r1 == rc && g1 == gc && b1 == bc && ((p2&0xFF)>>7) == 1 {
							if p2&0xFF != wBackground8 && c.cType == bgChunk {
								break
							}
							global.gWorld.pixels[pos] &= 0xFFFFFF7F
							newY++
						} else {
							break
						}
					}
					if newY < sameY {
						break
					} else {
						sameY = newY
					}
				} else {
					break
				}
			}

			px = x + c.bounds.X
			py = y + c.bounds.Y

			draw++

			// Convert to decimal
			r := float64(p>>24&0xFF) / 255.0
			g := float64(p>>16&0xFF) / 255.0
			b := float64(p>>8&0xFF) / 255.0
			a := float64(p&0xFF) / 255.0

			// Increase length of triangles if we need to draw more than we had before.
			// Add a buffer so we can skip a few increments.
			if draw*6 >= len(*c.triangles) {
				c.triangles.SetLen(draw*6 + 60)
			}

			// Size of triangle is given by how large the greedy algorithm found out.
			(*c.triangles)[i].Position = pixel.Vec{X: px, Y: py}
			(*c.triangles)[i+1].Position = pixel.Vec{X: px + sameX, Y: py}
			(*c.triangles)[i+2].Position = pixel.Vec{X: px + sameX, Y: py + sameY}
			(*c.triangles)[i+3].Position = pixel.Vec{X: px, Y: py}
			(*c.triangles)[i+4].Position = pixel.Vec{X: px, Y: py + sameY}
			(*c.triangles)[i+5].Position = pixel.Vec{X: px + sameX, Y: py + sameY}
			for n := 0; n < 6; n++ {
				(*c.triangles)[i+n].Color = pixel.RGBA{R: r, G: g, B: b, A: a}
			}

			i += 6
		}
	}

	// Reset the greedy bit
	for x := 0.0; x < c.bounds.Width; x++ {
		for y := 0.0; y < c.bounds.Height; y++ {
			global.gWorld.pixels[int(float64(global.gWorld.width)*(x+c.bounds.X)+(y+c.bounds.Y))] |= 0x00000080
		}
	}
	//	elapsed := time.Since(start)
	//	Debug("Build took %s", elapsed, "SKIP:", skip, "Draw:", draw, "Total:", len(*c.triangles)/6, "Decr:", 100.0-((float64(draw)*6.0)/(float64(wPixelsPerChunk*wPixelsPerChunk)*6.0)*100.0), "%")
	c.triangles.SetLen(draw * 6)
	c.batch.Dirty()
	c.dirty = false

	// If background we build our sprite
	if c.cType == bgChunk {
		canvas := pixelgl.NewCanvas(pixel.R(0, 0, float64(global.gWorld.width), float64(global.gWorld.height)))
		c.batch.Draw(canvas)
		c.sprite = pixel.NewSprite(canvas, pixel.R(c.bounds.X, c.bounds.Y, c.bounds.X+c.bounds.Width, c.bounds.Y+c.bounds.Height))
		c.triangles.SetLen(0)
		c.batch = nil
	}
}
