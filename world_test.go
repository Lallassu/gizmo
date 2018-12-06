package main

import (
	_ "github.com/faiface/pixel"
	_ "github.com/faiface/pixel/pixelgl"
	"math/rand"
	"os"
	"testing"
)

var w world

func Prepare() world {
	w := world{}
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

func BenchmarkWorldTestBit1(b *testing.B) {
	p := 0xFF44F4AF
	for i := 0; i < b.N; i++ {
		if p&0x00000080 == 0 {

		}
	}
}
func BenchmarkWorldTestBit2(b *testing.B) {
	p := 0xFF44F4AF
	for i := 0; i < b.N; i++ {
		if p&0xFF>>7 == 0 {

		}
	}
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
