//=============================================================
// ui.go
//-------------------------------------------------------------
// User Interface (HUD) for the game
//=============================================================
package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"strconv"
	"unicode"
)

//=============================================================
//
//=============================================================
type UI struct {
	regular *text.Text
	fps     *text.Text
	canvas  *pixelgl.Canvas
}

//=============================================================
//
//=============================================================
func (u *UI) create() {
	u.canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(wViewMax), float64(wViewMax)))
	u.canvas.Clear(pixel.RGBA{0, 0, 0, 0})

	ttf, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	fface := truetype.NewFace(ttf, &truetype.Options{
		Size:              8,
		GlyphCacheEntries: 1,
	})

	regular := text.NewAtlas(
		fface,
		text.ASCII, text.RangeTable(unicode.Latin),
	)
	u.regular = text.New(pixel.ZV, regular)
	u.regular.Color = pixel.RGBA{1, 0, 1, 1}

	u.fps = text.New(pixel.ZV, regular)
	u.fps.Color = pixel.RGBA{1, 0, 1, 1}
}

//=============================================================
//
//=============================================================
func (u *UI) updateFPS(fps int) {
	u.fps.Clear()
	u.fps.WriteString(fmt.Sprintf("FPS: %v", strconv.Itoa(fps)))
}

//=============================================================
//
//=============================================================
func (u *UI) draw(dt float64) {
	u.canvas.Clear(pixel.RGBA{0, 0, 0, 0})
	u.fps.Draw(u.canvas, pixel.IM.Moved(pixel.V(1, wViewMax/2+22)))
	u.canvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(global.gCamera.pos.X+wViewMax/2.0, global.gCamera.pos.Y+wViewMax/2.0)))
}

//=============================================================
//
//=============================================================
func (u *UI) showText(text string, pos pixel.Vec, time float64) {
	for _, x := range text {
		u.regular.WriteRune(rune(x))
	}
}
