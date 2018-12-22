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
	"math"
)

type object struct {
	phys
	graphics
	name  string
	oType itemType
	//	textureFile string
	//height      float64
	//width       float64
	//size        float64
	//pixels      []uint32
	//	scale       float64
	owner Entity
	//active      bool
	reloadTime  float64
	animateIdle bool
	//	sprite      *sprite
	//	canvas      *pixelgl.Canvas

	// batch     *pixel.Batch
	// triangles *pixel.TrianglesData
	// dirty     bool
}

//=============================================================
//
//=============================================================
func (o *object) create(x, y float64) {
	o.mass = 5
	//o.active = true
	o.animateIdle = false

	o.createGfx(x, y)
	o.createPhys(x, y, o.frameWidth, o.frameHeight)

	//	o.dirty = true
	//	o.triangles = pixel.MakeTrianglesData(100)
	//	o.batch = pixel.NewBatch(o.triangles, nil)

	// var img image.Image
	// img, o.width, o.height, o.size = loadTexture(o.textureFile)

	// // Initiate bounds for qt
	// o.bounds = &Bounds{
	// 	X:      x_,
	// 	Y:      y_,
	// 	Width:  o.width * o.scale,
	// 	Height: o.height * o.scale,
	// 	entity: Entity(o),
	// }

	// o.pixels = make([]uint32, int(o.size*o.size))

	// for x := 0.0; x < o.width; x++ {
	// 	for y := 0.0; y < o.height; y++ {
	// 		r, g, b, a := img.At(int(x), int(o.height-y)).RGBA()
	// 		if r == 0 && g == 0 && b == 0 && a == 0 {
	// 			continue
	// 		}
	// 		o.pixels[int(x*o.size+y)] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | a&0xFF
	// 	}
	// }

	// o.canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(o.width), float64(o.height)))

	// // build initial
	// o.build()

	// // Add object to QT
	// global.gWorld.AddObject(o.bounds)
}

//=============================================================
// Build
//=============================================================
//func (o *object) build() {
//	i := 0
//	rc := uint32(0)
//	gc := uint32(0)
//	bc := uint32(0)
//	p2 := uint32(0)
//	r1 := uint32(0)
//	g1 := uint32(0)
//	b1 := uint32(0)
//	draw := 0
//	same_x := 1.0
//	same_y := 1.0
//	pos := 0
//
//	for x := 0.0; x < o.width; x++ {
//		for y := 0.0; y < o.height; y++ {
//			p := o.pixels[int(x*o.size+y)]
//			// Skip visisted or empty
//			if p == 0 || p&0xFF>>7 == 0 {
//				continue
//			}
//			rc = p >> 24 & 0xFF
//			gc = p >> 16 & 0xFF
//			bc = p >> 8 & 0xFF
//			same_x = 1.0
//			same_y = 1.0
//
//			for l := x + 1; l < o.width; l++ {
//				// Check color
//				pos = int(l*o.size + y)
//				p2 = o.pixels[pos]
//				if p2 == 0 {
//					break
//				}
//				r1 = p2 >> 24 & 0xFF
//				g1 = p2 >> 16 & 0xFF
//				b1 = p2 >> 8 & 0xFF
//
//				if r1 == rc && g1 == gc && b1 == bc && ((p2&0xFF)>>7) == 1 {
//					// Same color and not yet visited!
//					o.pixels[pos] &= 0xFFFFFF7F
//					same_x++
//					new_y := 1.0
//					for k := y; k < o.height; k++ {
//						pos = int(l*o.size + k)
//						p2 = o.pixels[pos]
//						if p2 == 0 {
//							break
//						}
//						r1 = p2 >> 24 & 0xFF
//						g1 = p2 >> 16 & 0xFF
//						b1 = p2 >> 8 & 0xFF
//
//						if r1 == rc && g1 == gc && b1 == bc && ((p2&0xFF)>>7) == 1 {
//							o.pixels[pos] &= 0xFFFFFF7F
//							new_y++
//						} else {
//							break
//						}
//					}
//					if new_y < same_y {
//						break
//					} else {
//						same_y = new_y
//					}
//				} else {
//					break
//				}
//			}
//
//			draw++
//
//			// Convert to decimal
//			r := float64(p>>24&0xFF) / 255.0
//			g := float64(p>>16&0xFF) / 255.0
//			b := float64(p>>8&0xFF) / 255.0
//			a := float64(p&0xFF) / 255.0
//
//			// Increase length of triangles if we need to draw more than we had before.
//			// Add a buffer so we can skip a few increments.
//			if draw*6 >= len(*o.triangles) {
//				o.triangles.SetLen(draw*6 + 10)
//			}
//
//			// Size of triangle is given by how large the greedy algorithm found out.
//			(*o.triangles)[i].Position = pixel.Vec{x, y}
//			(*o.triangles)[i+1].Position = pixel.Vec{x + same_x, y}
//			(*o.triangles)[i+2].Position = pixel.Vec{x + same_x, y + same_y}
//			(*o.triangles)[i+3].Position = pixel.Vec{x, y}
//			(*o.triangles)[i+4].Position = pixel.Vec{x, y + same_y}
//			(*o.triangles)[i+5].Position = pixel.Vec{x + same_x, y + same_y}
//			for n := 0; n < 6; n++ {
//				(*o.triangles)[i+n].Color = pixel.RGBA{r, g, b, a}
//			}
//
//			i += 6
//		}
//	}
//
//	// Reset the greedy bit
//	for x := 0.0; x < o.width; x++ {
//		for y := 0.0; y < o.height; y++ {
//			pos = int(x*o.size + y)
//			if o.pixels[pos] != 0 {
//				o.pixels[pos] |= 0x00000080
//			}
//		}
//	}
//	o.triangles.SetLen(draw * 6)
//	o.batch.Dirty()
//	o.dirty = false
//	o.canvas.Clear(pixel.RGBA{0, 0, 0, 0})
//	o.batch.Draw(o.canvas)
//}

//=============================================================
//
//  Function to implement Entity interface
//
//=============================================================
//=============================================================
//
//=============================================================
func (o *object) hit(x_, y_, vx, vy float64, power int) {
	o.explode()
	return

	// x := int(math.Abs(float64(o.bounds.X - x_)))
	// y := int(math.Abs(float64(o.bounds.Y - y_)))

	// pow := power * power
	// for rx := x - power; rx <= x+power; rx++ {
	// 	xx := (rx - x) * (rx - x)
	// 	for ry := y - power; ry <= y+power; ry++ {
	// 		if ry < 0 {
	// 			continue
	// 		}
	// 		val := (ry-y)*(ry-y) + xx
	// 		if val < pow {
	// 			pos := int(o.size)*rx + ry
	// 			if pos >= 0 && pos < int(o.size*o.size) {
	// 				if o.pixels[pos] != 0 {
	// 					global.gParticleEngine.newParticle(
	// 						particle{
	// 							x:           float64(x_),
	// 							y:           float64(y_),
	// 							size:        1,
	// 							restitution: -0.1 - global.gRand.randFloat()/4,
	// 							life:        wParticleDefaultLife,
	// 							fx:          10,
	// 							fy:          10,
	// 							vx:          vx,
	// 							vy:          float64(5 - global.gRand.rand()),
	// 							mass:        1,
	// 							pType:       particleRegular,
	// 							color:       o.pixels[pos],
	// 							static:      true,
	// 						})
	// 					o.pixels[pos] = 0
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	// o.dirty = true
	// //o.build()
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
	o.explodeGfx(o.bounds.X, o.bounds.Y)
	global.gWorld.qt.Remove(o.bounds)
	//o.active = false
	//for x := 0.0; x < o.width; x++ {
	//	for y := 0.0; y < o.height; y++ {
	//		p := o.pixels[int(x*o.size+y)]
	//		if p == 0 {
	//			continue
	//		}

	//		if global.gRand.randFloat() < 0.1 {
	//			size := o.scale
	//			if o.scale < 0.5 {
	//				size = 0.5
	//			}
	//			global.gParticleEngine.newParticle(
	//				particle{
	//					x:           o.bounds.X + float64(x)*o.scale,
	//					y:           o.bounds.Y + float64(y)*o.scale,
	//					size:        size,
	//					restitution: -0.1 - global.gRand.randFloat()/4,
	//					life:        wParticleDefaultLife,
	//					fx:          float64(10 - global.gRand.rand()),
	//					fy:          float64(10 - global.gRand.rand()),
	//					vx:          float64(5 - global.gRand.rand()),
	//					vy:          float64(5 - global.gRand.rand()),
	//					mass:        1,
	//					pType:       particleRegular,
	//					color:       p,
	//					static:      true,
	//				})
	//			o.pixels[int(x*o.size+y)] = 0
	//		}
	//	}
	//}
	//if o.owner != nil {
	//	o.owner.throw()
	//}
}

//=============================================================
//
//=============================================================
func (o *object) throw() {

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
	//o.fx += o.owner.(*mob).dir * 10
	//o.vx += o.owner.(*mob).dir * 10

	o.force.X = math.Abs(o.owner.(*mob).velo.X) + 5
	o.force.Y = 5
	o.velo.Y = 5
	o.velo.X = o.owner.(*mob).velo.X / 2
	o.owner = nil
}

//=============================================================
// Get bounds
//=============================================================
func (o *object) getBounds() *Bounds {
	return o.bounds
}

//=============================================================
//
//=============================================================
func (o *object) draw(dt, elapsed float64) {
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
		o.canvas.Draw(global.gWin, pixel.IM.ScaledXY(pixel.ZV, pixel.V(o.scale, o.scale)).Moved(pixel.V(o.bounds.X+o.bounds.Width/2, offset+o.bounds.Y+o.bounds.Height/2)).Rotated(pixel.V(o.bounds.X+o.bounds.Width/2, o.bounds.Y+o.bounds.Height/2), o.rotation))
	} else {
		owner := o.owner.(*mob)
		offset := 0.0
		switch o.bounds.entity.(type) {
		case *item:
			offset = 10.0
		}
		o.canvas.Draw(global.gWin, pixel.IM.ScaledXY(pixel.ZV, pixel.V(o.scale*owner.dir, o.scale)).
			Moved(pixel.V(owner.bounds.X+owner.bounds.Width/2, offset+owner.bounds.Y+owner.bounds.Height/2-2)).
			Rotated(pixel.Vec{o.bounds.X + o.bounds.Width/2, o.bounds.Y + o.bounds.Height/2}, o.rotation*owner.dir))
		// Update oect positions based on mob
		o.bounds.X = owner.bounds.X
		o.bounds.Y = owner.bounds.Y

	}
}
