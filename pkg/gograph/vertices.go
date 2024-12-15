package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
)

type Vertex interface {
	GetValues() []float64
	SetValues([]float64)
	Size() int
	X() float64
	Y() float64
	Z() float64
	W() float64
	Id() int64
	GetEdges() []Edge
}

type SimpleVertex struct {
	Spatial gomath.Spatial
	Edges   []Edge
	id      int64
}

func NewVertex(spatial gomath.Spatial, edges ...Edge) SimpleVertex {
	return SimpleVertex{Spatial: spatial, Edges: edges, id: -1}
}

func (v SimpleVertex) DistanceTo(other SimpleVertex, distanceFunction ...gomath.DistanceFunction) float64 {
	return gomath.ToPoint(v.Spatial).DistanceTo(gomath.ToPoint(other.Spatial), distanceFunction...)
}

func (v SimpleVertex) GetValues() []float64 {
	return v.Spatial.GetValues()
}

func (v SimpleVertex) SetValues(values []float64) {
	v.Spatial.SetValues(values)
}

func (v SimpleVertex) Size() int {
	return v.Spatial.Size()
}

func (v SimpleVertex) X() float64 {
	return v.Spatial.X()
}

func (v SimpleVertex) Y() float64 {
	return v.Spatial.Y()
}

func (v SimpleVertex) Z() float64 {
	return v.Spatial.Z()
}

func (v SimpleVertex) W() float64 {
	return v.Spatial.W()
}

func (v SimpleVertex) Id() int64 {
	return v.id
}

func (v SimpleVertex) GetEdges() []Edge {
	return v.Edges
}

type VertexWrapper struct {
	Previous *VertexWrapper
	Inner    Vertex
	Next     *VertexWrapper
	Costs    map[string]CostEntry
}

func NewVertexWrapper(inner Vertex, costs map[string]CostEntry) VertexWrapper {
	return VertexWrapper{Inner: inner, Costs: costs}
}

func (v VertexWrapper) GetValues() []float64 {
	return v.Inner.GetValues()
}

func (v VertexWrapper) SetValues(newValues []float64) {
	v.Inner.SetValues(newValues)
}

func (v VertexWrapper) Size() int {
	return v.Inner.Size()
}

func (v VertexWrapper) X() float64 {
	return v.Inner.X()
}

func (v VertexWrapper) Y() float64 {
	return v.Inner.Y()
}

func (v VertexWrapper) Z() float64 {
	return v.Inner.Z()
}

func (v VertexWrapper) W() float64 {
	return v.Inner.W()
}

func (v VertexWrapper) Id() int64 {
	return v.Inner.Id()
}

func (v VertexWrapper) GetEdges() []Edge {
	return v.Inner.GetEdges()
}
