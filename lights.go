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
	"math/rand"
)

var blocks []block
var lights []light

var limd *imdraw.IMDraw
var lcanvas *pixelgl.Canvas

//=============================================================
// Specific light
//=============================================================
type light struct {
	position    pixel.Vec
	color       pixel.RGBA
	angleSpread float64
	angle       float64
	radius      float64
}

type block struct {
	position pixel.Vec
	width    float64
	height   float64
	visible  bool
}

func (l *light) findDistance(b_ block, angle, rLen_ float64, start_ bool, shortest_ float64, closestBlock_ block) (start bool, rLen, shortest float64, closestBlock block) {
	b := b_
	rLen = rLen_
	start = start_
	shortest = shortest_
	closestBlock = closestBlock_

	y := (b.position.Y + b.height/2) - l.position.Y
	x := (b.position.X + b.width/2) - l.position.X
	dist := math.Sqrt((y * y) + (x * x))

	if l.radius >= dist {
		rads := angle * (math.Pi / 180)
		pointPos := pixel.V(l.position.X, l.position.Y)

		pointPos.X += math.Cos(rads) * dist
		pointPos.Y += math.Sin(rads) * dist

		if pointPos.X > b.position.X && pointPos.X < b.position.X+b.width && pointPos.Y > b.position.Y && pointPos.Y < b.position.Y+b.height {
			if start || dist < shortest {
				start = false
				shortest = dist
				rLen = dist
				closestBlock = b
			}
			return
		}
	}
	return
}

func (l *light) shineLight() {
	curAngle := l.angle - (l.angleSpread / 2)
	dynLen := l.radius
	addTo := 1 / l.radius

	for ; curAngle < l.angle+(l.angleSpread/2); curAngle += (addTo * (180 / math.Pi)) * 2 {
		dynLen = l.radius

		start := true
		shortest := 0.0
		rLen := dynLen
		b := block{}

		for i := 0; i < len(blocks); i++ {
			start, rLen, shortest, b = l.findDistance(blocks[i], curAngle, rLen, start, shortest, b)
		}

		rads := curAngle * (math.Pi / 180)
		end := pixel.Vec{l.position.X, l.position.Y}

		b.visible = true
		end.X += math.Cos(rads) * rLen
		end.Y += math.Sin(rads) * rLen

		// DRAW IMD
		limd.Color = pixel.RGBA{0.2, 0, 0.2, 0.1}
		limd.Push(pixel.Vec{l.position.X, l.position.Y})
		limd.Push(pixel.Vec{end.X, end.Y})
		limd.Line(1)
		// ctx.beginPath()
		// ctx.moveTo(l.position.x, l.position.y)
		// ctx.lineTo(end.x, end.y)
		// ctx.closePath()
		// // ctx.clip();
		// ctx.stroke()
	}
}

var angle float64
var totaldt float64

func drawLights(dt float64) {
	angle += 0.6

	lcanvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(600, 600)))

}

func createLights() {
	lcanvas = pixelgl.NewCanvas(pixel.R(0, 0, 1000, 1000))
	lcanvas.Clear(pixel.RGBA{0.0, 0, 0, 0.0})
	limd = imdraw.New(lcanvas)

	for i := 0; i < 50; i++ {
		size := float64(rand.Intn(20) + 10)
		blocks = append(blocks, block{position: pixel.Vec{float64(rand.Intn(512)), float64(rand.Intn(512))}, width: size, height: size})
	}
	for _, b := range blocks {
		if b.visible {
			limd.Color = pixel.RGBA{1.0, 0, 1.0, 1}
			limd.Push(pixel.Vec{b.position.X, b.position.Y})
			limd.Push(pixel.Vec{b.position.X + b.width, b.position.Y + b.height})
			limd.Rectangle(0)
			b.visible = false
		} else {
			limd.Color = pixel.RGBA{0, 1.00, 0, 1}
			limd.Push(pixel.Vec{b.position.X, b.position.Y})
			limd.Push(pixel.Vec{b.position.X + b.width, b.position.Y + b.height})
			limd.Rectangle(0)
		}
	}
	lights = append(lights, light{position: pixel.Vec{300, 300}, angleSpread: 300, angle: 300, color: pixel.RGBA{0, 0, 0.4, 0.1}})
	for i, _ := range lights {
		lights[i].radius = 1000
		lights[i].shineLight()
	}
	limd.Draw(lcanvas)

}
