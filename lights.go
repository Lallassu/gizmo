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
}

//=============================================================
// Create a light
//=============================================================
func (l *light) create(x, y, angle, spread, radius float64, color pixel.RGBA) {
	l.canvas = pixelgl.NewCanvas(pixel.R(0, 0, radius+radius/2, radius+radius/2))
	l.canvas.SetComposeMethod(pixel.ComposeOver)
	l.imd = imdraw.New(l.canvas)

	l.radius = float32(radius)
	l.angle = angle
	l.angleSpread = spread
	l.color = color

	l.canvas.SetUniform("uPosX", &l.uPosX)
	l.canvas.SetUniform("uPosY", &l.uPosY)
	l.canvas.SetUniform("uRadius", &l.radius)
	l.canvas.SetFragmentShader(fragmentShaderLight)

	l.bounds = &Bounds{
		X:      x,
		Y:      y,
		Width:  float64(radius + radius/2),
		Height: float64(radius + radius/2),
		entity: Entity(l),
	}
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

	if l.redrawDt > 1/30 {
		l.redrawDt = 0
		l.uPosX = float32(l.bounds.X)
		l.uPosY = float32(l.bounds.Y)

		l.canvas.Clear(pixel.RGBA{0.0, 0.0, 0.0, 0.0})
		l.shine()
	}
	l.canvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(l.bounds.X+l.bounds.Width/2, l.bounds.Y+l.bounds.Height/2)))
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

	for curAngle := l.angle - (l.angleSpread / 2); curAngle < l.angle+(l.angleSpread/2); curAngle += addTo * (180 / math.Pi) * 10 {
		end := pixel.Vec{l.bounds.Width / 2, l.bounds.Height / 2}
		rads := curAngle * (math.Pi / 180)

		dist := float32(0.0)

		// Find next foreground.
		// Incr radius to make it fade away
		for !global.gWorld.IsRegular(end.X, end.Y) && dist < l.radius+(l.radius/2) {
			dist += 1
			end.X += math.Cos(rads)
			end.Y += math.Sin(rads)
		}
		if last.X == -1 {
			last = end
		}
		l.imd.Push(pixel.Vec{end.X, end.Y})
	}

	// Add the first position again so we close the polygon.
	l.imd.Push(pixel.Vec{last.X, last.Y})
	l.imd.Polygon(0)
	l.imd.Draw(l.canvas)
}

//func drawLights(dt float64) {
//	//li.angle += 0.6
//	li.redrawDt += dt
//
//	if li.redrawDt > 1/30 {
//		li.redrawDt = 0
//
//		li.position = global.gPlayer.getPosition()
//		li.position.Y += 10
//		uPosX = float32(li.position.X)
//		uPosY = float32(li.position.Y)
//
//		lcanvas.Clear(pixel.RGBA{0.0, 0.0, 0.0, 0.0})
//		li.shine()
//	}
//	lcanvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(512, 512)))
//
//}
//
//func createLights() {
//	lcanvas = pixelgl.NewCanvas(pixel.R(0, 0, 1024, 1024))
//	lcanvas.SetComposeMethod(pixel.ComposeOver)
//	limd = imdraw.New(lcanvas)
//
//	li = &light{position: global.gPlayer.getPosition(), angleSpread: 360, angle: 90, color: pixel.RGBA{0, 0, 0.4, 0.1}}
//	li.radius = 100
//
//	lcanvas.SetUniform("uPosX", &li.uPosX)
//	lcanvas.SetUniform("uPosY", &li.uPosY)
//	lcanvas.SetUniform("uRadius", &li.radius)
//	lcanvas.SetFragmentShader(fragmentShaderLight)
//
//}
