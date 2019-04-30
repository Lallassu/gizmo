package main

import (
	"github.com/faiface/pixel"
	"image"
	"math"
	"os"
)

type Map struct {
	sprite  *pixel.Sprite
	sprite2 *pixel.Sprite
}

func (m *Map) newMapFromImg() {
	var img image.Image
	var width float64
	var height float64
	// img, width, height, _ = loadTexture(fmt.Sprintf("%v%v", wAssetMapsPath, "map1_map.png"))

	file, err := os.Open("assets/maps/map1_map.png")
	if err != nil {
		Error("Error loading texture file:", err)
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	if err != nil {
		Error("Error decoding texture file:", err)
	}

	img2 := pixel.PictureDataFromImage(img)
	width = img2.Rect.Max.X
	height = img2.Rect.Max.Y

	size := width * width
	if height > width {
		size = height * height
	}
	global.gWorld.NewMap(width, height, size)
	for x := 0.0; x < width; x++ {
		for y := 0.0; y < height; y++ {
			c := img2.Pix[img2.Index(pixel.Vec{x, y})]

			if c.A == 0 {
				//global.gWorld.SetPixel(x, y, (r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | wBackground8&0xFF))
				continue
			}
			global.gWorld.SetPixel(int(x), int(y), (uint32(c.R)&0xFF<<24 | uint32(c.G)&0xFF<<16 | uint32(c.B)&0xFF<<8 | 0xFF))
		}
	}
	// Add "red" for each world piece.
	// ww := int(width)
	// hh := int(height)
	// tmp := 0
	// for x := 0; x <= ww; x++ {
	// 	for y := 0; y <= hh; y++ {
	// 		r, g, b, a := img.At(x, hh-y).RGBA()
	// 		if a == 0 {
	// 			//global.gWorld.SetPixel(x, y, (r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | wBackground8&0xFF))
	// 			continue
	// 		}
	// 		tmp++
	// 		global.gWorld.SetPixel(x, y, (r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | 0xFF))
	// 	}
	// }

	global.gWorld.buildAllChunks()
	global.gPlayer.create(float64(100), float64(100))

	for i := 0.0; i < 10.0; i++ {
		w := &weapon{}
		w.newWeapon(100+(i*20), 300.0+i, weaponAk47)
	}

	// Load the image bg image
	file, err = os.Open("assets/maps/map1_bg.png")
	if err != nil {
		Error("Error loading texture file:", err)
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	if err != nil {
		Error("Error decoding texture file:", err)
	}
	image2 := pixel.PictureDataFromImage(img)
	m.sprite = pixel.NewSprite(image2, image2.Bounds())

	file, err = os.Open("assets/maps/map1_fg.png")
	if err != nil {
		Error("Error loading texture file:", err)
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	if err != nil {
		Error("Error decoding texture file:", err)
	}
	image3 := pixel.PictureDataFromImage(img)

	m.sprite2 = pixel.NewSprite(image3, image3.Bounds())

}

func (m *Map) updateBG(dt float32) {
	v := global.gWin.Bounds().Center()
	v.X += global.gPlayer.getPosition().X
	v.Y = 200
	m.sprite.Draw(global.gWin, pixel.IM.Moved(v))
}
func (m *Map) updateFG(dt float32) {
	v := global.gWin.Bounds().Center()
	v.Y += -150 + math.Sin(v.X)
	m.sprite2.Draw(global.gWin, pixel.IM.Moved(v))
}
