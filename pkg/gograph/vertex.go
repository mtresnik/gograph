package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
)

type Vertex struct {
	Spatial gomath.Spatial
	id      int64
}

func (v Vertex) DistanceTo(other Vertex, distanceFunction ...gomath.DistanceFunction) float64 {
	return gomath.ToPoint(v.Spatial).DistanceTo(gomath.ToPoint(other.Spatial), distanceFunction...)
}

func (v Vertex) GetValues() []float64 {
	return v.Spatial.GetValues()
}

func (v Vertex) SetValues(values []float64) {
	v.Spatial.SetValues(values)
}

func (v Vertex) Size() int {
	return v.Spatial.Size()
}

func (v Vertex) X() float64 {
	return v.Spatial.X()
}

func (v Vertex) Y() float64 {
	return v.Spatial.Y()
}

func (v Vertex) Z() float64 {
	return v.Spatial.Z()
}

func (v Vertex) W() float64 {
	return v.Spatial.W()
}

func (v Vertex) Id() int64 {
	return v.id
}
