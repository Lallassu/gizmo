//=============================================================
// lights.go
//-------------------------------------------------------------
// Handle lights and pooling of lights
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"math"
	_ "time"
)

//=============================================================
// Light structure
//=============================================================
type light struct {
	bounds      *Bounds
	color       pixel.RGBA
	angleSpread float64
	angle       float64
	radius      float32
	redrawDt    float64
	imd         *imdraw.IMDraw
	canvas      *pixelgl.Canvas
	uPosX       float32
	uPosY       float32
	life        float64
	dynamic     bool
}

//=============================================================
// Create a light
// Life == -1, infinite life
//=============================================================
func (l *light) create(x, y, angle, spread, radius float64, color pixel.RGBA, dynamic bool, life float64) {
	l.canvas = pixelgl.NewCanvas(pixel.R(0, 0, radius*2, radius*2))
	l.canvas.SetComposeMethod(pixel.ComposeOver)
	l.imd = imdraw.New(l.canvas)

	l.radius = float32(radius)
	l.angle = angle
	l.angleSpread = spread
	l.color = color
	l.life = life
	l.dynamic = dynamic

	l.bounds = &Bounds{
		X:      x,
		Y:      y,
		Width:  float64(radius * 2),
		Height: float64(radius * 2),
		entity: Entity(l),
	}

	l.uPosX = float32(l.bounds.Width / 2)
	l.uPosY = float32(l.bounds.Height / 2)
	l.canvas.SetUniform("uPosX", &l.uPosY)
	l.canvas.SetUniform("uPosY", &l.uPosY)
	l.canvas.SetUniform("uRadius", &l.radius)
	l.canvas.SetFragmentShader(fragmentShaderLight)

	global.gWorld.AddObject(l.bounds)
}

//=============================================================
// Hit
//=============================================================
func (l *light) hit(x, y, vx, vy float64, power int) {

}

//=============================================================
// Get position
//=============================================================
func (l *light) getPosition() pixel.Vec {
	return pixel.Vec{l.bounds.X, l.bounds.Y}
}

//=============================================================
// Draw
//=============================================================
func (l *light) draw(dt, elapsed float64) {
	l.redrawDt += dt

	if !l.dynamic {
		if global.gRand.randFloat() < 0.01 {
			return
		}
	}

	if l.redrawDt > 1/20 {
		l.redrawDt = 0

		if l.dynamic {
			l.life -= dt
			if l.life <= 0 {
				global.gWorld.qt.Remove(l.bounds)
				return
			}
		}
		l.canvas.Clear(pixel.RGBA{0, 0, 0, 0})
		l.shine()
	}
	//l.canvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(l.bounds.X+l.bounds.Width/2, l.bounds.Y+l.bounds.Height/2)))
	l.canvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(l.bounds.X, l.bounds.Y)))
}

//=============================================================
// Shine!
//=============================================================
func (l *light) shine() {
	addTo := float64(1 / l.radius)

	l.imd.Clear()
	l.imd.Push(pixel.Vec{l.bounds.Width / 2, l.bounds.Height / 2})
	l.imd.Color = l.color
	last := pixel.Vec{-1, -1}

	bounds := []*Bounds{}
	for _, b := range global.gWorld.qt.RetrieveIntersections(&Bounds{X: l.bounds.X - float64(l.radius), Y: l.bounds.Y - float64(l.radius), Width: float64(l.radius * 2), Height: float64(l.radius * 2)}) {
		switch b.entity.(type) {
		case *chunk:
			continue
		case *light:
			continue
		case *mob:
			continue
		}
		bounds = append(bounds, b)
	}

	// Raytrace around position (Using a bit of non-granular approach to speed up things)
	for curAngle := l.angle - (l.angleSpread / 2); curAngle < l.angle+(l.angleSpread/2); curAngle += addTo * (180 / math.Pi) * 10 {
		end := pixel.Vec{l.bounds.X, l.bounds.Y}
		rads := curAngle * (math.Pi / 180)

		// Find next foreground.
		for !global.gWorld.IsRegular(end.X, end.Y) && math.Abs((end.X-l.bounds.X)) < float64(l.radius) && math.Abs(end.Y-l.bounds.Y) < float64(l.radius) {
			// Check if object.
			next := false
			for _, b := range bounds {
				if end.X >= b.X && end.X < b.X+b.Width {
					if end.Y >= b.Y && end.Y < b.Y+b.Height {
						next = true
						break
					}
				}
			}
			if next {
				break
			}
			end.X += math.Cos(rads)
			end.Y += math.Sin(rads)
		}
		if last.X == -1 {
			last = end
		}
		l.imd.Push(pixel.Vec{end.X - l.bounds.X + l.bounds.Width/2, end.Y - l.bounds.Y + l.bounds.Height/2})
	}

	// Add the first position again so we close the polygon.
	//l.imd.Push(pixel.Vec{last.X - l.bounds.X + l.bounds.Width/2, last.Y - l.bounds.Y + l.bounds.Height/2})
	l.imd.Polygon(0)
	l.imd.Draw(l.canvas)
}
