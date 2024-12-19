package gograph

import (
	"crypto/sha256"
	"encoding/binary"
	"github.com/mtresnik/gomath/pkg/gomath"
	"math"
)

type Vertex interface {
	GetValues() []float64
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
	RemoveEdge(edge Edge)
}

func GetEdge(from Vertex, to Vertex) Edge {
	for _, edge := range from.GetEdges() {
		if ToVertex(edge.To()).Hash() == to.Hash() {
			return edge
		}
	}
	return SimpleEdge{
		from:     from,
		to:       to,
		id:       -1,
		distance: gomath.EuclideanDistance(from, to),
	}
}

func VertexFromSpatial(spatial gomath.Spatial) Vertex {
	return &SimpleVertex{Spatial: spatial, Edges: []Edge{}, id: -1, hash: -1}
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
	hasher := sha256.New()

	values := vertex.GetValues()
	for _, value := range values {
		bits := math.Float64bits(value)
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, bits)
		_, _ = hasher.Write(buf)
	}

	if vertex.Id() > 0 {
		idBuf := make([]byte, 8)
		binary.BigEndian.PutUint64(idBuf, uint64(vertex.Id()))
		_, _ = hasher.Write(idBuf)
	}

	size := vertex.Size()
	sizeBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(sizeBuf, uint64(size))
	_, _ = hasher.Write(sizeBuf)

	fullHash := hasher.Sum(nil)
	hash1 := binary.BigEndian.Uint64(fullHash[:8])
	hash2 := binary.BigEndian.Uint64(fullHash[8:16])

	finalHash := int64(hash1 ^ hash2)
	return finalHash
}

type SimpleVertex struct {
	Spatial gomath.Spatial
	Edges   []Edge
	id      int64
	hash    int64
}

func NewSimpleVertex(spatial gomath.Spatial, edges ...Edge) SimpleVertex {
	return SimpleVertex{Spatial: spatial, Edges: edges, id: -1, hash: -1}
}

func (v *SimpleVertex) GetEdge(to Vertex) Edge {
	return GetEdge(v, to)
}

func (v *SimpleVertex) DistanceTo(other SimpleVertex, distanceFunction ...gomath.DistanceFunction) float64 {
	return gomath.ToPoint(v.Spatial).DistanceTo(gomath.ToPoint(other.Spatial), distanceFunction...)
}

func (v *SimpleVertex) GetValues() []float64 {
	return v.Spatial.GetValues()
}

func (v *SimpleVertex) Size() int {
	return v.Spatial.Size()
}

func (v *SimpleVertex) X() float64 {
	return v.Spatial.X()
}

func (v *SimpleVertex) Y() float64 {
	return v.Spatial.Y()
}

func (v *SimpleVertex) Z() float64 {
	return v.Spatial.Z()
}

func (v *SimpleVertex) W() float64 {
	return v.Spatial.W()
}

func (v *SimpleVertex) Id() int64 {
	return v.id
}

func (v *SimpleVertex) GetEdges() []Edge {
	return v.Edges
}

func (v *SimpleVertex) Hash() int64 {
	if v.hash <= 0 {
		v.hash = HashVertex(v)
	}
	return v.hash
}

func (v *SimpleVertex) AddEdge(edge Edge) {
	v.Edges = append(v.Edges, edge)
}

func (v *SimpleVertex) RemoveEdge(edge Edge) {
	for i, e := range v.Edges {
		if e.Hash() == edge.Hash() {
			v.Edges = append(v.Edges[:i], v.Edges[i+1:]...)
			return
		}
	}
}

type VertexWrapper struct {
	Previous *VertexWrapper
	Inner    Vertex
	Next     *VertexWrapper
	Costs    map[string]CostEntry
	Combined *CostEntry
}

func NewVertexWrapper(inner Vertex, costs map[string]CostEntry, pCostCombiner ...CostCombiner) *VertexWrapper {
	costCombiner := MultiplicativeCostCombiner
	if len(pCostCombiner) > 0 {
		costCombiner = pCostCombiner[0]
	}
	combined := costCombiner(costs)
	return &VertexWrapper{Inner: inner, Costs: costs, Combined: &combined}
}

func (v *VertexWrapper) GetCombined(pCostCombiner ...CostCombiner) *CostEntry {
	if v.Combined != nil {
		return v.Combined
	}
	costCombiner := MultiplicativeCostCombiner
	if len(pCostCombiner) > 0 {
		costCombiner = pCostCombiner[0]
	}
	combined := costCombiner(v.Costs)
	v.Combined = &combined
	return &combined
}

func (v *VertexWrapper) GetValues() []float64 {
	return v.Inner.GetValues()
}

func (v *VertexWrapper) Size() int {
	return v.Inner.Size()
}

func (v *VertexWrapper) X() float64 {
	return v.Inner.X()
}

func (v *VertexWrapper) Y() float64 {
	return v.Inner.Y()
}

func (v *VertexWrapper) Z() float64 {
	return v.Inner.Z()
}

func (v *VertexWrapper) W() float64 {
	return v.Inner.W()
}

func (v *VertexWrapper) Id() int64 {
	return v.Inner.Id()
}

func (v *VertexWrapper) GetEdges() []Edge {
	return v.Inner.GetEdges()
}

func (v *VertexWrapper) Hash() int64 {
	return v.Inner.Hash()
}

func (v *VertexWrapper) GetEdge(to Vertex) Edge {
	return v.Inner.GetEdge(to)
}

func (v *VertexWrapper) AddEdge(edge Edge) {
	v.Inner.AddEdge(edge)
}

func (v *VertexWrapper) RemoveEdge(edge Edge) {
	v.Inner.RemoveEdge(edge)
}
