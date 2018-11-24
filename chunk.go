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
	bounds *Bounds
	bdt    float64 // build dt
}

//=============================================================
// Impl. the Entity interface
//=============================================================
func (c *chunk) hit(x, y, vx, vy float64) bool {
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

//=============================================================
// Create a new chunk
//=============================================================
func (c *chunk) create(x, y float64) {
	c.canvas = pixelgl.NewCanvas(pixel.R(0, 0, wPixelsPerChunk, wPixelsPerChunk))
	c.dirty = true
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
	c.canvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(c.bounds.X+(c.bounds.Width)/2, c.bounds.Y+c.bounds.Height/2)))
}

//=============================================================
// Get chunk position
//=============================================================
func (c *chunk) getPosition() pixel.Vec {
	return pixel.Vec{c.bounds.X, c.bounds.Y}
}

//=============================================================
// Set chunk position
//=============================================================
func (c *chunk) setPosition(x, y float64) {
	c.bounds.X = x
	c.bounds.Y = y
}

//=============================================================
// Get bounds
//=============================================================
func (c *chunk) getBounds() *Bounds {
	return c.bounds
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
			//).Mul(pixel.Alpha(float64(p&0xFF) / 255.0))
			).Mul(pixel.Alpha(float64(p&0xFF) / 255.0))

			model.Push(
				pixel.V(float64(x*wPixelSize), float64(y*wPixelSize)),
				//pixel.V(float64(x*wPixelSize+wPixelSize), float64(y*wPixelSize+wPixelSize)),
				pixel.V(float64(x*wPixelSize+wPixelSize), float64(y*wPixelSize+wPixelSize)),
			)
			model.Rectangle(0)
		}
	}
	c.canvas.Clear(pixel.RGBA{0, 0, 0, 0})
	model.Draw(c.canvas)
	c.dirty = false
}
