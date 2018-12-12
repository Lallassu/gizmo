//=============================================================
// lights.go
//-------------------------------------------------------------
// Handle lights and pooling of lights
//=============================================================
package main

import ()

//=============================================================
// Light pool
//=============================================================
type lights struct {
	pool []light
	max  int
}

//=============================================================
// Specific light
//=============================================================
type light struct {
	pos    pixel.Vec
	color  pixel.RGBA
	spread float64
	imd    *imdraw.IMDraw
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
func (l *lights) init() {
	l.max = max
	l.pool = make([]light, wLightsMax)

	for i := 0; i < wLightsMax; i++ {
		nl := light{active: false}
		nl.imd = imdraw.New(nil)
		nl.imd.Color = pixel.Alpha(1)
		nl.imd.Push(pixel.ZV)
		nl.imd.Color = pixel.Alpha(0)
		for angle := -cl.spread / 2; angle <= cl.spread/2; angle += cl.spread / 64 {
			nl.imd.Push(pixel.V(1, 0).Rotated(angle))
		}
		nl.imd.Polygon(0)
		l.pool = append(l.pool, newl)
	}
	pe.idx = 0
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
	if newLight.size <= 0 {
		newLight.size = 1
	}
	newLight.active = true
	l.pool[l.idx : l.idx+1][0] = newLight
}

//=============================================================
// Update all active lights
//=============================================================
func (l *lights) update(dt, time float64) {
	for i, _ := range l.pool {
		if l.pool[i].active {
			l.pool[i].update(dt, time)
		}
	}
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

}
