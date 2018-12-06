//=============================================================
// chunk.go
//-------------------------------------------------------------
// Part of the world. Handles its part of world pixels.
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	_ "math/rand"
	_ "time"
)

//=============================================================
// Chunk Structure
//=============================================================
type chunk struct {
	dirty bool
	//	canvas    *pixelgl.Canvas
	batch     *pixel.Batch
	triangles *pixel.TrianglesData
	bounds    *Bounds
	bdt       float64 // build dt
}

//=============================================================
// Impl. the Entity interface
//=============================================================
func (c *chunk) hit(x, y, vx, vy float64, power int) bool {
	return false
}

func (c *chunk) explode() {
}

func (c *chunk) move(x, y float64) {
}

func (c *chunk) getMass() float64 {
	return 1.0
}

func (c *chunk) getType() entityType {
	return entityChunk
}

func (c *chunk) getPosition() pixel.Vec {
	return pixel.Vec{c.bounds.X, c.bounds.Y}
}

func (c *chunk) setPosition(x, y float64) {
	c.bounds.X = x
	c.bounds.Y = y
}

func (c *chunk) getBounds() *Bounds {
	return c.bounds
}

//=============================================================
// Create a new chunk
//=============================================================
func (c *chunk) create(x, y float64) {
	//c.canvas = pixelgl.NewCanvas(pixel.R(0, 0, wPixelsPerChunk/2, wPixelsPerChunk))
	//	c.canvas = pixelgl.NewCanvas(pixel.R(0, 0, 1, 1))
	c.dirty = true
	c.triangles = pixel.MakeTrianglesData(500) //wPixelsPerChunk * wPixelsPerChunk * 6)
	c.batch = pixel.NewBatch(c.triangles, nil) //c.canvas)
	c.bounds = &Bounds{
		X:      x,
		Y:      y,
		Width:  wPixelsPerChunk,
		Height: wPixelsPerChunk,
		entity: Entity(c),
	}
}

//=============================================================
// Draw the chunk
//=============================================================
func (c *chunk) draw(dt float64) {
	//	c.bdt += dt
	if dt == -1 {
		c.build()
		return
	}
	if c.dirty { //&& c.bdt > 0.2 {
		//c.bdt = 0
		c.build()
	}
	//c.canvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(c.bounds.X+(c.bounds.Width)/2, c.bounds.Y+c.bounds.Height/2)))
	c.batch.Draw(global.gWin)
}

//=============================================================
// Rebuild/Build the chunk.
//=============================================================
func (c *chunk) build() {
	//start := time.Now()
	i := 0
	r := 0.0
	g := 0.0
	b := 0.0
	a := 0.0
	//skip := 0
	draw := 0

	// TBD: Check alpha!!!

	for x := 0.0; x < c.bounds.Width; x++ {
		for y := 0.0; y < c.bounds.Height; y++ {
			p := global.gWorld.pixels[int(float64(global.gWorld.width)*(x+c.bounds.X)+(y+c.bounds.Y))]
			if p == 0 {
				continue
			}

			// Check if already visisted
			if ((p & 0xFF) >> 7) == 0 {
				//skip++
				continue
			}

			r = float64(p>>24&0xFF) / 255.0
			g = float64(p>>16&0xFF) / 255.0
			b = float64(p>>8&0xFF) / 255.0
			a = float64(p&0xFF) / 255.0

			// Greedy algorithm to check for range of colors.
			// Use first bit in alpha to check for if it has been visited or not.
			// It's not being used anyway. Or at least for now :)

			// First check how far we can go with the same pixel color
			same_x := 1.0
			same_y := 1.0
			pos := 0
			p2 := uint32(0)
			r1 := 0.0
			g1 := 0.0
			b1 := 0.0
			// For each X, walk as long as possible towards Y
			for l := x + 1; l < c.bounds.Width; l++ {
				// Check color
				pos = int(float64(global.gWorld.width)*(l+c.bounds.X) + (y + c.bounds.Y))
				p2 = global.gWorld.pixels[pos]
				if p2 == 0 {
					break
				}
				r1 = float64(p2>>24&0xFF) / 255.0
				g1 = float64(p2>>16&0xFF) / 255.0
				b1 = float64(p2>>8&0xFF) / 255.0

				if r1 == r && g1 == g && b1 == b && ((p2&0xFF)>>7) == 1 {
					// Same color and not yet visited!
					global.gWorld.pixels[pos] &= 0xFFFFFF7F
					//Debug("AFTER: ", strconv.FormatInt(int64(global.gWorld.pixels[int(float64(global.gWorld.width)*(l+c.bounds.X)+(y+c.bounds.Y))]), 2))
					same_x++
					new_y := 1.0
					for k := y; k < c.bounds.Height; k++ {
						pos = int(float64(global.gWorld.width)*(l+c.bounds.X) + (k + c.bounds.Y))
						p2 = global.gWorld.pixels[pos]
						r1 = float64(p2>>24&0xFF) / 255.0
						g1 = float64(p2>>16&0xFF) / 255.0
						b1 = float64(p2>>8&0xFF) / 255.0

						if r1 == r && g1 == g && b1 == b && ((p2&0xFF)>>7) == 1 {
							global.gWorld.pixels[pos] &= 0xFFFFFF7F
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

			px := x + c.bounds.X
			py := y + c.bounds.Y

			draw++

			// Increase length of triangles if we need to draw more than we had before.
			if draw*6 >= len(*c.triangles) {
				c.triangles.SetLen(draw*6 + 100)
			}
			//	r = rand.Float64()

			// Size of triangle is given by how large greedy algorithm found out.
			(*c.triangles)[i].Position = pixel.Vec{px, py}
			(*c.triangles)[i+1].Position = pixel.Vec{px + same_x, py}
			(*c.triangles)[i+2].Position = pixel.Vec{px + same_x, py + same_y}
			(*c.triangles)[i+3].Position = pixel.Vec{px, py}
			(*c.triangles)[i+4].Position = pixel.Vec{px, py + same_y}
			(*c.triangles)[i+5].Position = pixel.Vec{px + same_x, py + same_y}
			(*c.triangles)[i].Color = pixel.RGBA{r, g, b, a}
			(*c.triangles)[i+1].Color = pixel.RGBA{r, g, b, a}
			(*c.triangles)[i+2].Color = pixel.RGBA{r, g, b, a}
			(*c.triangles)[i+3].Color = pixel.RGBA{r, g, b, a}
			(*c.triangles)[i+4].Color = pixel.RGBA{r, g, b, a}
			(*c.triangles)[i+5].Color = pixel.RGBA{r, g, b, a}
			i += 6
		}
	}

	// Reset the greedy bit
	for x := 0.0; x < c.bounds.Width; x++ {
		for y := 0.0; y < c.bounds.Height; y++ {
			global.gWorld.pixels[int(float64(global.gWorld.width)*(x+c.bounds.X)+(y+c.bounds.Y))] |= 0x00000080
		}
	}
	c.triangles.SetLen(draw * 6)
	c.batch.Dirty()
	// elapsed := time.Since(start)
	// Debug("Build took %s", elapsed, "SKIP:", skip, "Draw:", draw, "Total:", len(*c.triangles)/6, "Decr:", 100.0-((float64(draw)*6.0)/(float64(wPixelsPerChunk*wPixelsPerChunk)*6.0)*100.0), "%")
	c.dirty = false
}
