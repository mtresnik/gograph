package gograph

import (
	"fmt"
	"github.com/mtresnik/gomath/pkg/gomath"
)

type Edge interface {
	From() gomath.Spatial
	To() gomath.Spatial
	Reverse() Edge
	Id() int64
}

type SimpleEdge struct {
	from Vertex
	to   Vertex
	id   int64
}

func (e SimpleEdge) From() gomath.Spatial {
	return e.from
}

func (e SimpleEdge) To() gomath.Spatial {
	return e.to
}

func (e SimpleEdge) Reverse() Edge {
	return SimpleEdge{e.to, e.from, e.id}
}

func (e SimpleEdge) Id() int64 {
	return e.id
}

func (e SimpleEdge) String() string {
	return fmt.Sprintf("Edge from %v to %v", e.From(), e.To())
}

type PolyEdge struct {
	Vertices []Vertex
	id       int64
}

func (e PolyEdge) From() gomath.Spatial {
	return e.Vertices[0]
}

func (e PolyEdge) To() gomath.Spatial {
	return e.Vertices[len(e.Vertices)-1]
}

func (e PolyEdge) Reverse() Edge {
	reversedVertices := make([]Vertex, len(e.Vertices))
	for i := 0; i < len(e.Vertices); i++ {
		reversedVertices[i] = e.Vertices[len(e.Vertices)-1-i]
	}
	return PolyEdge{Vertices: reversedVertices, id: e.id}
}

func (e PolyEdge) Id() int64 {
	return e.id
}
