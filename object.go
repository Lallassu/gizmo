//=============================================================
// object.go
//-------------------------------------------------------------
// Different objects that are subject to physics and flood fill
// for destuction into pieces. Much like mob, but special adds
// like FF and different physics and actions.
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"image"
)

type object struct {
	textureFile string
	img         image.Image
	model       *imdraw.IMDraw
	canvas      *pixelgl.Canvas
	bounds      *Bounds
	objectType  entityType
	mass        float64
	restitution float64
	height      int
	width       int
	size        int
	force       pixel.Vec
	pixels      []uint32
	prevPos     []pixel.Vec
}

//=============================================================
//
//=============================================================
func (o *object) create(x, y float64) {
	o.prevPos = make([]pixel.Vec, 100)

	o.mass = 50

	o.img, o.width, o.height, o.size = loadTexture(o.textureFile)

	// Initiate bounds for qt
	o.bounds = &Bounds{
		X:      x,
		Y:      y,
		Width:  float64(o.width),
		Height: float64(o.height),
		entity: Entity(o),
	}

	o.pixels = make([]uint32, o.size*o.size)

	for x := 0; x <= o.width; x++ {
		for y := 0; y <= o.height; y++ {
			r, g, b, a := o.img.At(x, o.height-y).RGBA()
			o.pixels[x*o.size+y] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | a&0xFF
		}
	}

	// Generate some CD pixel for faster CD check.
	//rand.Seed(time.Now().UTC().UnixNano())
	//for x := 0; x < 20; x++ {
	//	o.cdPixels = append(o.cdPixels, [2]uint32{uint32(rand.Intn(o.width)), uint32(rand.Intn(o.height))})
	//}

	o.canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(o.width), float64(o.height)))

	// build initial
	o.build()

	// Add object to QT
	global.gWorld.AddObject(o.bounds)
}

//=============================================================
// Build
//=============================================================
func (o *object) build() {
	o.model = imdraw.New(nil)
	for x := 0; x < o.width; x++ {
		for y := 0; y < o.height; y++ {
			p := o.pixels[x*o.size+y]
			if p == 0 {
				continue
			}

			o.model.Color = pixel.RGB(
				float64(p>>24&0xFF)/255.0,
				float64(p>>16&0xFF)/255.0,
				float64(p>>8&0xFF)/255.0,
			).Mul(pixel.Alpha(float64(p&0xFF) / 255.0))
			o.model.Push(
				pixel.V(float64(x*wPixelSize), float64(y*wPixelSize)),
				pixel.V(float64(x*wPixelSize+wPixelSize), float64(y*wPixelSize+wPixelSize)),
			)
			o.model.Rectangle(0)
		}
	}

	o.canvas.Clear(pixel.RGBA{0, 0, 0, 0})
	o.model.Draw(o.canvas)
}

//=============================================================
//
//  Function to implement Entity interface
//
//=============================================================
//=============================================================
//
//=============================================================
func (o *object) hit(x_, y_ float64) bool {
	//x := int(math.Abs(float64(o.bounds.X - x_)))
	//y := int(math.Abs(float64(o.bounds.Y - y_)))

	o.build()
	return true
}

//=============================================================
//
//=============================================================
func (o *object) explode() {
}

//=============================================================
//
//=============================================================
func (o *object) move(dx, dy float64) {
	// Add the force, movenment is handled in the physics function
	// o.force.X += dx * o.speed
	// o.force.Y += dy * o.speed
}

//=============================================================
//
//=============================================================
func (o *object) getPosition() pixel.Vec {
	return pixel.Vec{o.bounds.X, o.bounds.Y}
}

//=============================================================
//
//=============================================================
func (o *object) getMass() float64 {
	return o.mass
}

//=============================================================
//
//=============================================================
func (o *object) getType() entityType {
	return entityObject
}

//=============================================================
//
//=============================================================
func (o *object) setPosition(x, y float64) {
	o.bounds.X = x
	o.bounds.Y = y
}

//=============================================================
//
//=============================================================
func (o *object) saveMove() {
	o.prevPos = append(o.prevPos, pixel.Vec{o.bounds.X, o.bounds.Y})
	// TBD: Only remove every second or something
	if len(o.prevPos) > 100 {
		o.prevPos = o.prevPos[:100]
	}
}

//=============================================================
// Physics
//=============================================================
func (o *object) physics(dt float64) {
	o.saveMove()
}

//=============================================================
//
//=============================================================
func (o *object) draw(dt float64) {
	// Update physics
	o.physics(dt)

}
