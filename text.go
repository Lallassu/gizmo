package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"unicode"
)

type font struct {
	chars  map[string]*pixel.Sprite
	width  float64
	height float64
	cWidth float64
}

func (f *font) create() {
	img, _, _, _ := loadTexture(fmt.Sprintf("%v%v", wAssetMixedPath, "font1.png"))
	pic := pixel.PictureDataFromImage(img)
	f.width = pic.Bounds().Max.X
	f.height = pic.Bounds().Max.Y

	f.cWidth = math.Floor(f.width / 33)

	f.chars = make(map[string]*pixel.Sprite, 33)
	chars := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "v", "y", "x", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	for i, c := range chars {
		f.chars[c] = pixel.NewSprite(pic, pixel.R(float64(i)*f.cWidth, 0, float64(i+1)*f.cWidth, f.height))
	}
}

func (f *font) write(text string) *pixelgl.Canvas {
	canvas := pixelgl.NewCanvas(pixel.R(0, 0, float64(len(text))*f.cWidth, f.height))
	for i, c := range text {
		if f.chars[string(unicode.ToLower(c))] != nil {
			f.chars[string(unicode.ToLower(c))].Draw(canvas, pixel.IM.Moved(pixel.V(float64(i)*f.cWidth+f.cWidth/2, f.height/2)))
		}
	}
	return canvas
}
