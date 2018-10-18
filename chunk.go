//=============================================================
// chunk.go
//-------------------------------------------------------------
// Part of the world. Handles its part of world pixels.
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	_ "github.com/faiface/pixel/pixelgl"
)

//=============================================================
// Chunk Structure
//=============================================================
type chunk struct {
	dirty  bool
	model  *imdraw.IMDraw
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
	c.model = imdraw.New(nil)
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
	c.model.Draw(global.gWin)
}

//=============================================================
// Rebuild/Build the chunk.
//=============================================================
func (c *chunk) build() {
	model := imdraw.New(nil)
	for x := c.bounds.X; x < c.bounds.X+c.bounds.Width; x++ {
		for y := c.bounds.Y; y < c.bounds.Y+c.bounds.Height; y++ {
			p := global.gWorld.pixels[int(float64(global.gWorld.width)*x+y)]
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
	c.model = model
	c.dirty = false
}
