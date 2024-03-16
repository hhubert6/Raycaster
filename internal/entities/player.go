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
	PLAYER_SIZE    = 5
)

type Player struct {
	Pos   vec.Vec2
	Angle float64
}

func (p *Player) Move(dir Direction, gridMap [][]int) {
	var angle float64

	switch dir {
	case FORWARD:
		angle = p.Angle
	case BACKWARD:
		angle = p.Angle + math.Pi
	case LEFT:
		angle = p.Angle - math.Pi/2
	case RIGHT:
		angle = p.Angle + math.Pi/2
	}

	dx, dy := math.Cos(angle), math.Sin(angle)
	p.checkCollision(MOVEMENT_SPEED*dx, MOVEMENT_SPEED*dy, gridMap)
}

func (p *Player) checkCollision(dx, dy float64, gridMap [][]int) {
	d := math.Hypot(dx, dy) * PLAYER_SIZE

	horY, vertX := int(p.Pos[1]), int(p.Pos[0])
	var horX, vertY int

	if dx > 0 {
		horX = int(p.Pos[0] + d)
	} else {
		horX = int(p.Pos[0] - d)
	}
	if dy > 0 {
		vertY = int(p.Pos[1] + d)
	} else {
		vertY = int(p.Pos[1] - d)
	}

	if gridMap[horY][horX] == 0 {
		p.Pos[0] += dx
	}
	if gridMap[vertY][vertX] == 0 {
		p.Pos[1] += dy
	}
}

func (p *Player) Rotate(angle float64) {
	p.Angle += angle
}

func (p *Player) Dir() *vec.Vec2 {
	return &vec.Vec2{math.Cos(p.Angle), math.Sin(p.Angle)}
}

func (p *Player) SetPos(v vec.Vec2) {
	p.Pos = v
}
