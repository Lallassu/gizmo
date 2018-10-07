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
type Chunk struct {
	dirty  bool
	pixels []*uint32
	x      float64
	y      float64
	width  float64
	height float64
	model  *imdraw.IMDraw
}

//=============================================================
// Create a new chunk
//=============================================================
func (c *Chunk) create() {
	c.draw_model = imdraw.New(nil)
	c.dirty = true
}

//=============================================================
// Draw the chunk
//=============================================================
func (c *Chunk) draw(win *pixelgl.Window) {
	c.model.Draw(win)
}

//=============================================================
// Rebuild/Build the chunk.
//=============================================================
func (c *Chunk) build() {
	model := imdraw.New(nil)

	for x := 0; x < wPixelsPerChunk; x++ {
		for y := 0; y < wPixelsPerChunk; y++ {
			p := *c.Pixels[wPixelsPerChunk*x+y]
			if p == 0 {
				continue
			}

			model.Color = pixel.RGB(
				float64(p>>24&0xFF)/255.0,
				float64(p>>16&0xFF)/255.0,
				float64(p>>8&0xFF)/255.0,
			).Mul(pixel.Alpha(float64(p&0xFF) / 255.0))

			model.Push(
				pixel.V(float64(c.Posx+x*wPixelSize), float64(c.Posy+y*wPixelSize)),
				pixel.V(float64(c.Posx+x*wPixelSize+wPixelSize), float64(c.Posy+y*wPixelSize+wPixelSize)),
			)
			model.Rectangle(0)
		}
	}
	c.model = model
	c.dirty = false
}
