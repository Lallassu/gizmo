//=============================================================
// chunk.go
//-------------------------------------------------------------
// Part of the world. Handles its part of world pixels.
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"time"
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
	c.triangles = pixel.MakeTrianglesData(wPixelsPerChunk * wPixelsPerChunk * 6)
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
	c.bdt += dt
	if dt == -1 {
		c.build()
		return
	}
	if c.dirty && c.bdt > 0.2 {
		c.bdt = 0
		c.build()
	}
	//c.canvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(c.bounds.X+(c.bounds.Width)/2, c.bounds.Y+c.bounds.Height/2)))
	c.batch.Draw(global.gWin)
}

//=============================================================
// Rebuild/Build the chunk.
//=============================================================
func (c *chunk) build() {
	start := time.Now()
	i := 0
	r := 0.0
	g := 0.0
	b := 0.0
	a := 0.0
	for x := 0.0; x < c.bounds.Width; x++ {
		for y := 0.0; y < c.bounds.Height; y++ {
			p := global.gWorld.pixels[int(float64(global.gWorld.width)*(x+c.bounds.X)+(y+c.bounds.Y))]
			if p == 0 {
				continue
			}

			r = float64(p>>24&0xFF) / 255.0
			g = float64(p>>16&0xFF) / 255.0
			b = float64(p>>8&0xFF) / 255.0
			a = float64(p&0xFF) / 255.0

			px := (x + c.bounds.X)
			py := y + c.bounds.Y

			// TBD: Greedy algorithm to check for range of colors.

			(*c.triangles)[i].Position = pixel.Vec{px, py}
			(*c.triangles)[i+1].Position = pixel.Vec{px + 1, py}
			(*c.triangles)[i+2].Position = pixel.Vec{px + 1, py + 1}
			(*c.triangles)[i+3].Position = pixel.Vec{px, py}
			(*c.triangles)[i+4].Position = pixel.Vec{px, py + 1}
			(*c.triangles)[i+5].Position = pixel.Vec{px + 1, py + 1}
			(*c.triangles)[i].Color = pixel.RGBA{r, g, b, a}
			(*c.triangles)[i+1].Color = pixel.RGBA{r, g, b, a}
			(*c.triangles)[i+2].Color = pixel.RGBA{r, g, b, a}
			(*c.triangles)[i+3].Color = pixel.RGBA{r, g, b, a}
			(*c.triangles)[i+4].Color = pixel.RGBA{r, g, b, a}
			(*c.triangles)[i+5].Color = pixel.RGBA{r, g, b, a}
			i += 6
		}
	}
	c.batch.Dirty()
	elapsed := time.Since(start)
	Debug("Build took %s", elapsed)
	c.dirty = false
}

//func (c *chunk) buildFirst() {
//	return
//	start := time.Now()
//	model := imdraw.New(nil)
//	for x := 0.0; x < c.bounds.Width; x++ {
//		for y := 0.0; y < c.bounds.Height; y++ {
//			p := global.gWorld.pixels[int(float64(global.gWorld.width)*(x+c.bounds.X)+(y+c.bounds.Y))]
//			if p == 0 {
//				continue
//			}
//
//			model.Color = pixel.RGB(
//				float64(p>>24&0xFF)/255.0,
//				float64(p>>16&0xFF)/255.0,
//				float64(p>>8&0xFF)/255.0,
//			//).Mul(pixel.Alpha(float64(p&0xFF) / 255.0))
//			).Mul(pixel.Alpha(float64(p&0xFF) / 255.0))
//
//			model.Push(
//				pixel.V(float64(c.bounds.X+x*wPixelSize), float64(c.bounds.Y+y*wPixelSize)),
//				//pixel.V(float64(x*wPixelSize+wPixelSize), float64(y*wPixelSize+wPixelSize)),
//				pixel.V(float64(c.bounds.X+x*wPixelSize+wPixelSize), float64(c.bounds.Y+y*wPixelSize+wPixelSize)),
//			)
//			model.Rectangle(0)
//		}
//	}
//	//c.canvas.Clear(pixel.RGBA{0, 0, 0, 0})
//	elapsed := time.Since(start)
//	Debug("FIRST took %s", elapsed)
//	model.Draw(c.batch)
//	//c.batch.Dirty()
//	c.dirty = false
//}
