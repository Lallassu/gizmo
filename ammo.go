//=============================================================
// ammo.go
//-------------------------------------------------------------
// Ammunition
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"math"
)

type ammoEngine struct {
	canvas  *pixelgl.Canvas
	bullets []ammo
	idx     int
}

type ammo struct {
	x        float64
	y        float64
	size     float64
	color    uint32
	life     float64
	mass     float64
	ammoType ammoType
	active   bool
	fx       float64
	fy       float64
	power    int

	vx    float64
	vy    float64
	prevX float64
	prevY float64
	mdt   float64
}

//=============================================================
// Create the ammo engine pool
//=============================================================
func (pe *ammoEngine) create() {
	pe.bullets = make([]ammo, wAmmoMax)
	pe.canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(global.gWorld.height), float64(global.gWorld.width)))

	// Use a channel for ammos.
	for i := 0; i < wAmmoMax; i++ {
		p := ammo{active: false}
		pe.bullets = append(pe.bullets, p)
	}
	pe.idx = 0
}

//=============================================================
// Get new ammo
//=============================================================
func (pe *ammoEngine) newAmmo(p ammo) {
	pe.idx++
	if pe.idx >= len(pe.bullets) {
		pe.idx = 0
	}
	newp := pe.bullets[pe.idx : pe.idx+1][0]
	newp = p
	newp.active = true
	pe.bullets[pe.idx : pe.idx+1][0] = newp
}

//=============================================================
// Draw the canvas
//=============================================================
func (pe *ammoEngine) update(dt float64) {
	for i, _ := range pe.bullets {
		if pe.bullets[i].active {
			pe.bullets[i].update(dt)
		}
	}
	pe.build()
	pe.canvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(float64(global.gWorld.height/2), float64(global.gWorld.width/2))))
}

//=============================================================
// Build ammo canvas
//=============================================================
func (pe *ammoEngine) build() {
	model := imdraw.New(nil)
	for _, p := range pe.bullets {
		if p.active {
			model.Color = pixel.RGB(
				float64(p.color>>24&0xFF)/255.0,
				float64(p.color>>16&0xFF)/255.0,
				float64(p.color>>8&0xFF)/255.0,
			).Mul(pixel.Alpha(float64(p.color&0xFF) / 255.0))

			model.Push(
				pixel.V(float64(p.x), float64(p.y)),
				pixel.V(float64(p.x+p.size), float64(p.y+p.size)),
			)
			model.Rectangle(0)
		}
	}
	pe.canvas.Clear(pixel.RGBA{0, 0, 0, 0})
	model.Draw(pe.canvas)
}

//=============================================================
// Explode
//=============================================================
func (p *ammo) explode() {
	global.gWorld.Explode(p.x, p.y, p.power)
	p.active = false
	p.life = 0
	global.gParticleEngine.effectExplosion(p.x, p.y, p.power)
}

//=============================================================
// Update ammo
//=============================================================
func (p *ammo) update(dt float64) {
	if p.life <= 0 {
		p.active = false
		return
	}
	p.life -= dt
	ax := p.fx * dt * p.vx * p.mass
	ay := p.fy * dt * p.vy * p.mass

	p.prevX = p.x
	p.prevY = p.y

	// Take the largest
	// divide the other with it.
	// Lerp pixels (to get almost pixel perfect collisions)
	lx := 1.0
	ly := 1.0
	loops := 1
	if math.Abs(ax) > 1 || math.Abs(ay) > 1 {
		if math.Abs(ax) > math.Abs(ay) {
			if ax < 0 {
				lx *= -1
			}
			ly = ay / math.Abs(ax)
			loops = int(math.Abs(ax))
		} else {
			if ay < 0 {
				ly *= -1
			}
			lx = ax / math.Abs(ay)
			loops = int(math.Abs(ay))
		}
	} else {
		lx = ax
		ly = ay
	}

	for n := 0; n < loops; n++ {
		if global.gWorld.IsWall(p.x+lx, p.y+ly) { // || global.gWorld.IsObject(p.x+lx, p.y+ly) {
			p.explode()
			break
		} else {
			p.x += lx
			p.y += ly
		}
	}

	p.vy -= dt * p.fy

	if p.fx > 0 {
		p.fx -= dt * global.gWorld.gravity * p.mass
	} else {
		p.fx = 0
	}
	if p.fy > 0 && p.mass > 0 {
		p.fy -= dt * global.gWorld.gravity * p.mass
	}

	if p.prevX-p.x == 0 && p.prevY-p.y == 0 {
		p.mdt += dt
	} else {
		p.mdt = 0
	}

	if p.life <= 0 || p.mdt > 1 {
		p.explode()
	}
}
