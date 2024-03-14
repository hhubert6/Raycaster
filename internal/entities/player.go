package entities

import (
	"math"

	vec "github.com/hhubert6/Raycaster/internal/vector"
)

type Direction int

const (
	FORWARD Direction = iota
	BACKWARD
	RIGHT
	LEFT
	MOVEMENT_SPEED = 0.05
)

type Player struct {
	Pos   vec.Vec2
	Angle float64
}

func (p *Player) Move(dir Direction) {
	switch dir {
	case FORWARD:
		dx, dy := math.Cos(p.Angle), math.Sin(p.Angle)
		p.Pos[0] += MOVEMENT_SPEED * dx
		p.Pos[1] += MOVEMENT_SPEED * dy
	case BACKWARD:
		dx, dy := math.Cos(p.Angle), math.Sin(p.Angle)
		p.Pos[0] -= MOVEMENT_SPEED * dx
		p.Pos[1] -= MOVEMENT_SPEED * dy
	case LEFT:
		dx, dy := math.Cos(p.Angle-math.Pi/2), math.Sin(p.Angle-math.Pi/2)
		p.Pos[0] += MOVEMENT_SPEED * dx
		p.Pos[1] += MOVEMENT_SPEED * dy
	case RIGHT:
		dx, dy := math.Cos(p.Angle+math.Pi/2), math.Sin(p.Angle+math.Pi/2)
		p.Pos[0] += MOVEMENT_SPEED * dx
		p.Pos[1] += MOVEMENT_SPEED * dy
	}
}

func (p *Player) Rotate(angle float64) {
	p.Angle += angle
}

func (p *Player) SetPos(v vec.Vec2) {
	p.Pos = v
}
