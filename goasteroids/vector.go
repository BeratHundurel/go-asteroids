package goasteroids

import "math"

type Vector struct {
	X float64
	Y float64
}

func (v Vector) Normalize() Vector {
	magnitude := math.Sqrt(v.X*v.X + v.Y*v.Y)
	return Vector{X: v.X / magnitude, Y: v.Y / magnitude}
}
