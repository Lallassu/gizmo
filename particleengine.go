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
// Explosion effect
//=============================================================
func (pe *particleEngine) effectExplosion(x, y float64, size int) {
	// Create fire part
	for i := 0; i < size; i++ {
		r := 0xF5
		g := 50 + rand.Intn(180)
		b := 16
		a := 20 + rand.Intn(220)
		pe.newParticle(particle{
			color:       uint32((r & 0xFF << 24) | (g & 0xFF << 16) | (b & 0xFF << 8) | (a & 0xFF)),
			size:        rand.Float64() * 2,
			x:           x + float64(size/2-rand.Intn(size)),
			y:           y + float64(size/2-rand.Intn(size)),
			vx:          0,
			vy:          rand.Float64() * 10,
			fx:          10,
			fy:          10,
			life:        rand.Float64() * 5,
			restitution: 0,
		})
	}
	// Create smoke
	for i := 0; i < size*2; i++ {
		c := 50 + rand.Intn(205)
		a := 20 + rand.Intn(220)
		pe.newParticle(particle{
			color:       uint32((c & 0xFF << 24) | (c & 0xFF << 16) | (c & 0xFF << 8) | (a & 0xFF)),
			size:        rand.Float64() * 2,
			x:           x + float64(size/2-rand.Intn(size)),
			y:           y + float64(size/2-rand.Intn(size)),
			vx:          0,
			vy:          rand.Float64() * 10,
			fy:          -rand.Float64() * 10,
			fx:          0,
			life:        rand.Float64() * 1.5,
			mass:        -0.1,
			restitution: 0,
		})
	}
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
	newp.bounces = 0
	//newp.life = wParticleDefaultLife
	//newp.restitution = -0.3
	//newp.fx = 10
	//newp.fy = 10
	//newp.vx = float64(5 - rand.Intn(10))
	//newp.vy = float64(5 - rand.Intn(10))
	//newp.x = p.x
	//newp.y = p.y
	//newp.size = 1
	//newp.mass = 2 * rand.Float64()
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
