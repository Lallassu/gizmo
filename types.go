package main

type mapType int

const (
	mapEasy mapType = iota
	mapNormal
	mapHard
	mapHell
	mapWtf
)

type ammoType int

const (
	ammoRegular ammoType = iota
	ammoMissile
	ammoShotgun
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

type entityType int

const (
	entityPlayer entityType = iota
	entityEnemy
	entityObject
	entityChunk
)

type animationType int

const (
	animWalk animationType = iota
	animJump
	animClimb
	animShoot
	animIdle
)
