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
)

//=============================================================
// Light pool
//=============================================================
type lights struct {
	pool   []light
	idx    int
	canvas *pixelgl.Canvas
}

//=============================================================
// Specific light
//=============================================================
type light struct {
	pos    pixel.Vec
	color  pixel.RGBA
	spread float64
	imd    *imdraw.IMDraw
	canvas *pixelgl.Canvas
	angle  float64
	radius float64
	active bool
	life   float64
}

//=============================================================
//
// Light pool
//
//=============================================================

//=============================================================
// Init light pool
//=============================================================
func (l *lights) create() {
	l.pool = make([]light, wLightsMax)
	l.canvas = pixelgl.NewCanvas(global.gWin.Bounds())

	for i := 0; i < wLightsMax; i++ {
		nl := light{active: true}
		nl.canvas = pixelgl.NewCanvas(global.gWin.Bounds())
		nl.pos = pixel.Vec{100.0, 100.0 + float64(i)*10.0}
		nl.imd = imdraw.New(nil)
		nl.imd.Color = pixel.RGBA{1, 0, 0, 1}
		nl.imd.Push(pixel.ZV)
		nl.imd.Color = pixel.RGBA{0, 0, 0, 0}
		nl.spread = 10
		nl.radius = math.Pi
		for angle := -nl.spread / 2; angle <= nl.spread/2; angle += nl.spread / 64 {
			nl.imd.Push(pixel.V(1, 0).Rotated(angle))
		}
		nl.imd.Polygon(0)
		l.pool = append(l.pool, nl)
	}
	l.idx = 0
}

//=============================================================
// Create a new light from pool
//=============================================================
func (l *lights) newLight(newl light) {
	l.idx++
	if l.idx >= len(l.pool) {
		l.idx = 0
	}
	newLight := l.pool[l.idx : l.idx+1][0]

	newLight = newl
	newLight.active = true
	l.pool[l.idx : l.idx+1][0] = newLight
}

//=============================================================
// Update all active lights
//=============================================================
func (l *lights) update(dt, time float64) {
	l.canvas.Clear(pixel.Alpha(0))
	l.canvas.SetComposeMethod(pixel.ComposePlus)

	for i, _ := range l.pool {
		if l.pool[i].active {
			l.pool[i].update(dt, time)
			l.pool[i].canvas.Draw(l.canvas, pixel.IM.Moved(l.canvas.Bounds().Center()))
		}
	}
	global.gWin.SetColorMask(pixel.Alpha(1))
	l.canvas.Draw(global.gWin, pixel.IM.Moved(global.gWin.Bounds().Center()))
}

//=============================================================
//
// Individual Lights
//
//=============================================================

//=============================================================
// Update light
//=============================================================
func (l *light) update(dt, time float64) {
	l.canvas.Clear(pixel.Alpha(0))
	l.pos.X = l.pos.X + math.Sin(time/2)
	l.canvas.SetMatrix(pixel.IM.Scaled(pixel.ZV, l.radius).Rotated(pixel.ZV, l.angle).Moved(l.pos))
	l.canvas.SetColorMask(pixel.Alpha(1))
	l.canvas.SetComposeMethod(pixel.ComposePlus)
	l.imd.Draw(l.canvas)
}
