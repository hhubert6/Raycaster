package raycaster

import (
	"math"

	vec "github.com/hhubert6/Raycaster/internal/vector"
)

const (
	MAX_DISTANCE = 100.0
)

type Ray struct {
	Start, Dir, UnitStepSize, MapCheck, Length1D, Step vec.Vec2
}

func NewRayFromTarget(start, target vec.Vec2) Ray {
	rayDir := target.Sub(&start).Normalize()
	return NewRayFromDir(start, *rayDir)
}

func NewRayFromAngle(start vec.Vec2, angle float64) Ray {
	rayDir := vec.Vec2{math.Cos(angle), math.Sin(angle)}
	return NewRayFromDir(start, rayDir)
}

func NewRayFromDir(start, dir vec.Vec2) Ray {
	vRayUnitStepSize := vec.Vec2{math.Hypot(1, dir[1]/dir[0]), math.Hypot(1, dir[0]/dir[1])}
	vMapCheck := vec.Vec2{math.Floor(start[0]), math.Floor(start[1])}
	vRayLength1D := vec.Vec2{}
	vStep := vec.Vec2{}

	if dir[0] < 0 {
		vStep[0] = -1
		vRayLength1D[0] = (start[0] - vMapCheck[0]) * vRayUnitStepSize[0]
	} else {
		vStep[0] = 1
		vRayLength1D[0] = (vMapCheck[0] + 1 - start[0]) * vRayUnitStepSize[0]
	}
	if dir[1] < 0 {
		vStep[1] = -1
		vRayLength1D[1] = (start[1] - vMapCheck[1]) * vRayUnitStepSize[1]
	} else {
		vStep[1] = 1
		vRayLength1D[1] = (vMapCheck[1] + 1 - start[1]) * vRayUnitStepSize[1]
	}

	return Ray{start, dir, vRayUnitStepSize, vMapCheck, vRayLength1D, vStep}
}

func (r *Ray) Cast(gridMap [][]int) (offset float64, distance float64, side int, solidFound int) {
	distance = 0.0

	for solidFound == 0 && distance < MAX_DISTANCE {
		if r.Length1D[0] < r.Length1D[1] {
			r.MapCheck[0] += r.Step[0]
			distance = r.Length1D[0]
			side = 0
			r.Length1D[0] += r.UnitStepSize[0]
		} else {
			r.MapCheck[1] += r.Step[1]
			distance = r.Length1D[1]
			side = 1
			r.Length1D[1] += r.UnitStepSize[1]
		}

		if r.MapCheck[0] >= 0 && r.MapCheck[0] < float64(len(gridMap[0])) && r.MapCheck[1] >= 0 && r.MapCheck[1] < float64(len(gridMap)) {
			if gridMap[int(r.MapCheck[1])][int(r.MapCheck[0])] > 0 {
				solidFound = gridMap[int(r.MapCheck[1])][int(r.MapCheck[0])]
			}
		}
	}

	if solidFound > 0 {
		intersect := r.Start.Copy()
		intersect.Add(r.Dir.Scale(distance))
		offset := 0.0
		if side == 0 {
			if r.Step[0] > 0 {
				offset = intersect[1] - r.MapCheck[1]
			} else {
				offset = r.MapCheck[1] + 1 - intersect[1]
			}
		} else {
			if r.Step[1] > 0 {
				offset = r.MapCheck[0] + 1 - intersect[0]
			} else {
				offset = intersect[0] - r.MapCheck[0]
			}
		}
		return offset, distance, side, solidFound
	}
	return 0.0, math.Inf(1), side, 0
}
