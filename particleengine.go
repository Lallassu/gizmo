//=============================================================
// particleengine.go
//-------------------------------------------------------------
// Particles of different kinds.
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	_ "sync"
)

type particleEngine struct {
	canvas    *pixelgl.Canvas
	particles []particle
	idx       int
	batch     *pixel.Batch
	colors    []uint8
	colormap  map[uint32]int
	//	cMutex    *sync.Mutex
}

//=============================================================
// Blood effect
//=============================================================
func (pe *particleEngine) effectBlood(x, y, vx, vy float64, size int) {
	//	for i := 0; i < 1; i++ {
	r := 175 + global.gRand.rand()*5
	g := 10 + global.gRand.rand()*2
	b := 10 + global.gRand.rand()*2
	a := 255 //global.gRand.rand() * 255

	pe.newParticle(particle{
		x:           float64(x),
		y:           float64(y),
		size:        global.gRand.randFloat() * 3,
		restitution: -0.1 - global.gRand.randFloat()/4,
		life:        wParticleDefaultLife,
		fx:          5 + global.gRand.randFloat()*5,
		fy:          5 + global.gRand.randFloat()*5,
		vx:          vx / 2,
		vy:          float64(5 - global.gRand.rand()),
		mass:        2,
		pType:       particleBlood,
		color:       uint32(r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | a&0xFF),
		static:      true,
	})
	//}

}
func (pe *particleEngine) ammoShell(x, y, dir, size float64) {
	r := 0xFF
	g := 0xD7
	b := 0
	a := 255

	pe.newParticle(particle{
		color:       uint32(r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | a&0xFF),
		size:        size,
		x:           x,
		y:           y,
		vx:          -global.gRand.randFloat() * 3 * dir,
		vy:          2,
		fx:          5,
		fy:          5,
		life:        5,
		mass:        1,
		pType:       particleRegular,
		restitution: -0.3,
	})
}

func (pe *particleEngine) effectSmoke(x, y float64, size int) {
	for i := 0; i < size*2; i++ {
		c := 50 + global.gRand.rand()*20
		a := 20 + global.gRand.rand()*200
		pe.newParticle(particle{
			color:       uint32(c&0xFF<<24 | c&0xFF<<16 | c&0xFF<<8 | a&0xFF),
			size:        global.gRand.randFloat() * 2.5,
			x:           x + float64(size/2) - global.gRand.randFloat()*float64(size) + global.gRand.randFloat()*2,
			y:           y + float64(size/2) - global.gRand.randFloat()*float64(size) + global.gRand.randFloat()*2,
			vx:          0,
			vy:          global.gRand.randFloat() * 10,
			fy:          -global.gRand.randFloat() * 10,
			fx:          0,
			life:        global.gRand.randFloat() * 3.5,
			mass:        -0.1,
			pType:       particleSmoke,
			restitution: 0,
		})
	}
}

func (pe *particleEngine) effectExplosion(x, y float64, size int) {
	// Create fire part
	for i := 0; i < size*5; i++ {
		r := 0xF9
		g := 50 + global.gRand.rand()*14
		b := 16
		a := 20 + global.gRand.rand()*22

		pe.newParticle(particle{
			color:       uint32(r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | a&0xFF),
			size:        global.gRand.randFloat() * 2,
			x:           x,
			y:           y,
			vx:          float64(5 - global.gRand.rand()),
			vy:          float64(5 - global.gRand.rand()),
			fx:          10,
			fy:          10,
			life:        global.gRand.randFloat() * 2,
			mass:        1,
			pType:       particleFire,
			restitution: 0,
		})

	}
	// Create smoke
	for i := 0; i < size*2; i++ {
		c := 50 + global.gRand.rand()*20
		a := 20 + global.gRand.rand()*200
		pe.newParticle(particle{
			color:       uint32(c&0xFF<<24 | c&0xFF<<16 | c&0xFF<<8 | a&0xFF),
			size:        global.gRand.randFloat() * 2.5,
			x:           x + float64(size/2) - global.gRand.randFloat()*float64(size) + global.gRand.randFloat()*2,
			y:           y + float64(size/2) - global.gRand.randFloat()*float64(size) + global.gRand.randFloat()*2,
			vx:          0,
			vy:          global.gRand.randFloat() * 10,
			fy:          -global.gRand.randFloat() * 10,
			fx:          0,
			life:        global.gRand.randFloat() * 3.5,
			mass:        -0.1,
			pType:       particleSmoke,
			restitution: 0,
		})
	}
}

//=============================================================
// Add or verify that the color exists in batch canvas.
//=============================================================
func (pe *particleEngine) addColorToBatch(color uint32) {
	pos := color | 0xFF
	if _, ok := pe.colormap[pos]; !ok {
		r := color >> 24 & 0xFF
		g := color >> 16 & 0xFF
		b := color >> 8 & 0xFF
		pe.colors = append(pe.colors, uint8(r), uint8(g), uint8(b), 255.0)
		pe.colormap[pos] = (len(pe.colors) / 4) - 1

		//	pe.cMutex.Lock()
		pe.canvas.SetBounds(pixel.R(0, 0, float64(len(pe.colors)/4), 1))
		pe.canvas.SetPixels(pe.colors)
		//	pe.cMutex.Unlock()
	}
}

//=============================================================
// Create the particle engine pool
//=============================================================
func (pe *particleEngine) create() {
	pe.particles = make([]particle, wParticlesMax)
	pe.colormap = make(map[uint32]int)
	//	pe.cMutex = &sync.Mutex{}

	pe.canvas = pixelgl.NewCanvas(pixel.R(0, 0, 1, 1)) // Max seems to be 2^14 per row
	pe.batch = pixel.NewBatch(&pixel.TrianglesData{}, pe.canvas)

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

	// Check if color are defined if not,create and add to batch
	pe.addColorToBatch(p.color)

	// Make a shallow copy, no pointers in particle so we're fine.
	if p.size <= 0 {
		p.size = 1
	}
	newp = p
	newp.active = true

	pe.particles[pe.idx : pe.idx+1][0] = newp
}

//=============================================================
// Draw the canvas
//=============================================================
func (pe *particleEngine) update(dt float64) {
	pe.batch.Clear()

	//pe.cMutex.Lock()
	sprite := pixel.NewSprite(pe.canvas, pixel.R(0, 0, 1, 1))
	for i, _ := range pe.particles {
		if pe.particles[i].active {
			pe.particles[i].update(dt)
			color := pe.particles[i].color
			r := color >> 24 & 0xFF
			g := color >> 16 & 0xFF
			b := color >> 8 & 0xFF
			a := color & 0xFF

			pos := r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | 0xFF
			sprite.Set(pe.canvas, pixel.R(float64(pe.colormap[pos]), 0, float64(pe.colormap[pos]+1), 1))
			if pe.particles[i].pType == particleRegular {
				sprite.Draw(pe.batch, pixel.IM.Scaled(pixel.ZV, pe.particles[i].size).Moved(pixel.V(pe.particles[i].x, pe.particles[i].y)))
			} else {
				sprite.DrawColorMask(pe.batch, pixel.IM.Scaled(pixel.ZV, pe.particles[i].size).Moved(pixel.V(pe.particles[i].x, pe.particles[i].y)), pixel.RGBA{float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0, float64(a) / 255.0})
				if pe.particles[i].pType == particleFire {
					sprite.DrawColorMask(pe.batch, pixel.IM.Scaled(pixel.ZV, pe.particles[i].size*pe.particles[i].life*2).Moved(pixel.V(pe.particles[i].x, pe.particles[i].y)), pixel.RGBA{float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0, 0.1111})
				}
				// if pe.particles[i].pType == particleFire {
				// 	pe.imd.Color = pixel.RGBA{1, 0, 0, 1.0}
				// 	pe.imd.Push(pixel.Vec{pe.particles[i].x, pe.particles[i].y})
				// 	//	pe.imd.Push(pixel.Vec{pe.particles[i].x + 4, pe.particles[i].y + 4})
				// 	pe.imd.Circle(1, 2)
				// 	pe.trails = append(pe.trails, &trail{pos: pixel.Vec{pe.particles[i].x, pe.particles[i].y}, ts: pe.utime})
				// }
			}
		}
	}
	//	pe.cMutex.Unlock()

	pe.batch.Draw(global.gWin)
}
