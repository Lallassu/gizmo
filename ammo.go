package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// ammoEngine is the pool for ammunition
type ammoEngine struct {
	canvas  *pixelgl.Canvas
	bullets []ammo
	idx     int
	batch   *pixel.Batch
}

// ammo is the specific ammunition
type ammo struct {
	x      float64
	y      float64
	size   float64
	color  uint32
	life   float64
	mass   float64
	active bool
	fx     float64
	fy     float64
	power  int
	owner  entity

	vx    float64
	vy    float64
	prevX float64
	prevY float64
	mdt   float64
}

// create the ammo engine pool
func (pe *ammoEngine) create() {
	pe.bullets = make([]ammo, wAmmoMax)
	pe.canvas = pixelgl.NewCanvas(pixel.R(0, 0, 1, 1))
	pe.canvas.Clear(pixel.RGBA{R: 0, G: 0, B: 0, A: 1})
	// var fragmentShader = `
	//    #version 330 core
	//
	//    in vec2  vTexCoords;
	//    in vec2 vPosition;
	//    in vec4 vColor;
	//
	//    out vec4 fragColor;
	//
	//    uniform vec4 uTexBounds;
	//    uniform sampler2D uTexture;
	//
	//    void main() {
	//       vec4 color = vec4(0.0,0.0,0.0,1.0);
	// 	  color = vec4(1.0, 0.1,0.1,0.5);
	//    	  fragColor = color;
	//    }
	//	   `
	//	pe.canvas.SetFragmentShader(fragmentShader)
	//pe.canvas.SetSmooth(true)
	pe.batch = pixel.NewBatch(&pixel.TrianglesData{}, pe.canvas)

	for i := 0; i < wAmmoMax; i++ {
		p := ammo{active: false}
		pe.bullets = append(pe.bullets, p)
	}
	pe.idx = 0
}

// newAmmo initiates a new ammo from ammo pool
func (pe *ammoEngine) newAmmo(p ammo) {
	pe.idx++
	if pe.idx >= len(pe.bullets) {
		pe.idx = 0
	}
	newp := pe.bullets[pe.idx : pe.idx+1][0]
	newp = p
	newp.owner = p.owner
	newp.active = true
	pe.bullets[pe.idx : pe.idx+1][0] = newp
}

// update renders the ammo pool
func (pe *ammoEngine) update(dt float64) {
	pe.batch.Clear()
	for i := range pe.bullets {
		if pe.bullets[i].active {
			pe.bullets[i].update(dt)
			pe.canvas.Clear(pixel.RGBA{
				R: float64((pe.bullets[i].color >> 24 & 0xFF)) / 255.0,
				G: float64((pe.bullets[i].color >> 16 & 0xFF)) / 255.0,
				B: float64((pe.bullets[i].color >> 8 & 0xFF)) / 255.0,
				A: 0.1111})
			pe.canvas.Draw(pe.batch, pixel.IM.Scaled(pixel.ZV, 1).Moved(pixel.V(pe.bullets[i].x, pe.bullets[i].y)))
		}
	}
	pe.batch.Draw(global.gWin)
}

// explode explodes an ammo at current position
func (p *ammo) explode() {
	global.gWorld.Explode(p.x, p.y, p.power)
	p.active = false
	p.life = 0
	global.gParticleEngine.effectExplosion(p.x, p.y, p.power)
}

// update individual ammunition
func (p *ammo) update(dt float64) {
	if p.life <= 0 || !p.active {
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

	// Check if hit object.
	for _, x := range global.gWorld.qt.RetrieveIntersections(&Bounds{X: p.x, Y: p.y, Width: 1, Height: 1}) {
		skipHit := false
		switch item := x.entity.(type) {
		case *light:
			skipHit = true
		case *chunk:
			skipHit = true
		case *mob:
			switch owner := p.owner.(type) {
			case *mob:
				if item == owner {
					skipHit = true
				}
			}
		}
		if !skipHit {
			x.entity.hit(p.x, p.y, p.vx, p.vy, p.power)
			p.active = false
			break
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
