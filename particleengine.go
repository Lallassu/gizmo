//=============================================================
// particleengine.go
//-------------------------------------------------------------
// Particles of different kinds.
//=============================================================
package main

import (
	"github.com/golang-collections/collections/stack"
)

type particleEngine struct {
	canvas   *pixelgl.Canvas
	active   *stack.Stack
	inactive *stack.Stack
}

//=============================================================
// Draw the canvas
//=============================================================
func (p *particleEngine) create() {
	p.canvas = pixelgl.NewCanvas(pixel.R(0, 0, global.gWindowHeight, global.gWindowWidth))
	p.active = stack.New()
	p.inactive = stack.New()

	for i := 0; i < wParticlesMax; i++ {
		p.inactive.push(particle{active: false})
	}

}

//=============================================================
// Get new particle
//=============================================================
func (p *particleEngine) newParticle() {

}

//=============================================================
// Draw the canvas
//=============================================================
func (p *particleEngine) draw(dt float64) {
	for _, p := range p.particles {
		if p.active {
			p.update(dt)
		}
	}
	p.canvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(0, 0)))
}

//=============================================================
// Build particle canvas
//=============================================================
func (c *particleEngine) build() {
	model := imdraw.New(nil)
	for _, p := range p.particles {
		if p.active {
			model.Color = pixel.RGB(
				float64(p.color>>24&0xFF)/255.0,
				float64(p.color>>16&0xFF)/255.0,
				float64(p.color>>8&0xFF)/255.0,
			).Mul(pixel.Alpha(float64(p.color&0xFF) / 255.0))

			model.Push(
				pixel.V(float64(p.x*p.size), float64(p.y*p.size)),
				pixel.V(float64(p.x*p.size+p.size), float64(p.y*p.size+p.size)),
			)
			model.Rectangle(0)
		}
	}
	p.canvas.Clear(pixel.RGBA{0, 0, 0, 0})
	model.Draw(p.canvas)
}
