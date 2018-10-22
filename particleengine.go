//=============================================================
// particleengine.go
//-------------------------------------------------------------
// Particles of different kinds.
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"math/rand"
)

type particleEngine struct {
	canvas    *pixelgl.Canvas
	particles []particle
	idx       int
}

//=============================================================
// Draw the canvas
//=============================================================
func (pe *particleEngine) create() {
	pe.particles = make([]particle, wParticlesMax)
	pe.canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(global.gWindowHeight), float64(global.gWindowWidth)))

	// Use a channel for particles.
	for i := 0; i < wParticlesMax; i++ {
		p := particle{active: false}
		pe.particles = append(pe.particles, p)
	}
	pe.idx = 0
}

//=============================================================
// Get new particle
//=============================================================
func (pe *particleEngine) newParticle(p particle) {
	pe.idx++
	if pe.idx >= len(pe.particles) {
		pe.idx = 0
	}
	newp := pe.particles[pe.idx : pe.idx+1][0]
	// Make a shallow copy, no pointers in particle so we're fine.
	newp = p
	newp.active = true
	newp.life = wParticleDefaultLife
	newp.restitution = -0.3
	newp.fx = 10
	newp.fy = 10
	newp.vx = float64(5 - rand.Intn(10))
	newp.vy = float64(5 - rand.Intn(10))
	newp.bounces = 0
	newp.x = p.x
	newp.y = p.y
	newp.size = 1
	newp.mass = 2 * rand.Float64()
	pe.particles[pe.idx : pe.idx+1][0] = newp
}

//=============================================================
// Draw the canvas
//=============================================================
func (pe *particleEngine) update(dt float64) {
	for i, _ := range pe.particles {
		if pe.particles[i].active {
			pe.particles[i].update(dt)
		}
	}
	pe.build()
	pe.canvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(float64(global.gWindowHeight/2), float64(global.gWindowWidth/2))))
}

//=============================================================
// Build particle canvas
//=============================================================
func (pe *particleEngine) build() {
	model := imdraw.New(nil)
	for _, p := range pe.particles {
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
	pe.canvas.Clear(pixel.RGBA{0, 0, 0, 0})
	model.Draw(pe.canvas)
}
