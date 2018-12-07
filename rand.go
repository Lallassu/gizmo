//=============================================================
// rand.go
//-------------------------------------------------------------
// Fast Rand
// On-start generated random values used to speed up performance
// of rand() instead of using math.rand().
// Only used when performance is critical.
//=============================================================
package main

import (
	"math/rand"
	"time"
)

type fRand struct {
	randInt     []int
	randFloat64 []float64
	rintc       int
	rfloatc     int
	max         int
}

func (r *fRand) create(max int) {
	r.rintc = -1
	r.rfloatc = -1
	r.max = max

	r.randInt = make([]int, max)
	r.randFloat64 = make([]float64, max)

	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < max; i++ {
		r.randInt[i] = rand.Intn(10)
		r.randFloat64[i] = rand.Float64()
	}
}

func (r *fRand) rand() int {
	if r.rintc >= r.max-1 {
		r.rintc = -1
	}
	r.rintc++
	return r.randInt[r.rintc]
}

func (r *fRand) randFloat() float64 {
	if r.rfloatc >= r.max-1 {
		r.rfloatc = -1
	}
	r.rfloatc++
	return r.randFloat64[r.rfloatc]
}
