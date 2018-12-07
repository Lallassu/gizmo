package main

import (
	_ "github.com/faiface/pixel"
	_ "github.com/faiface/pixel/pixelgl"
	"math/rand"
	"os"
	"testing"
)

var w world
var rlist []int

func Prepare() world {
	w := world{}
	w.Init()
	w.pixels = make([]uint32, 500*500)
	w.size = 500 * 500
	for i := 0; i < 500*500; i++ {
		w.pixels[i] = uint32(rand.Intn(0xFFFFFFFF))
	}
	rlist = make([]int, 100000)
	for i := 0; i < 100000; i++ {
		rlist = append(rlist, rand.Intn(100))
	}
	return w
}

func TestMain(m *testing.M) {
	w = Prepare()
	os.Exit(m.Run())
}

func BenchmarkWorldTestRand1(b *testing.B) {
	cnt := 0
	for i := 0; i < b.N; i++ {
		getRand(&cnt)
	}
}
func getRand(cnt *int) int {
	(*cnt)++
	if *cnt >= len(rlist)-1 {
		(*cnt) = 0
	}
	(*cnt)++

	return rlist[*cnt]
}

func BenchmarkWorldTestRand2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getRand2()
	}
}
func getRand2() int {
	return rand.Intn(100)
}

func BenchmarkWorldIsBackground(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w.IsBackground(300.0, 200.0)
	}
}

func BenchmarkWorldIsShadow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w.IsShadow(300.0, 200.0)
	}
}

func BenchmarkWorldIsRegular(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w.IsRegular(300.0, 200.0)
	}
}

func BenchmarkWorldIsWall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w.IsWall(300.0, 200.0)
	}
}

func BenchmarkWorldIsLadder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w.IsLadder(300.0, 200.0)
	}
}

func BenchmarkWorldRemoveShadow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w.removeShadow(300.0, 200.0)
	}
}

func BenchmarkWorldAddShadow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w.addShadow(300.0, 200.0)
	}
}
