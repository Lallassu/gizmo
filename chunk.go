//=============================================================
// chunk.go
//-------------------------------------------------------------
// Part of the world. Handles its part of world pixels.
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

//=============================================================
// Chunk Structure
//=============================================================
type chunk struct {
	dirty  bool
	canvas *pixelgl.Canvas
	bounds Bounds
}

//=============================================================
// Support the Entity interface
//=============================================================
func (c *chunk) hit(x, y float64) bool {
	return false
}

func (c *chunk) explode() {
}

func (c *chunk) getMass() float64 {
	return 1.0
}

func (c *chunk) getType() entityType {
	return entityChunk
}

//=============================================================
// Create a new chunk
//=============================================================
func (c *chunk) create(x, y float64) {
	c.canvas = pixelgl.NewCanvas(pixel.R(0, 0, wPixelsPerChunk, wPixelsPerChunk))
	c.dirty = true
	c.bounds = Bounds{
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
	if c.dirty {
		c.build()
	}
	c.canvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(c.bounds.X, c.bounds.Y)))
}

//=============================================================
// Rebuild/Build the chunk.
//=============================================================
func (c *chunk) build() {
	model := imdraw.New(nil)
	for x := 0.0; x < c.bounds.Width; x++ {
		for y := 0.0; y < c.bounds.Height; y++ {
			p := global.gWorld.pixels[int(float64(global.gWorld.width)*(x+c.bounds.X)+(y+c.bounds.Y))]
			if p == 0 {
				continue
			}

			model.Color = pixel.RGB(
				float64(p>>24&0xFF)/255.0,
				float64(p>>16&0xFF)/255.0,
				float64(p>>8&0xFF)/255.0,
			).Mul(pixel.Alpha(float64(p&0xFF) / 255.0))

			model.Push(
				pixel.V(float64(x*wPixelSize), float64(y*wPixelSize)),
				pixel.V(float64(x*wPixelSize+wPixelSize), float64(y*wPixelSize+wPixelSize)),
			)
			model.Rectangle(0)
		}
	}
	model.Draw(c.canvas)
	c.dirty = false
}
