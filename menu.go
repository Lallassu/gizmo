//=============================================================
// menu.go
//-------------------------------------------------------------
// Menu for the game
//=============================================================
package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

//=============================================================
//
//=============================================================
type menu struct {
	items   []*menuItem
	visible bool
	uTime   float32
}

//=============================================================
//
//=============================================================
type menuItem struct {
	action   func()
	name     string
	canvas   *pixelgl.Canvas
	selected int32
}

//=============================================================
//
//=============================================================
func (m *menu) create() {
	m.items = make([]*menuItem, 0)
	m.addItem("New Game", func() { Debug("New GAME") })
	m.addItem("Continue", func() { Debug("Continue") })
	m.addItem("Quit", func() { Debug("Quit") })
	m.items[0].selected = 1
	m.visible = true
}

//=============================================================
//
//=============================================================
func (m *menu) addItem(str string, f func()) {
	canvas := pixelgl.NewCanvas(pixel.R(0, 0, 100, 100))
	canvas.SetUniform("uTime", &m.uTime)
	canvas.SetFragmentShader(fragmentShaderMenuItem)
	global.gFont.writeToCanvas(str, canvas)
	item := &menuItem{
		canvas:   canvas,
		name:     str,
		action:   f,
		selected: 0,
	}

	//item.canvas.SetUniform("uSelected", &item.selected)

	m.items = append(m.items, item)

}

//=============================================================
//
//=============================================================
func (m *menu) selectItem() {
	for _, item := range m.items {
		if item.selected == 1 {
			item.action()
			return
		}
	}
}

//=============================================================
//
//=============================================================
func (m *menu) moveUp() {
	for i, item := range m.items {
		if item.selected == 1 {
			if i > 0 {
				m.items[i-1].selected = 1
			} else {
				m.items[len(m.items)-1].selected = 1
			}
			m.items[i].selected = 0
			break
		}
	}
}

//=============================================================
//
//=============================================================
func (m *menu) moveDown() {
	for i, item := range m.items {
		if item.selected == 1 {
			if i < len(m.items)-1 {
				m.items[i+1].selected = 1
			} else {
				m.items[0].selected = 1
			}
			m.items[i].selected = 0
			break
		}
	}
}

//=============================================================
//
//=============================================================
func (m *menu) draw(dt float64) {
	if m.visible {
		m.uTime += float32(dt)
		offset_y := 40.0
		offset_x := 30.0
		// Draw in the middle of the screen.
		for i, item := range m.items {
			item.canvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(global.gCamera.pos.X+wViewMax/2.0-offset_x, global.gCamera.pos.Y+wViewMax/2-float64(i)*offset_y)))
		}
	}
}
