package main

type mapType int

const (
	mapEasy mapType = iota
	mapNormal
	mapHard
	mapHell
	mapWtf
)

type weaponType int

const (
	weaponAk47 weaponType = iota
	weaponRocket
	weaponGrenade
)

type particleType int

const (
	particleRegular particleType = iota
	particleBlood
	particleFire
	particleSmoke
)

type entityType int

const (
	entityPlayer entityType = iota
	entityEnemy
	entityObject
	entityChunk
)
