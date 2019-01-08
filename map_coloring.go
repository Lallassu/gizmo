//=============================================================
// map_coloring.go
//-------------------------------------------------------------
// Color palettes for maps
//=============================================================
package main

type mapColor struct {
	background     uint32
	backgroundSoft uint32
	foreground     uint32
	ladders        uint32
	borders        uint32
	entityCodes    map[objectType]entityColor
}

type entityColor struct {
	r uint32
	g uint32
	b uint32
	a uint32
}

func (m *mapColor) generateColors(maptype mapType) {
	// Randomize within palettes different
	// types of coloring schemes depending
	// on map type.
}

func (m *mapColor) create() {
	// default colors
	m.background = 0x000000FF
	m.backgroundSoft = 0x3e585cFF
	m.foreground = 0x3d6253FF
	m.ladders = 0x8b4513FF
	m.borders = 0xFF0000FF

	// Color codes for items in the map (image)
	// Must correspond with the image colors.
	m.entityCodes = make(map[objectType]entityColor)
	m.entityCodes[mobPlayer] = entityColor{0, 0xFFFF, 0, 0}
	m.entityCodes[itemCrate] = entityColor{0xFFFF, 0xFFFF, 0xFFFF, 0}
	m.entityCodes[weaponAk47] = entityColor{0xFFFF, 0xFFFF, 0, 0}
	m.entityCodes[mobEnemy1] = entityColor{0, 0, 0xFFFF, 0}
	m.entityCodes[lampRegular] = entityColor{0, 0, 0xFFFF, 0}
}
