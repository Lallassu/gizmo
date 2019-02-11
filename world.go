//=============================================================
// world.go
//-------------------------------------------------------------
// Keep control of map (all pixels)
// Destuction of map
// Additions to map
// Generation of map
// Map flood fill
// Quadtree for entities
//=============================================================
package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
)

//=============================================================
// World Structure
//=============================================================
type world struct {
	width      int
	height     int
	size       int
	qt         *Quadtree
	pixels     []uint32
	currentMap mapType
	gravity    float64
	bgSprite   *pixel.Sprite
	doors      []pixel.Vec
}

//=============================================================
//=============================================================
// World Public Functions
//=============================================================
//=============================================================

//=============================================================
// Initialize world first time.
//=============================================================
func (w *world) Init() {
	w.qt = &Quadtree{
		Bounds:     Bounds{X: 0, Y: 0, Width: float64(w.width), Height: float64(w.height)},
		MaxObjects: 4,
		MaxLevels:  8,
		Level:      0,
	}
	w.gravity = wGravity
	w.doors = make([]pixel.Vec, 0)
}

//=============================================================
// New Map
//=============================================================
func (w *world) NewMap(width, height, size float64) {
	w.qt.Clear()

	w.width = int(width)
	w.height = int(height)
	w.size = int(size)
	w.pixels = make([]uint32, int(size))

	// FG Chunks
	for x := 0; x < w.width; x += wPixelsPerChunk {
		for y := 0; y < w.height; y += wPixelsPerChunk {
			c := &chunk{cType: fgChunk}
			c.create(float64(x), float64(y), wPixelsPerChunk)
			w.qt.Insert(c.bounds)
		}
	}
}

func (w *world) buildAllChunks() {
	// Build all chunks first time.
	for _, v := range w.qt.RetrieveIntersections(&Bounds{X: 0, Y: 0, Width: float64(w.width), Height: float64(w.height)}) {
		v.entity.draw(-1, 0)
	}

	// One sprite for whole bg
	c := &chunk{cType: bgChunk}
	c.create(float64(0), float64(0), w.width)
	c.build()
	w.bgSprite = c.sprite
}

func (w *world) fitInWorld(size int) (pixel.Vec, bool) {
	posX := rand.Intn(w.width - 1)
	posY := rand.Intn(w.height - 1)
	for x := posX - size/2; x < posX+size/2; x++ {
		for y := posY - size/2; y < posY+size/2; y++ {
			pos := x*w.width + y
			if pos < len(w.pixels) && pos >= 0 {
				p := w.pixels[pos]
				if p&0xFF != wBackground8 {
					return pixel.Vec{}, false
				}
			}
		}
	}
	return pixel.Vec{X: float64(posX), Y: float64(posY)}, true
}

//func loadPicture(path string) (pixel.Picture, error) {
//	file, err := os.Open(path)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//	img, _, err := image.Decode(file)
//	if err != nil {
//		return nil, err
//	}
//	return pixel.PictureDataFromImage(img), nil
//}

//=============================================================
// Add object to world (QT)
//=============================================================
func (w *world) AddObject(obj *Bounds) {
	w.qt.Insert(obj)
}

//=============================================================
// Remove object from world (QT)
//=============================================================
func (w *world) RemoveObject(obj entity) {

}

//=============================================================
// Check if pixel is a background
//=============================================================
func (w *world) IsBackground(posX, posY float64) bool {
	x := int(posX)
	y := int(posY)
	pos := w.width*x + y
	if pos < w.size && pos >= 0 {
		if w.pixels[pos]&0xFF == wBackground8 ||
			w.pixels[pos]&0xFF == wBackgroundNew8 {
			return true
		}
	}
	return false
}

//=============================================================
// Check if pixel is a shadow
//=============================================================
func (w *world) IsShadow(posX, posY float64) bool {
	x := int(posX)
	y := int(posY)
	pos := w.width*x + y
	if pos < w.width*w.height && pos >= 0 {
		if w.pixels[pos]&0xFF == wShadow8 {
			return true
		}
	}
	return false
}

//=============================================================
// Check if pixel is regular
//=============================================================
func (w *world) IsRegular(posX, posY float64) bool {
	x := int(posX)
	y := int(posY)
	pos := w.width*x + y
	if pos < w.width*w.height && pos >= 0 {
		if w.pixels[pos]&0xFF == 0xFF {
			return true
		}
	}
	return false
}

//=============================================================
// Check if it's a wall
//=============================================================
func (w *world) IsWall(posX, posY float64) bool {
	x := int(posX)
	y := int(posY)
	pos := w.width*x + y
	if pos < w.width*w.height && pos >= 0 {
		if w.pixels[pos] != 0 && w.pixels[pos]&0xFF != wBackgroundNew8 && w.pixels[pos]&0xFF != wBackground8 && w.pixels[pos]&0xFF != wShadow8 && w.pixels[pos]&0xFF != wLadder8 {
			return true
		}
	}
	return false
}

//=============================================================
// Check if it's a ladder
//=============================================================
func (w *world) IsLadder(posX, posY float64) bool {
	x := int(posX)
	y := int(posY)
	pos := w.width*x + y
	if pos < w.width*w.height && pos >= 0 {
		if w.pixels[pos]&0xFF == wLadder8 {
			return true
		}
	}
	return false

}

//=============================================================
// Check if pixel exists
//=============================================================
func (w *world) PixelExists(x, y float64) bool {
	return true
}

//=============================================================
// Get color of the specified pixel
// Return -1 if not exist
//=============================================================
func (w *world) PixelColor(x, y float64) int32 {
	if x < 0 || y < 0 || x >= float64(w.width) || y >= float64(w.height) {
		return -1
	}
	return int32(w.pixels[uint32(int(x)*w.width+int(y))])
}

//=============================================================
// Draw
//=============================================================
func (w *world) Draw(dt, elapsed float64) {
	w.bgSprite.Draw(global.gWin, pixel.IM.Moved(pixel.V(float64(w.width)/2, float64(w.height)/2)))

	// Draw objects in QT around player position only.
	pos := pixel.Vec{X: 0, Y: 0}
	if global.gCamera.follow != nil {
		pos = global.gCamera.follow.getPosition()
	}
	for _, v := range w.qt.RetrieveIntersections(&Bounds{X: pos.X - wViewMax/2, Y: pos.Y - wViewMax/2, Width: wViewMax, Height: wViewMax}) {
		v.entity.draw(dt, elapsed)
	}
}

//=============================================================
// Add pixel with color (replace if already exists)
//=============================================================
func (w *world) AddPixel(x, y int, color uint32) {
	pos := w.width*x + y
	if pos < w.width*w.height && pos >= 0 {
		w.pixels[pos] = color
		w.markChunkDirty(x, y)
	}
}

//=============================================================
// Set pixel without rebuilding chunk
//=============================================================
func (w *world) SetPixel(x, y int, color uint32) {
	pos := w.width*x + y
	if pos < w.width*w.height && pos >= 0 {
		w.pixels[pos] = color
	}
}

//=============================================================
// Remove a pixel from the world map
//=============================================================
func (w *world) RemovePixel(x, y int) {
	pos := w.width*x + y
	if pos < w.width*w.height && pos >= 0 {
		if w.pixels[pos]&0xFF == wStaticColor8 ||
			w.pixels[pos]&0xFF == wBackground8 ||
			w.pixels[pos]&0xFF == wBackgroundNew8 ||
			w.pixels[pos]&0xFF == wLadder8 ||
			w.pixels[pos]&0xFF == wShadow8 {
			return
		}

		// Remove shadow
		for i := 0; i < wShadowLength; i++ {
			pos2 := (x+i)*w.width + y - i
			if pos2 < w.width*w.height && pos2 >= 0 {
				w.removeShadow(x+i, y-i)
			}
		}

		// Particle
		if w.IsRegular(float64(x), float64(y)) {
			global.gParticleEngine.newParticle(
				particle{
					x:           float64(x),
					y:           float64(y),
					size:        1,
					restitution: -0.1 - global.gRand.randFloat()/4,
					life:        wParticleDefaultLife,
					fx:          10 + float64(5-global.gRand.rand()),
					fy:          10 + float64(5-global.gRand.rand()),
					vx:          float64(5 - global.gRand.rand()),
					vy:          float64(5 - global.gRand.rand()),
					mass:        1,
					pType:       particleRegular,
					color:       w.pixels[pos],
					static:      true,
				})
		}

		// Set bg pixel.
		if w.pixels[pos] != 0 {
			v := global.gMapColor.backgroundSoft
			v &= wBackgroundNew32
			w.pixels[pos] = v
		}
		w.markChunkDirty(x, y)
	}
}

//=============================================================
// Explode in world
// Also hits objects in the world.
//=============================================================
func (w *world) Explode(posX, posY float64, power int) {
	//	global.gSounds.play("explosion")

	global.gCamera.shake(pixel.V(posX, posY), power)
	x := int(posX)
	y := int(posY)
	pow := power * power
	ff := make([]pixel.Vec, 50)
	for rx := x - power; rx <= x+power; rx++ {
		vx := (rx - x) * (rx - x)
		for ry := y - power; ry <= y+power; ry++ {
			if ry < 0 {
				continue
			}
			val := (ry-y)*(ry-y) + vx
			if val < pow {
				w.RemovePixel(rx, ry)
				//w.ObjectHit(float64(rx), float64(ry))
				//for _, v := range w.qt.RetrieveIntersections(&Bounds{X: float64(rx), Y: float64(ry), Width: 1, Height: 1}) {
				//	v.entity.hit(x_, y_)
				//}
			} else {
				ff = append(ff, pixel.Vec{X: float64(rx), Y: float64(ry)})
			}
		}
	}

	// Retrieve with power * 4 for shockwave effect.
	for _, v := range w.qt.RetrieveIntersections(&Bounds{X: posX - float64(power*4), Y: posY - float64(power*4), Width: float64(power * 8), Height: float64(power * 8)}) {
		v.entity.hit(posX, posY, 0, 0, power)
	}

	// Add shadows
	for n := 0; n < len(ff); n++ {
		ffx := int(ff[n].X)
		ffy := int(ff[n].Y)
		pp := ffx*w.width + ffy
		if pp >= 0 && pp < w.width*w.height {
			if w.pixels[pp]&0xFF == 0xFF {
				for i := 0; i < wShadowLength; i++ {
					pos2 := (ffx+i)*w.width + ffy - i
					if pos2 < w.width*w.height && pos2 >= 0 {
						if w.pixels[pos2]&0xFF == wBackground8 ||
							w.pixels[pos2]&0xFF == wBackgroundNew8 {
							w.addShadow(ffx+i, ffy-i)
						}
					}
				}
			}
		}
	}
	l := light{}
	if power > 10 {
		l.create(posX, posY, 300, 360, float64(2*power), pixel.RGBA{R: 0.8, G: 0.3, B: 0, A: 0.1}, true, 0.1)
	}

	// Floodfill
	// pixels := make([]Vec, 0)
	// for i := 0; i < len(ff); i++ {
	// 	pixels = append(pixels, w.FloodFill(ff[i].X, ff[i].Y)...)
	// }

	// for i := 0; i < len(pixels); i++ {
	// 	w.UnMarkPixelVisited(pixels[i].X, pixels[i].Y)
	// }}
}

//=============================================================
//=============================================================
// World Internal Functions
//=============================================================
//=============================================================

//=============================================================
// Flood fill in map
//=============================================================
func (w *world) floodFill(x, y int) {

}

//=============================================================
// Remove shadows from map on given position
//=============================================================
func (w *world) removeShadow(x, y int) {
	pos := w.width*x + y
	if pos < w.size && pos >= 0 {
		if w.pixels[pos]&0xFF == wShadow8 {
			r := uint32(math.Floor(float64(w.pixels[pos]>>24&0xFF) * wShadowDepth))
			g := uint32(math.Floor(float64(w.pixels[pos]>>16&0xFF) * wShadowDepth))
			b := uint32(math.Floor(float64(w.pixels[pos]>>8&0xFF) * wShadowDepth))
			w.pixels[pos] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | wBackgroundNew8
			w.markChunkDirty(x, y)
		}
	}
}

//=============================================================
// Add shadows to map on given position
//=============================================================
func (w *world) addShadow(x, y int) {
	pos := w.width*x + y
	if pos < w.width*w.height && pos >= 0 {
		if w.pixels[pos]&0xFF != wShadow8 && w.pixels[pos]&0xFF != 0xFF {
			r := uint32(math.Ceil(float64(w.pixels[pos]>>24&0xFF) / wShadowDepth))
			g := uint32(math.Ceil(float64(w.pixels[pos]>>16&0xFF) / wShadowDepth))
			b := uint32(math.Ceil(float64(w.pixels[pos]>>8&0xFF) / wShadowDepth))
			w.pixels[pos] = r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | wShadow8
			w.markChunkDirty(x, y)
		}
	}
}

//=============================================================
// Mark chunk as dirty to rebuild it
//=============================================================
func (w *world) markChunkDirty(x, y int) {
	// Get all chunks in this area.
	for _, v := range w.qt.RetrieveIntersections(&Bounds{X: float64(x), Y: float64(y), Width: 3, Height: 3}) {
		switch item := v.entity.(type) {
		case *chunk:
			item.dirty = true
		}
	}
}
