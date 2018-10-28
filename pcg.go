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
}

func (p *pcg) Flower(x, y int) {
	//color := global.gWorld.coloring.getFlower()
	stem := rand.Intn(10)
	for i := 0; i < stem; i++ {
		global.gWorld.AddPixel(x, y+i, 0xFF00FFFF)
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
				global.gWorld.AddPixel(
					x+i,
					y+j,
					uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
				)
			} else {
				global.gWorld.AddPixel(
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
				global.gWorld.AddPixel(
					x-i,
					y-j,
					uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
				)
			} else {
				global.gWorld.AddPixel(
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
			if i == 5 || j == 5 || (j == 1 && i == 1) || (j == 1 && i == 0) || (j == 0 && i == 1) {
				r = 0x1e
				g = 0x2c
				b = 0x26
			} else if i == 0 && j == 0 {
				r = 0x2b
				g = 0x37
				b = 0x31
			} else if i == 4 || j == 4 || (j == 2 && i == 1) || (i == 2 && j == 1) || (j == 2 && i == 2) || (j == 2 && i == 0) || (j == 0 && i == 2) {
				r = 0x44
				g = 0x55
				b = 0x49
			} else if i == 3 || j == 3 {
				r = 0x3a
				g = 0x4e
				b = 0x42
			}
			if left {
				global.gWorld.AddPixel(
					x+i,
					y+j,
					uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
				)
			} else {
				global.gWorld.AddPixel(
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
			if i == 5 || j == 5 || (j == 1 && i == 1) || (j == 1 && i == 0) || (j == 0 && i == 1) {
				r = 0x1e
				g = 0x2c
				b = 0x26
			} else if i == 0 && j == 0 {
				r = 0x2b
				g = 0x37
				b = 0x31
			} else if i == 4 || j == 4 || (j == 2 && i == 1) || (i == 2 && j == 1) || (j == 2 && i == 2) || (j == 2 && i == 0) || (j == 0 && i == 2) {
				r = 0x44
				g = 0x55
				b = 0x49
			} else if i == 3 || j == 3 {
				r = 0x3a
				g = 0x4e
				b = 0x42
			}
			if left {
				global.gWorld.AddPixel(
					x+i,
					y-j,
					uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
				)
			} else {
				global.gWorld.AddPixel(
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
		global.gWorld.AddPixel(
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
			global.gWorld.AddPixel(
				x,
				y+i,
				uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
			)
		} else {
			global.gWorld.AddPixel(
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

		global.gWorld.AddPixel(
			x,
			y-i,
			uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
		)
	}
	if p.floorCnt == 50 {
		p.floorCnt = 0
	}
}

func (p *pcg) GenerateDoor(x, y int) {
	r := 0x73
	g := 0x80
	b := 0x62
	a := wBackground8
	doorLen := 30
	doorHeight := 40

	for i := 0; i < doorLen; i++ {
		for j := 0; j < doorHeight; j++ {
			// Frame
			if i == 0 || j == doorHeight-1 || i == doorLen-1 {
				r = 0x3a
				g = 0x3a
				b = 0x3a
			} else {
				r = 0x51
				g = 0x71
				b = 0x74
			}
			if (i == doorLen-2 && j < doorHeight-1) || (i == doorLen-3 && j < doorHeight-1) {
				r = 0x4c
				g = 0x4c
				b = 0x4c
			}
			if j == doorHeight-2 && i > 0 && i < doorLen-2 {
				r = 0x1b
				g = 0x1b
				b = 0x1b
			}
			// Shadow
			if j > doorHeight-5 || i < 3 {
				r /= 3
				g /= 3
				b /= 3
			}
			// Handle
			if j == doorHeight/3 && i > 1 && i < 5 {
				r = 0x52
				g = 0x67
				b = 0x69
			}
			// Handle
			if j == doorHeight/3-1 && i > 1 && i < 4 {
				r = 0x35
				g = 0x44
				b = 0x46
			}
			// Hinge
			if (j == doorHeight/4 || j == doorHeight/4-1 || j == doorHeight-(doorHeight/4-1) || j == doorHeight-(doorHeight/4)) && i == doorLen-3 {
				r = 0x4c
				g = 0x59
				b = 0x5c
			}
			// Hinge
			if (j == doorHeight/4 || j == doorHeight/4-1 || j == doorHeight-(doorHeight/4-1) || j == doorHeight-(doorHeight/4)) && i == doorLen-4 {
				r = 0x23
				g = 0x23
				b = 0x23
			}
			global.gWorld.AddPixel(
				x+i,
				y+j,
				uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|a),
			)
		}
	}
}

func (p *pcg) GenerateLine(x, y int) {
	r := 0x73
	g := 0x80
	b := 0x62
	a := wBackground8
	lineSize := 10

	for i := 0; i < lineSize; i++ {
		global.gWorld.AddPixel(
			x,
			y-i,
			uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|a),
		)
	}
	global.gWorld.AddPixel(
		x,
		y,
		uint32(r/2&0xFF<<24|g/2&0xFF<<16|b/2&0xFF<<8|a),
	)
	global.gWorld.AddPixel(
		x,
		y-lineSize,
		uint32(r/2&0xFF<<24|g/2&0xFF<<16|b/2&0xFF<<8|a),
	)
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
		global.gWorld.AddPixel(
			x,
			y-i,
			uint32(r&0xFF<<24|g_&0xFF<<16|b&0xFF<<8|a&0xFF),
		)
	}
}
