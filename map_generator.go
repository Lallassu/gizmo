//=============================================================
// map_generator.go
//-------------------------------------------------------------
//=============================================================
package main

import (
	"math/rand"
	"time"
)

type generator struct {
	cell      [][]uint32
	minx      int
	maxx      int
	minz      int
	maxz      int
	wide      int
	deep      int
	cell_size int
}

func (g *generator) NewWorld(size_x, size_z int) []uint32 {
	rand.Seed(time.Now().Unix())
	g.cell_size = 20
	g.wide = size_x / g.cell_size
	g.deep = size_z / g.cell_size
	g.minx, g.minz = 1, 1
	g.maxx = g.wide - 2
	g.maxz = g.deep - 2

	g.cell = make([][]uint32, g.wide)

	for x := 0; x < g.wide; x++ {
		g.cell[x] = make([]uint32, g.deep)
	}

	num_steps := 100
	step_length := 10

	g.randomWalk(num_steps, step_length)
	g.cleanDeadEnds()
	g.cleanDeadEnds()

	Debug("Generating...")
	pixels := make([]uint32, 0)
	for x := 0; x < g.wide; x++ {
		for z := 0; z < g.deep; z++ {
			if g.cell[x][z] == 1 {
				for i := 0; i < g.cell_size; i++ {
					for j := 0; j < g.cell_size; j++ {
						//world_.AddPixel(x*g.cell_size+i, z*g.cell_size+j, uint32(0xFF0000FF))
						pixels = append(pixels, uint32(x*g.cell_size+i))
						pixels = append(pixels, uint32(z*g.cell_size+j))
					}
				}
			}
		}
	}
	Debug("Map generation complete.")
	return pixels
}

func (g *generator) randomWalk(num_steps, step_length int) {
	px := g.wide / 2
	pz := g.deep / 2

	for i := 0; i < num_steps; i++ {
		g.cell[px][pz] = 1

		if rand.Float64() < 0.5 {
			g.makeRoom(px, pz, 1, step_length/2+1)
		}

		dx := []int{0, 1, 0, -1}
		dz := []int{1, 0, -1, 0}

		d := rand.Intn(4)

		nx := px + dx[d]*step_length
		nz := pz + dz[d]*step_length

		for !g.validCell(nx, nz) {
			d = rand.Intn(4)

			nx = px + dx[d]*step_length
			nz = pz + dz[d]*step_length
		}

		for x := px; x != nx; x += dx[d] {
			g.cell[x][pz] = 1
		}

		for z := pz; z != nz; z += dz[d] {
			g.cell[px][z] = 1
		}

		px = nx
		pz = nz
	}
}

func (g *generator) random(min, max int) int {
	return rand.Intn(max-min) + min
}

func (g *generator) validCell(x, z int) bool {
	if x < g.minx || x > g.maxx {
		return false
	}
	if z < g.minz || z > g.maxz {
		return false
	}
	return true
}

func (g *generator) makeRoom(x, z, minw, maxw int) {
	room_width := g.random(minw, maxw)
	room_depth := g.random(minw, maxw)

	for cx := x - room_width; cx <= x+room_width; cx++ {
		for cz := z; cz <= z+room_depth; cz++ {
			if g.validCell(cx, cz) {
				g.cell[cx][cz] = 1
			}
		}
	}
}

func (g *generator) cleanDeadEnds() {
	dx := []int{0, 1, 0, -1}
	dz := []int{1, 0, -1, 0}

	for x := 0; x < g.wide; x++ {
		for z := 0; z < g.deep; z++ {
			if g.cell[x][z] == 1 {
				sum := 0

				for i := 0; i < 4; i++ {
					nx := x + dx[i]
					nz := z + dz[i]

					if g.validCell(nx, nz) {
						if g.cell[nx][nz] == 1 {
							sum++
						}
					}
				}

				if sum <= 1 {
					g.cell[x][z] = 0
				}

			}
		}
	}

	for x := g.wide - 1; x >= 0; x-- {
		for z := g.deep - 1; z >= 0; z-- {
			if g.cell[x][z] == 1 {
				sum := 0

				for i := 0; i < 4; i++ {
					nx := x + dx[i]
					nz := z + dz[i]

					if g.validCell(nx, nz) {
						if g.cell[nx][nz] == 1 {
							sum++
						}
					}
				}

				if sum <= 1 {
					g.cell[x][z] = 0
				}
			}
		}
	}

}
