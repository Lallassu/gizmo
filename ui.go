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
	fps           *text.Text
	middleText    *text.Text
	canvas        *pixelgl.Canvas
	miniMapScale  float64
	middleTextStr string
}

//=============================================================
//
//=============================================================
func (u *UI) create() {
	u.canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(wViewMax), float64(wViewMax)))
	u.canvas.Clear(pixel.RGBA{0, 0, 0, 0})
	u.miniMapScale = 0.08

	ttf, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	fface := truetype.NewFace(ttf, &truetype.Options{
		Size:              8,
		GlyphCacheEntries: 1,
	})

	ffaceMiddle := truetype.NewFace(ttf, &truetype.Options{
		Size:              30,
		GlyphCacheEntries: 1,
	})

	regularMiddle := text.NewAtlas(
		ffaceMiddle,
		text.ASCII, text.RangeTable(unicode.Latin),
	)

	regular := text.NewAtlas(
		fface,
		text.ASCII, text.RangeTable(unicode.Latin),
	)
	u.fps = text.New(pixel.ZV, regular)
	u.fps.Color = pixel.RGBA{1, 0, 1, 1}

	u.middleText = text.New(pixel.ZV, regularMiddle)
	u.middleText.Color = pixel.RGBA{1, 0, 0, 1}
}

//=============================================================
// Mini map
//=============================================================
func (u *UI) updateMiniMap() {
	pos := global.gPlayer.getPosition()
	canvas := pixelgl.NewCanvas(pixel.R(0, 0, 1, 1))
	canvas.Clear(pixel.RGBA{1.0, 0, 0, 0.5})
	offset_x := float64(global.gWorld.width/2) * u.miniMapScale
	offset_y := float64(global.gWorld.height/2) * u.miniMapScale
	global.gWorld.bgSprite.Draw(u.canvas, pixel.IM.ScaledXY(pixel.V(u.miniMapScale, u.miniMapScale), pixel.V(u.miniMapScale, u.miniMapScale)).Moved(pixel.V(offset_x, offset_y)))
	canvas.Draw(u.canvas, pixel.IM.Moved(pixel.V(u.miniMapScale*pos.X+offset_x-float64(global.gWorld.width/2)*u.miniMapScale, u.miniMapScale*pos.Y+offset_y-float64(global.gWorld.height/2)*u.miniMapScale)))
}

//=============================================================
//
//=============================================================
func (u *UI) setMiddleText(text string) {
	u.middleTextStr = text
	u.middleText.Clear()
	u.middleText.WriteString(fmt.Sprintf("%v", text))
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
	u.updateMiniMap()
	u.fps.Draw(u.canvas, pixel.IM.Moved(pixel.V(1, wViewMax/2+22)))
	u.canvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(global.gCamera.pos.X+wViewMax/2.0, global.gCamera.pos.Y+wViewMax/2.0)))

	u.middleText.Draw(u.canvas, pixel.IM.Moved(pixel.V(float64(wViewMax/2-len(u.middleTextStr)/2), wViewMax/2+22)))
}
