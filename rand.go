//=============================================================
// rand.go
//-------------------------------------------------------------
// Fast Rand
// On-start generated random values used to speed up performance
// of rand() instead of using math.rand().
//=============================================================
package main

import (
	"math/rand"
	"time"
)

type fRand struct {
	randInt10   []int
	randInt100  []int
	randFloat64 []float64
	r10c        int
	r100c       int
	rfloatc     int
	max         int
}

func (r *fRand) create(max int) {
	r.r10c = -1
	r.r100c = -1
	r.rfloatc = -1
	r.max = max

	r.randInt10 = make([]int, max)
	r.randInt100 = make([]int, max)
	r.randFloat64 = make([]float64, max)

	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < max; i++ {
		//r.randInt10 = append(r.randInt10, rand.Intn(10))
		r.randInt10[i] = rand.Intn(10)
		r.randInt100[i] = rand.Intn(100)
		r.randFloat64[i] = rand.Float64()
	}
}

func (r *fRand) rand10() int {
	if r.r10c >= r.max-1 {
		r.r10c = -1
	}
	r.r10c++
	return r.randInt10[r.r10c]
}

func (r *fRand) rand100() int {
	if r.r100c >= r.max {
		r.r100c = -1
	}
	r.r100c++
	return r.randInt100[r.r100c]
}

func (r *fRand) randFloat() float64 {
	if r.rfloatc >= r.max {
		r.rfloatc = -1
	}
	r.rfloatc++
	return r.randFloat64[r.rfloatc]
}
