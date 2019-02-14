package main

import (
	"math/rand"

	"github.com/faiface/pixel"
)

// ai used for mobs
type ai struct {
	entity     entity
	dirX       float64
	dirY       float64
	objList    []pixel.Vec
	updateTime float64
}

// creates initiates a new ai entity.
func (a *ai) create(e entity) {
	a.entity = e
	a.dirX = 0.01
	a.objList = []pixel.Vec{}
}

//updateObjectList updates the information where weapons/objects exists in the world.
func (a *ai) updateObjectList() {
	// Get all weapons within view range.
	m, ok := a.entity.(*mob)
	if !ok {
		return
	}

	a.objList = []pixel.Vec{}
	for _, v := range global.gWorld.qt.RetrieveIntersections(&Bounds{X: m.bounds.X, Y: m.bounds.Y, Width: 300, Height: 300}) {
		//if v.entity.getType() == entityObject {
		switch item := v.entity.(type) {
		case *weapon:
			if item.isFree() {
				pos := item.getPosition()
				a.objList = append(a.objList, pos)
			}
		}
		//	}
	}
}

// findWeapon let the AI try to find a weapon
func (a *ai) findWeapon(dt float64) {

	// If at weapon position. Call m.pickup()
	closest := 0.0
	find := -1
	ePos := a.entity.getPosition()
	for i, o := range a.objList {
		dist := distance(o, ePos)
		if closest > dist || i == 0 {
			closest = dist
			find = i
		}
	}
	if find == -1 {
		return
	}

	if closest < 20 {
		if m, ok := a.entity.(*mob); ok {
			m.pickup()
		}
		return
	}

	// Go towards closest
	if ePos.X > a.objList[find].X {
		a.dirX = -dt
	} else {
		a.dirX = dt
	}

	if ePos.Y > a.objList[find].Y {
		a.dirY = -dt
	} else {
		a.dirY = dt
	}

}

func (a *ai) update(dt, time float64) {
	var m *mob
	var ok bool
	if m, ok = a.entity.(*mob); !ok {
		return
	}

	if m.carry == nil {
		a.updateTime += dt
		if a.updateTime > 5.0 {
			a.updateObjectList()
			a.updateTime = 0
		} else if a.updateTime < 1.0 {
			a.findWeapon(dt)
		}
	} else {
		// Find player
		distToPlayer := distance(m.getPosition(), global.gPlayer.getPosition())
		if distToPlayer < 100 {
			if rand.Float64() < 0.2 {
				m.shoot()
			}
		}
	}

	if m.hitLeftWall {
		a.dirX = -dt
	} else if m.hitRightWall {
		a.dirX = dt
		//} else if m.IsOnLadder() {
		//a.dir_y = dt
	} else {
		a.dirY = -dt
	}

	// Jump randomly
	if rand.Float64() < 0.01 {
		a.dirY = dt
	}

	m.move(a.dirX, a.dirY)
}
