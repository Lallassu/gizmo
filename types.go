package main

type mapType int

const (
	mapEasy mapType = iota
	mapNormal
	mapHard
	mapHell
	mapWtf
)

type particleType int

const (
	particleRegular particleType = iota
	particleBlood
	particleFire
	particleSmoke
)

type chunkType int

const (
	bgChunk chunkType = iota
	fgChunk
)

type animationType int

const (
	animWalk animationType = iota
	animJump
	animClimb
	animShoot
	animIdle
)

type objectType int

const (
	itemCrate objectType = iota
	mobPlayer
	mobEnemy1
	explosiveRegularMine
	explosiveClusterMine
	weaponAk47
	weaponP90
	weaponShotgun
)
