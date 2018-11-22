package main

type mapType int

const (
	mapEasy mapType = iota
	mapNormal
	mapHard
	mapHell
	mapWtf
)

type objectType int

const (
	weaponAk47 objectType = iota
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

type animationType int

const (
	animWalk animationType = iota
	animJump
	animClimb
	animShoot
	animIdle
)
