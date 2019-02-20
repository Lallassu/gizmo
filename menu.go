package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type menu struct {
	items            []*menuItem
	logo             *menuItem
	videoModes       []pixelgl.VideoMode
	currentVideoMode int
}

type menuItem struct {
	action   func()
	nextItem func()
	prevItem func()
	name     string
	canvas   *pixelgl.Canvas
	scale    float64
	selected int32
}

func (m *menu) create() {
	m.items = make([]*menuItem, 0)
	m.videoModes = make([]pixelgl.VideoMode, 0)
	for _, mode := range pixelgl.PrimaryMonitor().VideoModes() {
		m.videoModes = append(m.videoModes, mode)
	}
}

func (m *menu) createMain() {
	m.logo = &menuItem{
		canvas:   pixelgl.NewCanvas(pixel.R(0, 0, 1, 1)),
		name:     "Gizmo",
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
			// Change to current video mode
			mode := m.videoModes[m.currentVideoMode]
			global.gVariableConfig.WindowWidth = mode.Width
			global.gVariableConfig.WindowHeight = mode.Height
			global.gWin.SetBounds(pixel.R(0, 0, float64(mode.Width), float64(mode.Height)))
			centerWindow(global.gWin)
			global.gVariableConfig.SaveConfiguration()
		})
	m.items[len(m.items)-1].prevItem = func() {
		m.currentVideoMode--
		if m.currentVideoMode < 0 {
			m.currentVideoMode = len(m.videoModes) - 1
		}
		mode := m.videoModes[m.currentVideoMode]
		m.updateSelectedItemText(fmt.Sprintf("%20v: %-10v", "Resolution", fmt.Sprintf("%v x %v", mode.Width, mode.Height)))
	}
	m.items[len(m.items)-1].nextItem = func() {
		m.currentVideoMode++
		if m.currentVideoMode >= len(m.videoModes) {
			m.currentVideoMode = 0
		}
		mode := m.videoModes[m.currentVideoMode]
		m.updateSelectedItemText(fmt.Sprintf("%20v: %-10v", "Resolution", fmt.Sprintf("%v x %v", mode.Width, mode.Height)))
	}
	m.addItem(0.5, fmt.Sprintf("%20v: %-10v", "V-sync", global.gVariableConfig.Vsync),
		func() {
			if global.gVariableConfig.Vsync {
				global.gVariableConfig.Vsync = false
			} else {
				global.gVariableConfig.Vsync = true
			}

			global.gWin.SetVSync(global.gVariableConfig.Vsync)

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

			// TBD: Toggle fullscreen
			if global.gVariableConfig.Fullscreen {
				global.gWin.SetMonitor(pixelgl.PrimaryMonitor())
				//	global.gWin.SetBounds(global.gWin.Bounds())
			} else {
				global.gWin.SetMonitor(nil)
				global.gWin.SetBounds(pixel.R(0, 0, float64(global.gVariableConfig.WindowWidth), float64(global.gVariableConfig.WindowHeight)))
			}

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

func (m *menu) addItem(scale float64, str string, fAction func()) {
	item := &menuItem{
		canvas:   pixelgl.NewCanvas(pixel.R(0, 0, 1, 1)),
		name:     str,
		action:   fAction,
		selected: 0,
		scale:    scale,
	}
	item.canvas.SetUniform("uSelected", &item.selected)
	item.canvas.SetUniform("uTime", &global.uTime)
	item.canvas.SetFragmentShader(fragmentShaderMenuItem)

	m.items = append(m.items, item)
}

// nextItem used for a specific setting such as resolution
func (m *menu) nextItem() {
	for i, item := range m.items {
		if item.selected == 1 {
			if m.items[i].nextItem != nil {
				m.items[i].nextItem()
			}
			return
		}
	}
}

// prevItem used for a specific setting such as resolution
func (m *menu) prevItem() {
	for i, item := range m.items {
		if item.selected == 1 {
			if m.items[i].prevItem != nil {
				m.items[i].prevItem()
			}
			return
		}
	}
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
	if m.logo != nil {
		m.logo.canvas.Clear(pixel.RGBA{R: 0, G: 0, B: 0, A: 0})
		global.gFont.writeToCanvas(m.logo.name, m.logo.canvas)
		offsetX := (float64(global.gVariableConfig.WindowWidth) / global.gCamera.zoom) / 2 //- m.logo.canvas.Bounds().Max.X/2
		offsetY := (float64(global.gVariableConfig.WindowHeight) / global.gCamera.zoom) - m.logo.canvas.Bounds().Max.Y
		m.logo.canvas.Draw(global.gWin, pixel.IM.Scaled(pixel.ZV, m.logo.scale).Moved(pixel.V(global.gCamera.pos.X+offsetX, global.gCamera.pos.Y+offsetY)))
		//	extraOffset = 50
	}
	for i := range m.items {
		m.items[i].canvas.Clear(pixel.RGBA{R: 0, G: 0, B: 0, A: 0})
		global.gFont.writeToCanvas(m.items[i].name, m.items[i].canvas)
		//itemScale := m.items[i].scale
		offsetX := (float64(global.gVariableConfig.WindowWidth) / global.gCamera.zoom) / 2 //- m.logo.canvas.Bounds().Max.X/2
		offsetY := (float64(global.gVariableConfig.WindowHeight) / 1.5 / global.gCamera.zoom) - m.items[i].canvas.Bounds().Max.Y*float64(i)*m.items[i].scale
		m.items[i].canvas.Draw(global.gWin, pixel.IM.Scaled(pixel.ZV, m.items[i].scale).Moved(pixel.V(global.gCamera.pos.X+offsetX, global.gCamera.pos.Y+offsetY)))
	}
}
