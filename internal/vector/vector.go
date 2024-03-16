package vec

import "math"

type Vec2 [2]float64

func (v *Vec2) Add(other *Vec2) *Vec2 {
	v[0] += other[0]
	v[1] += other[1]
	return v
}

func (v *Vec2) Sub(other *Vec2) *Vec2 {
	v[0] -= other[0]
	v[1] -= other[1]
	return v
}

func (v *Vec2) Length() float64 {
	return math.Hypot(v[0], v[1])
}

func (v *Vec2) Normalize() *Vec2 {
	d := math.Hypot(v[0], v[1])
	v[0] /= d
	v[1] /= d
	return v
}

func (v *Vec2) Copy() *Vec2 {
	return &Vec2{v[0], v[1]}
}

func (v *Vec2) Scale(f float64) *Vec2 {
	v[0] *= f
	v[1] *= f
	return v
}
