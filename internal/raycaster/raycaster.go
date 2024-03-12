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
	rayDir := target
	rayDir.Sub(start)
	rayDir.Normalize()
	vRayUnitStepSize := vec.Vec2{math.Hypot(1, rayDir[1]/rayDir[0]), math.Hypot(1, rayDir[0]/rayDir[1])}

	vMapCheck := vec.Vec2{math.Floor(start[0]), math.Floor(start[1])}
	vRayLength1D := vec.Vec2{}

	vStep := vec.Vec2{}

	if rayDir[0] < 0 {
		vStep[0] = -1
		vRayLength1D[0] = (start[0] - vMapCheck[0]) * vRayUnitStepSize[0]
	} else {
		vStep[0] = 1
		vRayLength1D[0] = (vMapCheck[0] + 1 - start[0]) * vRayUnitStepSize[0]
	}

	if rayDir[1] < 0 {
		vStep[1] = -1
		vRayLength1D[1] = (start[1] - vMapCheck[1]) * vRayUnitStepSize[1]
	} else {
		vStep[1] = 1
		vRayLength1D[1] = (vMapCheck[1] + 1 - start[1]) * vRayUnitStepSize[1]
	}

	return Ray{start, rayDir, vRayUnitStepSize, vMapCheck, vRayLength1D, vStep}
}

func (r *Ray) Cast(gridMap [][]int) (intersectionPoint vec.Vec2, intersected bool) {
	solidFound := false
	distance := 0.0

	for !solidFound && distance < MAX_DISTANCE {
		if r.Length1D[0] < r.Length1D[1] {
			r.MapCheck[0] += r.Step[0]
			distance = r.Length1D[0]
			r.Length1D[0] += r.UnitStepSize[0]
		} else {
			r.MapCheck[1] += r.Step[1]
			distance = r.Length1D[1]
			r.Length1D[1] += r.UnitStepSize[1]
		}

		if r.MapCheck[0] >= 0 && r.MapCheck[0] < float64(len(gridMap[0])) && r.MapCheck[1] >= 0 && r.MapCheck[1] < float64(len(gridMap)) {
			if gridMap[int(r.MapCheck[1])][int(r.MapCheck[0])] == 1 {
				solidFound = true
			}
		}
	}

	if solidFound {
		intersection := r.Start.Copy()
		vDist := r.Dir.Copy()
		vDist.Scale(distance)
		intersection.Add(vDist)
		return intersection, true
	}
	return vec.Vec2{}, false
}
