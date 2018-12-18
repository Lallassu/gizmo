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
	"math"
)

type object struct {
	static      bool
	name        string
	textureFile string
	sprite      *sprite
	img         image.Image
	model       *imdraw.IMDraw
	canvas      *pixelgl.Canvas
	bounds      *Bounds
	mass        float64
	restitution float64
	height      int
	width       int
	size        int
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

	if !o.static {
		tmp_w := 0.0
		tmp_h := 0.0
		tmp_s := 0.0

		o.img, tmp_w, tmp_h, tmp_s = loadTexture(o.textureFile)
		o.width = int(tmp_w)
		o.height = int(tmp_h)
		o.size = int(tmp_s)

		// Initiate bounds for qt
		o.bounds = &Bounds{
			X:      x_,
			Y:      y_,
			Width:  float64(o.width) * o.scale,
			Height: float64(o.height) * o.scale,
			entity: Entity(o),
		}

		o.pixels = make([]uint32, o.size*o.size)

		for x := 0; x < o.width; x++ {
			for y := 0; y < o.height; y++ {
				r, g, b, a := o.img.At(x, o.height-y).RGBA()
				o.pixels[x*o.size+y] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | a&0xFF
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
	} else {
		// o.sprite = &sprite{name: o.name, scale: o.scale}
		// o.width, o.height = global.gTextures.spriteInfo(o.name)
		// Debug("SPRITE:", o.width, o.height)

		// // Load sprite included in batch
		// o.bounds = &Bounds{
		// 	X:      x_,
		// 	Y:      y_,
		// 	Width:  float64(o.width) * o.scale,
		// 	Height: float64(o.height) * o.scale,
		// 	entity: Entity(o),
		// }

		// global.gTextures.addObject(o.sprite)
	}

	// Add object to QT
	global.gWorld.AddObject(o.bounds)
}

//=============================================================
// Build
//=============================================================
func (o *object) build() {
	if !o.static {
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
	if !o.static {
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
					pos := o.size*rx + ry
					if pos >= 0 && pos < o.size*o.size {
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
		o.build()
	} else {
		// TBD: Keep track on hits before explode?
		o.explode()
	}
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
	// if !o.static {
	// 	if o.objectType == objectCrate {
	// 		o.active = false
	// 		for x := 0; x < o.width; x++ {
	// 			for y := 0; y < o.height; y++ {
	// 				p := o.pixels[x*o.size+y]
	// 				o.pixels[x*o.size+y] = 0

	// 				global.gParticleEngine.newParticle(
	// 					particle{
	// 						x:           o.bounds.X + float64(x),
	// 						y:           o.bounds.Y + float64(y),
	// 						size:        1,
	// 						restitution: -0.1 - global.gRand.randFloat()/4,
	// 						life:        wParticleDefaultLife,
	// 						fx:          10,
	// 						fy:          10,
	// 						vx:          float64(5 - global.gRand.rand()),
	// 						vy:          float64(5 - global.gRand.rand()),
	// 						mass:        1,
	// 						pType:       particleRegular,
	// 						color:       p,
	// 						static:      true,
	// 					})
	// 			}
	// 		}
	// 	}
	// } else {

	// }
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

	o.reloadTime += dt

	if o.owner == nil {
		o.physics(dt)
		if !o.static {
			if o.falling || !o.animateIdle {
				o.canvas.Draw(global.gWin, pixel.IM.ScaledXY(pixel.ZV, pixel.V(o.scale, o.scale)).Moved(pixel.V(o.bounds.X+o.bounds.Width/2, o.bounds.Y+o.bounds.Height/2)))
			} else {
				// Animate up/down
				offset := 5 + math.Sin(o.reloadTime)*3
				o.canvas.Draw(global.gWin, pixel.IM.ScaledXY(pixel.ZV, pixel.V(o.scale, o.scale)).Moved(pixel.V(o.bounds.X+o.bounds.Width/2, offset+o.bounds.Y+o.bounds.Height/2)))
			}
		} else {
			o.sprite.pos = pixel.Vec{o.bounds.X, o.bounds.Y}
		}
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
// Unstuck the objet if stuck.
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
