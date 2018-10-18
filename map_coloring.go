//=============================================================
// map_coloring.go
//-------------------------------------------------------------
// Color palettes for maps
//=============================================================
package main

type mapColor struct {
	background uint32
	foreground uint32
	ladders    uint32
	borders    uint32
}

func GenerateMapColor(maptype mapType) *mapColor {
	// Randomize within palettes different
	// types of coloring schemes depending
	// on map type.
	m := &mapColor{
		background: 0xAAF55FFF,
		foreground: 0xFFAADAFF,
		ladders:    0x00FFAAFF,
		borders:    0xFF0000FF,
	}
	return m
}

func (m *mapColor) getBackground() uint32 {
	return m.background
}

func (m *mapColor) getForeground() uint32 {
	return m.foreground
}

func (m *mapColor) getLadder() uint32 {
	return m.ladders
}

func (m *mapColor) getBorder() uint32 {
	return m.borders
}
