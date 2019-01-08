//=============================================================
// pcg.go
//-------------------------------------------------------------
// Procedurally Generated Stuff
//=============================================================
package main

import (
	_ "math"
	"math/rand"
)

type pcg struct {
	floorHeight int
	incr        float64
	floorCnt    int
	bgPlateCnt  int
	doorCnt     int
	lampCnt     int
	airCnt      int
}

func (p *pcg) Flower(x, y int) {
	//color := global.gWorld.coloring.getFlower()
	stem := rand.Intn(10)
	for i := 0; i < stem; i++ {
		global.gWorld.SetPixel(x, y+i, 0xFF00FFFF)
		global.gWorld.addShadow(x+1, y+i-1)
	}
}

func (p *pcg) MetalCornerDown(x, y int, left bool) {
	r := 0
	g := 0
	b := 0
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			if i == 0 || j == 0 {
				r = 0x2b
				g = 0x37
				b = 0x31
			} else if i == 1 || j == 1 {
				r = 0x22
				g = 0x2c
				b = 0x27
			} else if i == 2 || j == 2 {
				r = 0x44
				g = 0x55
				b = 0x49
			} else if i == 3 || j == 3 {
				r = 0x3a
				g = 0x4e
				b = 0x42
			} else if i == 4 || j == 4 {
				r = 0x44
				g = 0x55
				b = 0x49
			} else if i == 5 || j == 5 {
				r = 0x1e
				g = 0x2c
				b = 0x26
			}
			if left {
				global.gWorld.SetPixel(
					x+i,
					y+j,
					uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
				)
			} else {
				global.gWorld.SetPixel(
					x-i,
					y+j,
					uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
				)
			}
		}
	}
}

func (p *pcg) MetalCornerUp(x, y int, left bool) {
	r := 0
	g := 0
	b := 0
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			if i == 0 || j == 0 {
				r = 0x2b
				g = 0x37
				b = 0x31
			} else if i == 1 || j == 1 {
				r = 0x22
				g = 0x2c
				b = 0x27
			} else if i == 2 || j == 2 {
				r = 0x44
				g = 0x55
				b = 0x49
			} else if i == 3 || j == 3 {
				r = 0x3a
				g = 0x4e
				b = 0x42
			} else if i == 4 || j == 4 {
				r = 0x44
				g = 0x55
				b = 0x49
			} else if i == 5 || j == 5 {
				r = 0x1e
				g = 0x2c
				b = 0x26
			}
			if left {
				global.gWorld.SetPixel(
					x-i,
					y-j,
					uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
				)
			} else {
				global.gWorld.SetPixel(
					x+i,
					y-j,
					uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
				)
			}
		}
	}
}

func (p *pcg) MetalCornerRight(x, y int, left bool) {
	r := 0
	g := 0
	b := 0
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			if i == 5 || j == 5 || j == 1 && i == 1 || j == 1 && i == 0 || j == 0 && i == 1 {
				r = 0x1e
				g = 0x2c
				b = 0x26
			} else if i == 0 && j == 0 {
				r = 0x2b
				g = 0x37
				b = 0x31
			} else if i == 4 || j == 4 || j == 2 && i == 1 || i == 2 && j == 1 || j == 2 && i == 2 || j == 2 && i == 0 || j == 0 && i == 2 {
				r = 0x44
				g = 0x55
				b = 0x49
			} else if i == 3 || j == 3 {
				r = 0x3a
				g = 0x4e
				b = 0x42
			}
			if left {
				global.gWorld.SetPixel(
					x+i,
					y+j,
					uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
				)
			} else {
				global.gWorld.SetPixel(
					x-i,
					y+j,
					uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
				)
			}
		}
	}
}

func (p *pcg) MetalCornerLeft(x, y int, left bool) {
	r := 0
	g := 0
	b := 0
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			if i == 5 || j == 5 || j == 1 && i == 1 || j == 1 && i == 0 || j == 0 && i == 1 {
				r = 0x1e
				g = 0x2c
				b = 0x26
			} else if i == 0 && j == 0 {
				r = 0x2b
				g = 0x37
				b = 0x31
			} else if i == 4 || j == 4 || j == 2 && i == 1 || i == 2 && j == 1 || j == 2 && i == 2 || j == 2 && i == 0 || j == 0 && i == 2 {
				r = 0x44
				g = 0x55
				b = 0x49
			} else if i == 3 || j == 3 {
				r = 0x3a
				g = 0x4e
				b = 0x42
			}
			if left {
				global.gWorld.SetPixel(
					x+i,
					y-j,
					uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
				)
			} else {
				global.gWorld.SetPixel(
					x-i,
					y-j,
					uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
				)
			}
		}
	}
}

func (p *pcg) MetalWall(x, y int, leftSide bool) {
	r := 0
	g := 0
	b := 0
	x_ := 0
	for i := 0; i < 6; i++ {
		switch i {
		case 0:
			r = 0x2b
			g = 0x37
			b = 0x31
		case 1:
			r = 0x22
			g = 0x2c
			b = 0x27
		case 2, 4:
			r = 0x44
			g = 0x55
			b = 0x49
		case 3:
			r = 0x3a
			g = 0x4e
			b = 0x42
		case 5:
			r = 0x1e
			g = 0x2c
			b = 0x26
		}
		if leftSide {
			x_ = x + i
		} else {
			x_ = x - i
		}
		global.gWorld.SetPixel(
			x_,
			y,
			uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
		)
		// for n := 0; n < wShadowLength; n++ {
		// 	global.gWorld.addShadow(x+n, y-i-n)
		// }
	}
}

func (p *pcg) MetalFlat(x, y int, floor bool) {
	r := 0
	g := 0
	b := 0
	for i := 0; i < 6; i++ {
		switch i {
		case 0:
			r = 0x2b
			g = 0x37
			b = 0x31
		case 1:
			r = 0x22
			g = 0x2c
			b = 0x27
		case 2, 4:
			r = 0x44
			g = 0x55
			b = 0x49
		case 3:
			r = 0x3a
			g = 0x4e
			b = 0x42
		case 5:
			r = 0x1e
			g = 0x2c
			b = 0x26
		}
		if !floor {
			global.gWorld.SetPixel(
				x,
				y+i,
				uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
			)
		} else {
			global.gWorld.SetPixel(
				x,
				y-i,
				uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
			)
		}
	}
}

func (p *pcg) MetalFloor(x, y int) {
	r := 0
	g := 0
	b := 0
	p.floorCnt += 1
	for i := 0; i < 4; i++ {
		switch i {
		case 0:
			if p.floorCnt == 50 {
				r = 0x59
				g = 0x7e
				b = 0x6f
			} else {
				r = 0x7d
				g = 0xa6
				b = 0x91
			}
		case 1, 2:
			if p.floorCnt == 50 {
				r = 0x16
				g = 0x27
				b = 0x23
			} else {
				r = 0x59
				g = 0x7e
				b = 0x6f
			}
		case 3:
			r = 0x16
			g = 0x27
			b = 0x23
		}

		global.gWorld.SetPixel(
			x,
			y-i,
			uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
		)
	}
	if p.floorCnt == 50 {
		p.floorCnt = 0
	}
}

func (p *pcg) GenerateDoor(x, y int) bool {
	p.doorCnt++

	if p.doorCnt < 200 {
		return false
	}
	p.doorCnt = 0

	r := 0x73
	g := 0x80
	b := 0x62
	a := wBackground8
	doorType := rand.Float64()
	doorLight := rand.Float64()

	for i := 0; i < wDoorLen; i++ {
		for j := 0; j < wDoorHeight; j++ {
			// Frame
			if i == 0 || j == wDoorHeight-1 || i == wDoorLen-1 {
				r = 0x3a
				g = 0x3a
				b = 0x3a
			} else {
				r = 0x51
				g = 0x71
				b = 0x74
			}
			if i == wDoorLen-2 && j < wDoorHeight-1 || i == wDoorLen-3 && j < wDoorHeight-1 {
				r = 0x4c
				g = 0x4c
				b = 0x4c
			}
			if j == wDoorHeight-2 && i > 0 && i < wDoorLen-2 {
				r = 0x1b
				g = 0x1b
				b = 0x1b
			}
			// Shadow
			if j > wDoorHeight-5 || i < 3 {
				r /= 3
				g /= 3
				b /= 3
			}
			// Handle
			if j == wDoorHeight/3 && i > 1 && i < 5 {
				r = 0x52
				g = 0x67
				b = 0x69
			}
			// Handle
			if j == wDoorHeight/3-1 && i > 1 && i < 4 {
				r = 0x35
				g = 0x44
				b = 0x46
			}
			if j == wDoorHeight/3-1 && i > 1 && i < 6 {
				r = 0x35
				g = 0x44
				b = 0x46
			}

			// Middle of door
			if doorType > 0.5 {
				if j > wDoorHeight/2 && j < wDoorHeight-wDoorHeight/3 && i >= wDoorLen/3 && i < wDoorLen-wDoorLen/3 {
					// Jail door
					if i%2 == 0 {
						r = 0x00
						g = 0x00
						b = 0x00
					} else {
						r = 0x55
						g = 0x55
						b = 0x55
					}
					// Door window shadow
					if j == wDoorHeight/2-1 || j == wDoorHeight-wDoorHeight/3-1 || i == wDoorLen/3 {
						r = 0x00
						g = 0x00
						b = 0x00
					}
				}
			} else {
				if j > wDoorHeight/4 && j < wDoorHeight-wDoorHeight/3 && i >= wDoorLen/3 && i < wDoorLen-wDoorLen/3 {
					// Light
					if doorLight < 0.5 {
						r = 0xf6
						g = 0xab
						b = 0x34
					} else {
						r = 0x33 - i
						g = 0x33 - i
						b = 0x33 - i
					}
					// Door window shadow
					if j == wDoorHeight/4-1 || j == wDoorHeight-wDoorHeight/3-1 || i == wDoorLen/3 {
						r = 0x00
						g = 0x00
						b = 0x00
					}
				}
			}
			// Hinge
			if (j == wDoorHeight/4 || j == wDoorHeight/4-1 || j == wDoorHeight-(wDoorHeight/4-1) || j == wDoorHeight-wDoorHeight/4) && i == wDoorLen-3 {
				r = 0x4c
				g = 0x59
				b = 0x5c
			}
			// Hinge
			if (j == wDoorHeight/4 || j == wDoorHeight/4-1 || j == wDoorHeight-(wDoorHeight/4-1) || j == wDoorHeight-wDoorHeight/4) && i == wDoorLen-4 {
				r = 0x23
				g = 0x23
				b = 0x23
			}
			global.gWorld.SetPixel(
				x+i,
				y+j,
				uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|a),
			)
		}
	}
	return true
}

func (p *pcg) GenerateLamp(x, y int) {
	// p.lampCnt++
	// if p.lampCnt < 100 {
	// 	return
	// }
	// p.lampCnt = 0

	lr := 0x3b
	lg := 0x3a
	lb := 0x39
	for i := -5; i < 5; i++ {
		for j := 0; j <= 4; j++ {
			if i == -5 || i == 0 && j == 0 {
				lr = 0xa4
				lg = 0x81
				lb = 0x61
			} else {
				lr = 0x3b
				lg = 0x3a
				lb = 0x39
			}
			if i > -4 && i < 3 && j == 4 || i > -3 && i < 4 && j == 3 {
				lr = 0xff - rand.Intn(10)
				lg = 0xd6 - rand.Intn(20)
				lb = 0x2f
			}
			global.gWorld.SetPixel(
				x+i,
				y-j,
				uint32(lr&0xFF<<24|lg&0xFF<<16|lb&0xFF<<8|0xFF),
			)
		}
	}

	// for j := 5; j < 60; j++ {
	// 	for i := -j * 2 / 2; i < j*2/2; i++ {
	// 		c := global.gWorld.PixelColor(float64(x+i), float64(y-j))
	// 		r := c>>24&0xFF + int32(60-j)
	// 		g := c>>16&0xFF + int32(60-j)
	// 		b := c >> 8 & 0xFF
	// 		a := c & 0xFF

	// 		//	if global.gWorld.IsShadow(float64(x+i), float64(y-j)) || global.gWorld.IsBackground(float64(x+i), float64(y-j)) {
	// 		global.gWorld.SetPixel(
	// 			x+i,
	// 			y-j,
	// 			uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|a),
	// 		)
	// 		//	}
	// 	}
	// }
}

func (p *pcg) GenerateLine(x, y int) {
	r := 0x73
	g := 0x80
	b := 0x62
	a := wBackground8
	if global.gWorld.IsShadow(float64(x), float64(y)) {
		a = wShadow8
		r = 0x73 / 2
		g = 0x80 / 2
		b = 0x62 / 2
	}
	lineSize := 10

	for i := 0; i < lineSize; i++ {
		global.gWorld.SetPixel(
			x,
			y-i,
			uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|a),
		)
	}
	global.gWorld.SetPixel(
		x,
		y,
		uint32(r/2&0xFF<<24|g/2&0xFF<<16|b/2&0xFF<<8|a),
	)
	global.gWorld.SetPixel(
		x,
		y-lineSize,
		uint32(r/2&0xFF<<24|g/2&0xFF<<16|b/2&0xFF<<8|a),
	)
}

func (p *pcg) GenerateBottomAirIntake(x, y int) {
	p.airCnt++
	if p.airCnt < 250 {
		return
	}
	p.airCnt = 0

	r := 0x36
	g := 0x36
	b := 0x34
	a := wBackground8

	for i := 0; i < 6; i++ {
		for j := 1; j < 4; j++ {
			if global.gWorld.IsShadow(float64(x+i), float64(y+j)) {
				a = wShadow8
				r = 0x36
				g = 0x36
				b = 0x34
			}
			if i%2 == 0 {
				r = 0x36 / 2
				g = 0x36 / 2
				b = 0x34 / 2
			}
			if i == 0 || j == 3 {
				r = 0x0 + i
				g = 0x0 + i
				b = 0x0 + i
			}

			global.gWorld.SetPixel(
				x+i,
				y+j,
				uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|a),
			)
		}
	}
}

func (p *pcg) GenerateBottomLine(x, y int) {
	r := 0x36
	g := 0x36
	b := 0x34
	a := wBackground8
	if global.gWorld.IsShadow(float64(x), float64(y)) {
		a = wShadow8
		r = 0x36 / 2
		g = 0x36 / 2
		b = 0x34 / 2
	}
	lineSize := 3

	for i := 0; i < lineSize; i++ {
		global.gWorld.SetPixel(
			x,
			y-i,
			uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|a),
		)
	}
	// global.gWorld.SetPixel(
	// 	x,
	// 	y,
	// 	uint32(r/2&0xFF<<24|g/2&0xFF<<16|b/2&0xFF<<8|a),
	// )

}

func (p *pcg) GenerateBricks(x, y int) {
	a := wBackground8
	length := 10 + rand.Intn(30)
	height := 10 + rand.Intn(30)
	for i := 0; i < height; i++ {
		from := -length/3 - rand.Intn(length/2)
		to := length/3 + rand.Intn(length/2)
		for j := from; j < to; j++ {
			c := global.gWorld.PixelColor(float64(x+i), float64(y-j))
			r := int(c>>24&0xFF + int32(80-j))
			g := int(c>>16&0xFF + int32(80-j))
			b := int(c >> 8 & 0xFF)
			if j%10 == 0 || i%5 == 0 {
				rnd := rand.Intn(20)
				r = 0x00 + rnd
				g = 0x00 + rnd
				b = 0x00 + rnd
			} else {
				rnd := 10
				r -= rnd
				g -= rnd
				b -= rnd
			}
			if j < from+5 || j == 0 {
				r /= 2
				g /= 2
				b /= 2
			}
			if global.gWorld.IsBackground(float64(x+i), float64(y+j)) {
				global.gWorld.SetPixel(
					x+j,
					y+i,
					uint32(r/2&0xFF<<24|g/2&0xFF<<16|b/2&0xFF<<8|a),
				)
			}
		}
	}

}

func (p *pcg) GrassFloor(x, y int) {
	// Grass floor
	p.floorHeight = 5 + rand.Intn(4)
	r := 0x33
	g := 0xFF
	b := 0x33
	a := 0xFF
	for i := 0; i < p.floorHeight; i++ {
		g_ := g
		if i == p.floorHeight-1 {
			g_ -= 100
		} else if i == p.floorHeight-2 {
			g_ -= 50
		}
		global.gWorld.SetPixel(
			x,
			y-i,
			uint32(r&0xFF<<24|g_&0xFF<<16|b&0xFF<<8|a&0xFF),
		)
	}
}
