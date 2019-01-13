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
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"strconv"
	"unicode"
)

//=============================================================
//
//=============================================================
type UI struct {
	fps                *text.Text
	middleText         *text.Text
	lifeText           *text.Text
	canvas             *pixelgl.Canvas
	miniMapSprite      *pixel.Sprite
	miniMapFrameCanvas *pixelgl.Canvas
	miniMapCanvas      *pixelgl.Canvas
	miniMapScale       float64
	middleTextStr      string
	deathScreenTimer   float64
	uPos               mgl32.Vec2
	uTime              float32
	lifeCanvas         *pixelgl.Canvas
}

//=============================================================
//
//=============================================================
func (u *UI) create() {
	u.lifeCanvas = global.gFont.write("100")

	u.canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(wViewMax), float64(wViewMax)))
	u.canvas.Clear(pixel.RGBA{0, 0, 0, 0})

	u.miniMapCanvas = pixelgl.NewCanvas(pixel.R(0, 0, 100, 100))
	u.miniMapCanvas.SetUniform("uPos", &u.uPos)
	u.miniMapCanvas.SetUniform("uTime", &u.uTime)
	u.miniMapCanvas.SetFragmentShader(fragmentShaderMinimap)

	img, _, _, _ := loadTexture(fmt.Sprintf("%v%v", wAssetObjectsPath, "minimap.png"))
	pic := pixel.PictureDataFromImage(img)
	u.miniMapSprite = pixel.NewSprite(pic, pic.Bounds())
	u.miniMapFrameCanvas = pixelgl.NewCanvas(pic.Bounds())

	ttf, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	fface := truetype.NewFace(ttf, &truetype.Options{
		Size:              wFPSTextSize,
		GlyphCacheEntries: 1,
	})

	ffaceMiddle := truetype.NewFace(ttf, &truetype.Options{
		Size:              wMiddleTextSize,
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
	u.middleText.Color = pixel.RGBA{1, 1, 1, 1}

	u.lifeText = text.New(pixel.ZV, regular)
	u.lifeText.Color = pixel.RGBA{1, 0, 0, 1}
}

//=============================================================
// Mini map
//=============================================================
func (u *UI) updateMiniMap() {
	u.miniMapScale = 0.25 / (float64(global.gWorld.width) / 1024)
	pos := global.gPlayer.getPosition()
	//canvas := pixelgl.NewCanvas(pixel.R(0, 0, 1, 1))
	//canvas.Clear(pixel.RGBA{1.0, 0, 0, 0.5})

	offset_x := float64(global.gWorld.width/2) * u.miniMapScale
	offset_y := float64(global.gWorld.height/2) * u.miniMapScale
	//offset_x2 := offset_x //- float64(global.gWorld.width/2)*u.miniMapScale
	//offset_y2 := offset_y //- float64(global.gWorld.height/2)*u.miniMapScale

	u.uPos = mgl32.Vec2{float32(offset_x / 2), float32(offset_y / 2)}

	bounds := u.miniMapFrameCanvas.Bounds()
	//u.miniMapSprite.Draw(u.miniMapFrameCanvas, pixel.IM.Moved(pixel.V(offset_x2, offset_y2)).ScaledXY(pixel.V(0.8, 0.8), pixel.V(0.8, 0.8)))
	u.miniMapSprite.Draw(u.miniMapFrameCanvas, pixel.IM.Moved(pixel.V(bounds.Max.X/2, bounds.Max.Y/2)).ScaledXY(pixel.V(0.6, 0.6), pixel.V(0.6, 0.6)))

	//global.gWorld.bgSprite.Draw(u.miniMapCanvas, pixel.IM.ScaledXY(pixel.V(u.miniMapScale, u.miniMapScale), pixel.V(u.miniMapScale/2, u.miniMapScale/2)).Moved(pixel.V(u.miniMapScale*pos.X+offset_x-float64(global.gWorld.width/2)*u.miniMapScale, u.miniMapScale*pos.Y+offset_y-float64(global.gWorld.height/2)*u.miniMapScale)))
	bounds = global.gWorld.bgSprite.Frame()
	global.gWorld.bgSprite.Draw(u.miniMapCanvas,
		pixel.IM.ScaledXY(
			pixel.V(u.miniMapScale, u.miniMapScale),
			pixel.V(u.miniMapScale/2, u.miniMapScale/2),
		).Moved(
			pixel.V((bounds.Max.X*u.miniMapScale)/2+pos.X*u.miniMapScale, (bounds.Max.Y*u.miniMapScale)/2+pos.Y*u.miniMapScale)))
	//canvas.Draw(u.canvas, pixel.IM.Moved(pixel.V(u.miniMapScale*pos.X+offset_x-float64(global.gWorld.width/2)*u.miniMapScale, u.miniMapScale*pos.Y+offset_y-float64(global.gWorld.height/2)*u.miniMapScale)))
	//canvas.Draw(u.canvas, pixel.IM.Moved(pixel.V(offset_x, offset_y)))
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

	u.uTime += float32(dt)

	// Draw death screen
	color := pixel.RGBA{}
	if global.gPlayer.life == 0 {
		u.deathScreenTimer += dt
		u.setMiddleText(wDeathScreenText)

		red := u.deathScreenTimer / 10
		if red > 0.5 {
			red = 0.5
		}
		color = pixel.RGBA{red, 0, 0, u.deathScreenTimer / 10}
	} else {
		u.deathScreenTimer = 0
	}
	u.canvas.Clear(color)
	u.miniMapCanvas.Clear(color)
	u.miniMapFrameCanvas.Clear(color)

	u.updateMiniMap()

	u.fps.Draw(u.canvas, pixel.IM.Moved(pixel.V(1, wViewMax/2+22)))

	u.middleText.Draw(u.canvas, pixel.IM.Moved(pixel.V(float64(wViewMax/2-((len(u.middleTextStr)/3)*wMiddleTextSize)), wViewMax/3)))

	//u.lifeText.Clear()
	//u.lifeText.WriteString(fmt.Sprintf("Life: %v", global.gPlayer.life))
	//u.lifeText.Draw(u.canvas, pixel.IM.Scaled(pixel.ZV, 0.25).Moved(pixel.V(1, wViewMax/2+40)))

	u.lifeCanvas.Draw(u.canvas, pixel.IM.Moved(pixel.V(100, wViewMax/2+40)))

	bounds := u.miniMapFrameCanvas.Bounds()
	u.miniMapFrameCanvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(global.gCamera.pos.X+bounds.Max.X/2, global.gCamera.pos.Y+bounds.Max.Y/2)))
	bounds = u.miniMapCanvas.Bounds()
	u.miniMapCanvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(global.gCamera.pos.X, global.gCamera.pos.Y)))
	//u.miniMapCanvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(global.gCamera.pos.X+wViewMax/2.0, global.gCamera.pos.Y+wViewMax/2.0)))

	u.canvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(global.gCamera.pos.X+wViewMax/2.0, global.gCamera.pos.Y+wViewMax/2.0)))

}
