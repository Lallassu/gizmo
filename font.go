package main

import (
	"fmt"
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// font is a font mapping system
type font struct {
	chars   map[string]*pixel.Sprite
	cHeight float64
	cWidth  float64
}

// create  Inits font system
func (f *font) create() {
	img, _, _, _ := loadTexture(fmt.Sprintf("%v%v", wAssetMixedPath, "font.png"))
	pic := pixel.PictureDataFromImage(img)

	// For this specific image.
	f.cWidth = 18
	f.cHeight = 32

	f.chars = make(map[string]*pixel.Sprite, int(pic.Bounds().Max.X/f.cWidth*pic.Bounds().Max.Y/f.cHeight))

	// This maps directly to the font image.
	chars := []string{
		" ", "!", "\"", "#", "$", "%", "", "'", "(", ")", "x", "+", ",", "-", ".", "/", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ":", ";", "<", "=", ">", "?",
		"@", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "[", "\\", "]", "^", "_",
		"'", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "{", "|", "}", "~", " ",
	}

	y := (pic.Bounds().Max.Y / f.cHeight) - 2
	x := 0
	for i, c := range chars {
		if i%(int(pic.Bounds().Max.X)/int(f.cWidth)) == 0 && i > 0 {
			y--
			x = 0
		}
		f.chars[c] = pixel.NewSprite(pic, pixel.R(float64(x)*f.cWidth, y*f.cHeight, float64(x+1)*f.cWidth, y*f.cHeight+f.cHeight))
		x++
	}
}

// write text to canvas and return a canvas
func (f *font) write(text string) *pixelgl.Canvas {
	canvas := pixelgl.NewCanvas(pixel.R(0, 0, float64(len(text))*f.cWidth, f.cHeight))
	for i, c := range text {
		if f.chars[string(unicode.ToLower(c))] != nil {
			f.chars[string(unicode.ToLower(c))].Draw(canvas, pixel.IM.Moved(pixel.V(float64(i)*f.cWidth+f.cWidth/2, f.cHeight/2)))
		}
	}
	return canvas
}

// Draw text to given canvas (pre-created canvas for shaders)
func (f *font) writeToCanvas(text string, canvas *pixelgl.Canvas) {
	canvas.SetBounds(pixel.R(0, 0, float64(len(text))*f.cWidth, f.cHeight))
	for i, c := range text {
		if f.chars[string(unicode.ToLower(c))] != nil {
			f.chars[string(unicode.ToLower(c))].Draw(canvas, pixel.IM.Moved(pixel.V(float64(i)*f.cWidth+f.cWidth/2, f.cHeight/2)))
		}
	}
}
