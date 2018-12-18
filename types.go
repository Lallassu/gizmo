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
	ak47 weaponType = iota
	p90
	shotgun
)

type objectType int

const (
	objectWeapon objectType = iota
	objectCrate
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
