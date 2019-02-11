package main

import (
	"encoding/json"
	"image"
	"os"
	"regexp"

	"github.com/faiface/pixel"
)

// Textures loaded from packed texture atlas and JSON descr.
type textures struct {
	batch   *pixel.Batch
	image   *pixel.PictureData
	sprites map[string]*pixel.Sprite
	objects []*sprite
}

// Sprite contains an object where position is volatile
type sprite struct {
	name  string
	pos   pixel.Vec
	scale float64
}

// JSON parsing for texture packed information
type frames struct {
	Frames []files `json:"frames"`
	Meta   meta    `json:"meta"`
}

type meta struct {
	App         string `json:"app"`
	Version     string `json:"version"`
	Image       string `json:"image"`
	Format      string `json:"format"`
	Size        frame  `json:"size"`
	Scale       string `json:"scale"`
	Smartupdate string `json:"smartupdate"`
}

type files struct {
	Filename         string `json:"filename"`
	Frame            frame  `json:"frame"`
	Rotated          bool   `json:"rotated"`
	Trimmed          bool   `json:"trimmed"`
	SpriteSourceSize frame  `json:"spriteSourceSize"`
	SourceSize       frame  `json:"sourceSize"`
}

type frame struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	W float64 `json:"w"`
	H float64 `json:"h"`
}

// Load and create textures
func (t *textures) load(jsonfile string) {
	t.sprites = make(map[string]*pixel.Sprite)

	// Read json file to names (tr .png from name)
	jfile, err := os.Open(jsonfile)
	defer jfile.Close()
	if err != nil {
		Error("Read JSON file failed:", err.Error())
	}

	jsonParser := json.NewDecoder(jfile)
	result := frames{}
	err = jsonParser.Decode(&result)
	if err != nil {
		Error("Error decoding JSON:", err)
	}

	// Load the image
	file, err := os.Open(result.Meta.Image)
	if err != nil {
		Error("Error loading texture file:", err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		Error("Error decoding texture file:", err)
	}
	t.image = pixel.PictureDataFromImage(img)

	// Create a batch from image
	t.batch = pixel.NewBatch(&pixel.TrianglesData{}, t.image)

	// Create sprite for each frame t.sprites.
	reg := regexp.MustCompile(`\..*`)
	for _, f := range result.Frames {
		name := reg.ReplaceAllString(f.Filename, "${1}")
		Debug("F:", f.Filename, pixel.R(f.Frame.X, f.Frame.Y, f.Frame.X+f.Frame.W, f.Frame.Y+f.Frame.H))
		t.sprites[name] = pixel.NewSprite(t.image, pixel.R(f.Frame.X, f.Frame.Y, f.Frame.X+f.Frame.W, f.Frame.Y+f.Frame.H))
	}
}

// Add new sprite object to draw
func (t *textures) addObject(o *sprite) {
	t.objects = append(t.objects, o)
}

// Remove object from drawing list
func (t *textures) removeObject(o sprite) {
	// TBD
}

// Get info for sprite
func (t *textures) spriteInfo(name string) (int, int) {
	f := t.sprites[name].Frame()
	return int(f.Max.X - f.Min.X), int(f.Max.Y - f.Min.Y)
}

// Draw the batch based on the packed texture atlas
func (t *textures) update(dt float64) {
	t.batch.Clear()
	for _, o := range t.objects {
		t.sprites[o.name].Draw(t.batch, pixel.IM.Scaled(pixel.ZV, o.scale).Moved(o.pos))
	}
	t.batch.Draw(global.gWin)
}
