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
	_ "github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"image"
	"math"
)

type object struct {
	name        string
	textureFile string
	mass        float64
	restitution float64
	height      float64
	width       float64
	size        float64
	force       pixel.Vec
	pixels      []uint32
	prevPos     []pixel.Vec
	bounces     int
	vx          float64
	vy          float64
	fx          float64
	fy          float64
	scale       float64
	owner       Entity
	rotation    float64
	active      bool
	reloadTime  float64
	falling     bool
	animateIdle bool
	sprite      *sprite
	//  img         image.Image
	//  model       *imdraw.IMDraw
	canvas *pixelgl.Canvas
	bounds *Bounds

	batch     *pixel.Batch
	triangles *pixel.TrianglesData
	dirty     bool
}

//=============================================================
//
//=============================================================
func (o *object) create(x_, y_ float64) {
	o.prevPos = make([]pixel.Vec, 100)
	o.rotation = 0.1
	o.mass = 5
	o.restitution = -0.3
	o.fx = 1
	o.fy = 1
	o.vx = 1
	o.vy = 1
	o.active = true
	o.animateIdle = false

	o.dirty = true
	o.triangles = pixel.MakeTrianglesData(100)
	o.batch = pixel.NewBatch(o.triangles, nil)

	var img image.Image
	img, o.width, o.height, o.size = loadTexture(o.textureFile)

	// Initiate bounds for qt
	o.bounds = &Bounds{
		X:      x_,
		Y:      y_,
		Width:  o.width * o.scale,
		Height: o.height * o.scale,
		entity: Entity(o),
	}

	o.pixels = make([]uint32, int(o.size*o.size))

	for x := 0.0; x < o.width; x++ {
		for y := 0.0; y < o.height; y++ {
			r, g, b, a := img.At(int(x), int(o.height-y)).RGBA()
			if r == 0 && g == 0 && b == 0 && a == 0 {
				continue
			}
			o.pixels[int(x*o.size+y)] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | a&0xFF
		}
	}

	o.canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(o.width), float64(o.height)))
	//var fragmentShader = `
	// #version 330 core

	// in vec2  vTexCoords;
	// in vec2 vPosition;
	// in vec4 vColor;

	// out vec4 fragColor;

	// uniform float uCenterX;
	// uniform float uCenterY;

	// uniform vec4 uTexBounds;
	// uniform sampler2D uTexture;

	// void main() {
	//    vec4 color = vec4(0.0,0.0,0.0,0.0);
	//    float d = sqrt(pow(vPosition.x-uCenterX, 2) + pow(vPosition.y-uCenterY, 2));
	//	if (d <= 0) {
	//		d = 1;
	//	}
	//    color = vec4(vColor.g*(uCenterX/100), vColor.g, vColor.b, vColor.a/d);
	//    fragColor = color;
	// }
	//    `
	//o.canvas.SetFragmentShader(fragmentShader)

	// build initial
	o.build()

	// Add object to QT
	global.gWorld.AddObject(o.bounds)
}

//=============================================================
// Build
//=============================================================
func (o *object) build() {

	// if !o.static {
	// 	o.model = imdraw.New(nil)
	// 	for x := 0; x < o.width; x++ {
	// 		for y := 0; y < o.height; y++ {
	// 			p := o.pixels[x*o.size+y]
	// 			if p == 0 {
	// 				continue
	// 			}

	// 			o.model.Color = pixel.RGB(
	// 				float64(p>>24&0xFF)/255.0,
	// 				float64(p>>16&0xFF)/255.0,
	// 				float64(p>>8&0xFF)/255.0,
	// 			).Mul(pixel.Alpha(float64(p&0xFF) / 255.0))
	// 			o.model.Push(
	// 				pixel.V(float64(x*wPixelSize), float64(y*wPixelSize)),
	// 				pixel.V(float64(x*wPixelSize+wPixelSize), float64(y*wPixelSize+wPixelSize)),
	// 			)
	// 			o.model.Rectangle(0)
	// 		}
	// 	}

	// 	o.canvas.Clear(pixel.RGBA{0, 0, 0, 0})
	// 	o.model.Draw(o.canvas)
	// }
	i := 0
	rc := uint32(0)
	gc := uint32(0)
	bc := uint32(0)
	p2 := uint32(0)
	r1 := uint32(0)
	g1 := uint32(0)
	b1 := uint32(0)
	draw := 0
	same_x := 1.0
	same_y := 1.0
	pos := 0

	for x := 0.0; x < o.width; x++ {
		for y := 0.0; y < o.height; y++ {
			p := o.pixels[int(x*o.size+y)]
			// Skip visisted or empty
			if p == 0 || p&0xFF>>7 == 0 {
				continue
			}
			rc = p >> 24 & 0xFF
			gc = p >> 16 & 0xFF
			bc = p >> 8 & 0xFF
			same_x = 1.0
			same_y = 1.0

			for l := x + 1; l < o.width; l++ {
				// Check color
				pos = int(l*o.size + y)
				p2 = o.pixels[pos]
				if p2 == 0 {
					break
				}
				r1 = p2 >> 24 & 0xFF
				g1 = p2 >> 16 & 0xFF
				b1 = p2 >> 8 & 0xFF

				if r1 == rc && g1 == gc && b1 == bc && ((p2&0xFF)>>7) == 1 {
					// Same color and not yet visited!
					o.pixels[pos] &= 0xFFFFFF7F
					same_x++
					new_y := 1.0
					for k := y; k < o.height; k++ {
						pos = int(l*o.size + k)
						p2 = o.pixels[pos]
						if p2 == 0 {
							break
						}
						r1 = p2 >> 24 & 0xFF
						g1 = p2 >> 16 & 0xFF
						b1 = p2 >> 8 & 0xFF

						if r1 == rc && g1 == gc && b1 == bc && ((p2&0xFF)>>7) == 1 {
							o.pixels[pos] &= 0xFFFFFF7F
							new_y++
						} else {
							break
						}
					}
					if new_y < same_y {
						break
					} else {
						same_y = new_y
					}
				} else {
					break
				}
			}

			draw++

			// Convert to decimal
			r := float64(p>>24&0xFF) / 255.0
			g := float64(p>>16&0xFF) / 255.0
			b := float64(p>>8&0xFF) / 255.0
			a := float64(p&0xFF) / 255.0

			// Increase length of triangles if we need to draw more than we had before.
			// Add a buffer so we can skip a few increments.
			if draw*6 >= len(*o.triangles) {
				o.triangles.SetLen(draw*6 + 10)
			}

			// Size of triangle is given by how large the greedy algorithm found out.
			(*o.triangles)[i].Position = pixel.Vec{x, y}
			(*o.triangles)[i+1].Position = pixel.Vec{x + same_x, y}
			(*o.triangles)[i+2].Position = pixel.Vec{x + same_x, y + same_y}
			(*o.triangles)[i+3].Position = pixel.Vec{x, y}
			(*o.triangles)[i+4].Position = pixel.Vec{x, y + same_y}
			(*o.triangles)[i+5].Position = pixel.Vec{x + same_x, y + same_y}
			for n := 0; n < 6; n++ {
				(*o.triangles)[i+n].Color = pixel.RGBA{r, g, b, a}
			}

			i += 6
		}
	}

	// Reset the greedy bit
	for x := 0.0; x < o.width; x++ {
		for y := 0.0; y < o.height; y++ {
			pos = int(x*o.size + y)
			if o.pixels[pos] != 0 {
				o.pixels[pos] |= 0x00000080
			}
		}
	}
	o.triangles.SetLen(draw * 6)
	o.batch.Dirty()
	o.dirty = false
	o.canvas.Clear(pixel.RGBA{0, 0, 0, 0})
	o.batch.Draw(o.canvas)
}

//=============================================================
//
//  Function to implement Entity interface
//
//=============================================================
//=============================================================
//
//=============================================================
func (o *object) hit(x_, y_, vx, vy float64, power int) bool {
	o.explode()
	return true

	x := int(math.Abs(float64(o.bounds.X - x_)))
	y := int(math.Abs(float64(o.bounds.Y - y_)))

	pow := power * power
	for rx := x - power; rx <= x+power; rx++ {
		xx := (rx - x) * (rx - x)
		for ry := y - power; ry <= y+power; ry++ {
			if ry < 0 {
				continue
			}
			val := (ry-y)*(ry-y) + xx
			if val < pow {
				pos := int(o.size)*rx + ry
				if pos >= 0 && pos < int(o.size*o.size) {
					if o.pixels[pos] != 0 {
						global.gParticleEngine.newParticle(
							particle{
								x:           float64(x_),
								y:           float64(y_),
								size:        1,
								restitution: -0.1 - global.gRand.randFloat()/4,
								life:        wParticleDefaultLife,
								fx:          10,
								fy:          10,
								vx:          vx,
								vy:          float64(5 - global.gRand.rand()),
								mass:        1,
								pType:       particleRegular,
								color:       o.pixels[pos],
								static:      true,
							})
						o.pixels[pos] = 0
					}
				}
			}
		}
	}
	o.dirty = true
	//o.build()
	return true
}

//=============================================================
//
//=============================================================
func (o *object) isFree() bool {
	if o.owner == nil {
		return true
	}
	return false
}

//=============================================================
//
//=============================================================
func (o *object) explode() {
	o.active = false
	for x := 0.0; x < o.width; x++ {
		for y := 0.0; y < o.height; y++ {
			p := o.pixels[int(x*o.size+y)]
			if p == 0 {
				continue
			}

			size := o.scale
			if o.scale < 0.5 {
				size = 0.5
			}
			global.gParticleEngine.newParticle(
				particle{
					x:           o.bounds.X + float64(x)*o.scale,
					y:           o.bounds.Y + float64(y)*o.scale,
					size:        size,
					restitution: -0.1 - global.gRand.randFloat()/4,
					life:        wParticleDefaultLife,
					fx:          float64(10 - global.gRand.rand()),
					fy:          float64(10 - global.gRand.rand()),
					vx:          float64(5 - global.gRand.rand()),
					vy:          float64(5 - global.gRand.rand()),
					mass:        1,
					pType:       particleRegular,
					color:       p,
					static:      true,
				})
			o.pixels[int(x*o.size+y)] = 0
		}
	}
	o.owner = nil
	global.gWorld.qt.Remove(o.bounds)
}

//=============================================================
//
//=============================================================
func (o *object) pickup() {
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
func (o *object) setPosition(x, y float64) {
	o.bounds.X = x
	o.bounds.Y = y
}

//=============================================================
//
//=============================================================
func (o *object) setOwner(e Entity) {
	o.owner = e
}

//=============================================================
//
//=============================================================
func (o *object) removeOwner() {
	o.owner = nil
	o.fx = 1
	o.fy = 1
	o.bounces = 4
}

//=============================================================
// Get bounds
//=============================================================
func (o *object) getBounds() *Bounds {
	return o.bounds
}

//=============================================================
// Physics
//=============================================================
func (o *object) physics(dt float64) {
	// Only apply physics if not hold by an entity
	if o.owner != nil {
		return
	}

	if global.gWorld.IsWall(o.bounds.X, o.bounds.Y) {
		o.bounces++
		if o.bounces <= 4 {
			if o.vy < 0 {
				o.vy *= o.restitution
			} else {
				if o.vx > 0 {
					o.vx *= -o.restitution
					o.vy *= -o.restitution
				} else if o.vx < 0 {
					o.vx *= -o.restitution
					o.vy *= -o.restitution
				}
			}
		} else {
			o.fx = 0
			o.fy = 0
		}
	}
	ax := o.fx * dt * o.vx * o.mass
	ay := o.fy * dt * o.vy * o.mass
	o.bounds.X += ax
	o.bounds.Y += ay

	o.vy -= dt * o.fy

	if o.fx > 0 {
		o.fx -= dt * global.gWorld.gravity * o.mass
	} else {
		o.fx = 0
	}
	o.falling = false
	if o.fy > 0 {
		o.fy -= dt * global.gWorld.gravity * o.mass
		o.falling = true
	}

}

//=============================================================
//
//=============================================================
func (o *object) draw(dt, elapsed float64) {
	if !o.active {
		return
	}

	if o.dirty {
		o.build()
	}

	// This should be kept in weapon somehow...
	// But currently weapon has no draw.
	o.reloadTime += dt

	if o.owner == nil {
		o.physics(dt)
		offset := 0.0
		if !(o.falling || !o.animateIdle) {
			// Animate up/down
			offset = 5 + math.Sin(o.reloadTime)*3
		}
		o.canvas.Draw(global.gWin, pixel.IM.ScaledXY(pixel.ZV, pixel.V(o.scale, o.scale)).Moved(pixel.V(o.bounds.X+o.bounds.Width/2, offset+o.bounds.Y+o.bounds.Height/2)))
		o.unStuck(dt)
	} else {
		owner := o.owner.(*mob)
		o.canvas.Draw(global.gWin, pixel.IM.ScaledXY(pixel.ZV, pixel.V(o.scale*owner.dir, o.scale)).
			Moved(pixel.V(owner.bounds.X+owner.bounds.Width/2, owner.bounds.Y+owner.bounds.Height/2-2)).
			Rotated(pixel.Vec{o.bounds.X + o.bounds.Width/2, o.bounds.Y + o.bounds.Height/2}, o.rotation*owner.dir))
		// Update oect positions based on mob
		o.bounds.X = owner.bounds.X
		o.bounds.Y = owner.bounds.Y

	}
}

//=============================================================
// Unstuck the object if stuck.
//=============================================================
func (o *object) unStuck(dt float64) {
	bottom := false
	top := false
	offset := 1.0

	// Check bottom pixels
	for x := o.bounds.X; x < o.bounds.X+o.bounds.Width; x += 2 {
		if global.gWorld.IsRegular(x, o.bounds.Y+offset) {
			bottom = true
			break
		}
	}

	//Check top pixels
	for x := o.bounds.X; x < o.bounds.X+o.bounds.Width; x += 2 {
		if global.gWorld.IsRegular(x, o.bounds.Y+o.bounds.Height-offset) {
			top = true
			break
		}
	}

	if bottom {
		o.bounds.Y += 10 * o.mass * dt
	} else if top {
		o.bounds.Y -= 10 * o.mass * dt
	}
}
