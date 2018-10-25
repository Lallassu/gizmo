//=============================================================
// pcg.go
//-------------------------------------------------------------
// Procedurally Generated Stuff
//=============================================================
package main

import (
	"math/rand"
)

type pcg struct {
	floorHeight int
	incr        float64
}

func (p *pcg) Flower(x, y int) {
	//color := global.gWorld.coloring.getFlower()
	stem := rand.Intn(10)
	for i := 0; i < stem; i++ {
		global.gWorld.AddPixel(x, y+i, 0xFF00FFFF)
		global.gWorld.addShadow(x+1, y+i-1)
	}
}

func (p *pcg) MetalFloor(x, y int) {
	r := 0
	g := 0
	b := 0
	for i := 0; i < 10; i++ {
		switch i {
		case 0:
			r = 0x79
			g = 0xa5
			b = 0x91
		case 1, 2:
			r = 0x57
			g = 0x7d
			b = 0x6f
		case 3:
			r = 0x14
			g = 0x27
			b = 0x23
		case 4, 6:
			r = 0x44
			g = 0x55
			b = 0x49
		case 5:
			r = 0x3a
			g = 0x4e
			b = 0x42
		case 8:
			r = 0x22
			g = 0x2c
			b = 0x27
		case 9:
			r = 0x2b
			g = 0x37
			b = 0x31

		}
		global.gWorld.AddPixel(
			x,
			y-i,
			uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
		)
	}
}

func (p *pcg) MetalRoof(x, y int) {
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
		global.gWorld.AddPixel(
			x,
			y-i,
			uint32(r&0xFF<<24|g&0xFF<<16|b&0xFF<<8|0xFF),
		)
		for n := 0; n < wShadowLength; n++ {
			global.gWorld.addShadow(x+n, y-i-n)
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
		global.gWorld.AddPixel(
			x,
			y-i,
			uint32(r&0xFF<<24|g_&0xFF<<16|b&0xFF<<8|a&0xFF),
		)
	}
}
