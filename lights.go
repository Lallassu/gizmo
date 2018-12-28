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

var lights []light

var limd *imdraw.IMDraw
var lcanvas *pixelgl.Canvas
var li *light

var fragmentShaderLight = `
             #version 330 core
             
             in vec2  vTexCoords;
             in vec4  vColor;
             
             out vec4 fragColor;
             
			 uniform float uPosX;
			 uniform float uPosY;
             uniform vec4 uTexBounds;
             uniform sampler2D uTexture;
             
             void main() {
				vec4 c = vColor;
			    c *= 2;
			    vec3 fc = vec3(1.0, 0.3, 0.1);
	            vec2 borderSize = vec2(0.1); 

	            vec2 rectangleSize = vec2(1.0) - borderSize; 

				float dist = distance(vec2(uPosX, uPosY), vTexCoords);

	            //float distanceField = length(max(abs(c.x)-rectangleSize,0.0) / borderSize);

	            //float alpha = 1.0 - distanceField;
			    fc *= abs(0.3 / (sin( c.x + sin(c.y)+ 1.3 ) * 1.0) );
                fragColor = vec4(fc, 0.1);
             }
             
			 `
var uPosX float32
var uPosY float32

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

func (l *light) shineLight() {
	//dynLen := l.radius
	//start := time.Now()
	addTo := 1 / l.radius

	limd.Clear()
	for curAngle := l.angle - (l.angleSpread / 2); curAngle < l.angle+(l.angleSpread/2); curAngle += addTo * (180 / math.Pi) * 2 {
		// Find next foreground.
		end := l.position
		rads := curAngle * (math.Pi / 180)
		dist := 0.0
		for !global.gWorld.IsRegular(end.X, end.Y) && dist < l.radius {
			dist += 1
			end.X += math.Cos(rads)
			end.Y += math.Sin(rads)
		}
		limd.Color = pixel.RGBA{0.1, 0.1, 0.0, 0}
		limd.Push(pixel.Vec{l.position.X, l.position.Y})
		limd.Push(pixel.Vec{end.X, end.Y})
		limd.Line(1)

		//Debug("FROM:", l.position, "TO:", end)
	}
	limd.Draw(lcanvas)
	//elapsed := time.Since(start)
	//Debug("Build took %s", elapsed)
}

var angle float64
var totaldt float64

func drawLights(dt float64) {
	//li.angle += 0.6

	// TBD: Don't update if position hasn't changed

	lcanvas.Clear(pixel.RGBA{0, 0, 0, 0})
	li.position = global.gPlayer.getPosition()
	li.position.Y += 10
	uPosX = float32(li.position.X)
	uPosY = float32(li.position.Y)
	li.shineLight()
	lcanvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(1000, 1000)))

}

func createLights() {
	lcanvas = pixelgl.NewCanvas(pixel.R(0, 0, 2000, 2000))
	lcanvas.SetUniform("uPosX", &uPosX)
	lcanvas.SetUniform("uPosY", &uPosY)
	lcanvas.SetFragmentShader(fragmentShaderLight)
	lcanvas.Clear(pixel.RGBA{0.0, 0.01, 0, 0.1})
	limd = imdraw.New(lcanvas)

	li = &light{position: global.gPlayer.getPosition(), angleSpread: 360, angle: 300, color: pixel.RGBA{0, 0, 0.4, 0.1}}
	li.radius = 100

}
