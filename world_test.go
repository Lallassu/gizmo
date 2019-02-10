package main

import (
	"math/rand"
	"os"
	"testing"

	_ "github.com/faiface/pixel"
	_ "github.com/faiface/pixel/pixelgl"
)

var w world

func Prepare() world {
	w := world{}
	global.gRand.create(100000)
	w.Init()
	w.pixels = make([]uint32, 500*500)
	w.size = 500 * 500
	for i := 0; i < 500*500; i++ {
		w.pixels[i] = uint32(rand.Intn(0xFFFFFFFF))
	}
	return w
}

func TestMain(m *testing.M) {
	w = Prepare()
	os.Exit(m.Run())
}

func BenchmarkOwnRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		global.gRand.rand()
	}
}

func BenchmarkMathRand2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getRand2()
	}
}
func getRand2() int {
	return rand.Intn(10)
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
