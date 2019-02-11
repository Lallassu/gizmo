package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type menu struct {
	items []*menuItem
	logo  *menuItem
}

type menuItem struct {
	action   func()
	name     string
	canvas   *pixelgl.Canvas
	scale    float64
	selected int32
}

func (m *menu) create() {
	m.items = make([]*menuItem, 0)
}

func (m *menu) createMain() {
	m.logo = &menuItem{
		canvas:   pixelgl.NewCanvas(pixel.R(0, 0, 1, 1)),
		name:     "GoD",
		selected: 2,
		scale:    2.0,
	}
	m.logo.canvas.SetUniform("uSelected", &m.logo.selected)
	m.logo.canvas.SetUniform("uTime", &global.uTime)
	m.logo.canvas.SetFragmentShader(fragmentShaderMenuItem)

	m.create()
	m.addItem(1.0, "New Game", func() { setup() })
	m.addItem(1.0, "Continue", func() { Debug("Continue") })
	m.addItem(1.0, "Options", func() {
		global.gActiveMenu = global.gOptionsMenu
	})
	m.addItem(1.0, "About", func() { Debug("About") })
	m.addItem(1.0, "Quit", func() { global.gController.quit = true })
	m.items[0].selected = 1
}

func (m *menu) createOptions() {
	m.create()
	m.addItem(0.7, "Controls", func() {
		global.gActiveMenu = global.gControllerMenu
	})
	m.addItem(0.7, "Display", func() {
		global.gActiveMenu = global.gDisplayMenu
	})
	m.addItem(0.7, "Game", func() {
		global.gActiveMenu = global.gGameMenu
	})
	m.addItem(0.8, "Back", func() {
		global.gActiveMenu = global.gMainMenu

	})
	m.items[0].selected = 1
}

func (m *menu) createController() {
	m.create()
	m.addItem(0.5, fmt.Sprintf("%10v: %-10v", "Shoot", "KEY-X"), func() {})
	m.addItem(0.5, fmt.Sprintf("%10v: %-10v", "Jump", "KEY-X"), func() {})
	m.addItem(0.5, fmt.Sprintf("%10v: %-10v", "Duck", "KEY-X"), func() {})
	m.addItem(0.5, fmt.Sprintf("%10v: %-10v", "Left", "KEY-X"), func() {})
	m.addItem(0.5, fmt.Sprintf("%10v: %-10v", "Right", "KEY-X"), func() {})
	m.addItem(0.5, fmt.Sprintf("%10v: %-10v", "Climb", "KEY-X"), func() {})
	m.addItem(0.5, fmt.Sprintf("%10v: %-10v", "Action", "KEY-X"), func() {})
	m.addItem(0.5, fmt.Sprintf("%10v: %-10v", "Drop", "KEY-X"), func() {})
	m.addItem(0.5, fmt.Sprintf("%10v: %-10v", "Pickup", "KEY-X"), func() {})
	m.addItem(0.6, "Back", func() {
		global.gActiveMenu = global.gOptionsMenu
	})
	m.items[0].selected = 1
}

// Handle display settings
func (m *menu) createDisplay() {
	m.create()
	m.addItem(0.5, fmt.Sprintf("%20v: %-10v", "Resolution", fmt.Sprintf("%v x %v", global.gVariableConfig.WindowWidth, global.gVariableConfig.WindowHeight)),
		func() {
			// TBD: List of resolutions to toggle between
			global.gVariableConfig.SaveConfiguration()
		})
	m.addItem(0.5, fmt.Sprintf("%20v: %-10v", "V-sync", global.gVariableConfig.Vsync),
		func() {
			if global.gVariableConfig.Vsync {
				global.gVariableConfig.Vsync = false
			} else {
				global.gVariableConfig.Vsync = true
			}
			m.updateSelectedItemText(fmt.Sprintf("%20v: %-10v", "V-sync", global.gVariableConfig.Vsync))
			global.gVariableConfig.SaveConfiguration()
		})
	m.addItem(0.5, fmt.Sprintf("%20v: %-10v", "Fullscreen", global.gVariableConfig.Fullscreen),
		func() {
			if global.gVariableConfig.Fullscreen {
				global.gVariableConfig.Fullscreen = false
			} else {
				global.gVariableConfig.Fullscreen = true
			}
			m.updateSelectedItemText(fmt.Sprintf("%20v: %-10v", "Fullscreen", global.gVariableConfig.Fullscreen))
			global.gVariableConfig.SaveConfiguration()
		})
	m.addItem(0.5, fmt.Sprintf("%20v: %-10v", "Undecorated Window", global.gVariableConfig.UndecoratedWindow),
		func() {
			if global.gVariableConfig.UndecoratedWindow {
				global.gVariableConfig.UndecoratedWindow = false
			} else {
				global.gVariableConfig.UndecoratedWindow = true
			}
			m.updateSelectedItemText(fmt.Sprintf("%20v: %-10v", "Undecorated Window", global.gVariableConfig.UndecoratedWindow))
			global.gVariableConfig.SaveConfiguration()
		})
	m.addItem(0.7, "Back", func() {
		global.gActiveMenu = global.gOptionsMenu

	})
	m.items[0].selected = 1
}

func (m *menu) updateSelectedItemText(text string) {
	for i, v := range m.items {
		if v.selected == 1 {
			m.items[i].name = text
			break
		}
	}
}

func (m *menu) createGame() {
	m.create()
	m.addItem(0.5, fmt.Sprintf("%10v: %-10v", "Max Particles", global.gVariableConfig.MaxParticles), func() {})
	m.addItem(0.7, "Back", func() {
		global.gActiveMenu = global.gOptionsMenu

	})
	m.items[0].selected = 1
}

func (m *menu) addItem(scale float64, str string, f func()) {
	item := &menuItem{
		canvas:   pixelgl.NewCanvas(pixel.R(0, 0, 1, 1)),
		name:     str,
		action:   f,
		selected: 0,
		scale:    scale,
	}
	item.canvas.SetUniform("uSelected", &item.selected)
	item.canvas.SetUniform("uTime", &global.uTime)
	item.canvas.SetFragmentShader(fragmentShaderMenuItem)

	m.items = append(m.items, item)
}

func (m *menu) selectItem() {
	for i, item := range m.items {
		if item.selected == 1 {
			m.items[i].action()
			return
		}
	}
}

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

func (m *menu) draw(dt, elapsed float64) {
	offsetX := 30.0
	extraOffset := 0.0
	if m.logo != nil {
		m.logo.canvas.Clear(pixel.RGBA{R: 0, G: 0, B: 0, A: 0})
		global.gFont.writeToCanvas(m.logo.name, m.logo.canvas)
		m.logo.canvas.Draw(global.gWin, pixel.IM.Scaled(pixel.ZV, m.logo.scale).Moved(pixel.V(global.gCamera.pos.X+wViewMax/2.0-offsetX, global.gCamera.pos.Y+wViewMax/20+200)))
		extraOffset = 50
	}
	for i := range m.items {
		m.items[i].canvas.Clear(pixel.RGBA{R: 0, G: 0, B: 0, A: 0})
		offsetY := 30 * m.items[i].scale
		global.gFont.writeToCanvas(m.items[i].name, m.items[i].canvas)
		m.items[i].canvas.Draw(global.gWin, pixel.IM.Scaled(pixel.ZV, m.items[i].scale).Moved(pixel.V(global.gCamera.pos.X+wViewMax/2.0-offsetX, global.gCamera.pos.Y+wViewMax/2-float64(i)*offsetY-extraOffset)))
	}
}
