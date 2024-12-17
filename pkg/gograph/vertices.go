package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
	"hash/fnv"
	"math"
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
	Hash() int64
	GetEdge(to Vertex) Edge
	AddEdge(edge Edge)
}

type VertexListener interface {
	Visit(v Vertex)
}

func getEdge(from Vertex, to Vertex) Edge {
	for _, edge := range from.GetEdges() {
		if ToVertex(edge.To()).Hash() == to.Hash() {
			return edge
		}
	}
	return SimpleEdge{
		from:     from,
		to:       to,
		id:       -1,
		distance: gomath.EuclideanDistance{}.Eval(from, to),
	}
}

func VertexFromSpatial(spatial gomath.Spatial) Vertex {
	return SimpleVertex{Spatial: spatial, Edges: []Edge{}, id: -1, hash: -1}
}

func ToVertex(vertex interface{}) Vertex {
	if v, ok := vertex.(Vertex); ok {
		return v
	}
	if s, ok := vertex.(gomath.Spatial); ok {
		return VertexFromSpatial(s)
	}
	panic("Cannot cast to Vertex")
}

func VertexHashOrId(vertex Vertex) int64 {
	if vertex.Id() > 0 {
		return vertex.Id()
	}
	return HashVertex(vertex)
}

func HashVertex(vertex Vertex) int64 {
	hasher := fnv.New64a()

	for _, value := range vertex.GetValues() {
		bits := math.Float64bits(value)
		buf := make([]byte, 8)
		for i := 0; i < 8; i++ {
			buf[i] = byte(bits >> (i * 8))
		}
		_, _ = hasher.Write(buf)
	}

	return int64(hasher.Sum64())
}

type SimpleVertex struct {
	Spatial gomath.Spatial
	Edges   []Edge
	id      int64
	hash    int64
}

func NewVertex(spatial gomath.Spatial, edges ...Edge) SimpleVertex {
	return SimpleVertex{Spatial: spatial, Edges: edges, id: -1, hash: -1}
}

func (v SimpleVertex) GetEdge(to Vertex) Edge {
	return getEdge(v, to)
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

func (v SimpleVertex) Hash() int64 {
	if v.hash == -1 {
		v.hash = HashVertex(v)
	}
	return v.hash
}

func (v SimpleVertex) AddEdge(edge Edge) {
	if v.Edges == nil {
		v.Edges = make([]Edge, 0)
	}
	v.Edges = append(v.Edges, edge)
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

func (v VertexWrapper) Hash() int64 {
	return v.Inner.Hash()
}

func (v VertexWrapper) GetEdge(to Vertex) Edge {
	return v.Inner.GetEdge(to)
}

func (v VertexWrapper) AddEdge(edge Edge) {
	v.Inner.AddEdge(edge)
}
