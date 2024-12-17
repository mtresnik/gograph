package gograph

import (
	"hash/fnv"
	"maps"
	"time"
)

type Graph interface {
	Id() int64
	GetEdge(id int64) Edge
	GetVertex(id int64) Vertex
	AddEdge(e Edge)
	ContainsEdge(e Edge) bool
	ContainsVertex(v Vertex) bool
	AddVertex(v Vertex)
	GetVertices() []Vertex
	GetEdges() []Edge
	RemoveEdge(e Edge)
	RemoveVertex(v Vertex)
	Size() int
	Clear()
	Hash() int64
}

func GraphHashOrId(graph Graph) int64 {
	if graph.Id() > 0 {
		return graph.Id()
	}
	return graph.Hash()
}

type SimpleGraph struct {
	id       int64
	edges    map[int64]Edge
	vertices map[int64]Vertex
	hash     int64
}

func NewSimpleGraph() *SimpleGraph {
	return &SimpleGraph{
		id:       time.Now().UnixNano(),
		edges:    make(map[int64]Edge),
		vertices: make(map[int64]Vertex),
		hash:     -1,
	}
}

func (g *SimpleGraph) Id() int64 {
	return g.id
}

func (g *SimpleGraph) Size() int {
	return len(g.vertices)
}

func (g *SimpleGraph) Clear() {
	g.edges = make(map[int64]Edge)
	g.vertices = make(map[int64]Vertex)
}

func (g *SimpleGraph) GetEdge(id int64) Edge {
	return g.edges[id]
}

func (g *SimpleGraph) GetVertex(id int64) Vertex {
	return g.vertices[id]
}

func (g *SimpleGraph) AddEdge(e Edge) {
	g.edges[EdgeHashOrId(e)] = e
}

func (g *SimpleGraph) ContainsEdge(e Edge) bool {
	return g.edges[EdgeHashOrId(e)] != nil
}

func (g *SimpleGraph) ContainsVertex(v Vertex) bool {
	return g.vertices[VertexHashOrId(v)] != nil
}

func (g *SimpleGraph) AddVertex(v Vertex) {
	g.vertices[VertexHashOrId(v)] = v
}

func (g *SimpleGraph) GetVertices() []Vertex {
	vertices := make([]Vertex, 0, len(g.vertices))
	for _, v := range g.vertices {
		vertices = append(vertices, v)
	}
	return vertices
}

func (g *SimpleGraph) GetEdges() []Edge {
	edges := make([]Edge, 0, len(g.edges))
	for _, e := range g.edges {
		edges = append(edges, e)
	}
	return edges
}

func (g *SimpleGraph) RemoveEdge(e Edge) {
	delete(g.edges, EdgeHashOrId(e))
}

func (g *SimpleGraph) RemoveVertex(v Vertex) {
	delete(g.vertices, VertexHashOrId(v))
}

func (g *SimpleGraph) SetId(id int64) {
	g.id = id
}

func (g *SimpleGraph) Hash() int64 {
	if g.hash != -1 {
		return g.hash
	}
	vertexKeysIter := maps.Keys(g.vertices)
	vertexKeys := make([]int64, 0, len(g.vertices))
	for vertKey := range vertexKeysIter {
		vertexKeys = append(vertexKeys, vertKey)
	}
	edgeKeysIter := maps.Keys(g.edges)
	edgeKeys := make([]int64, 0, len(g.edges))
	for edgeKey := range edgeKeysIter {
		edgeKeys = append(edgeKeys, edgeKey)
	}
	allKeys := make([]int64, 0, len(g.vertices)+len(g.edges))
	allKeys = append(allKeys, vertexKeys...)
	allKeys = append(allKeys, edgeKeys...)

	hasher := fnv.New64a()
	for _, key := range allKeys {
		var buf [8]byte
		for i := 0; i < 8; i++ {
			buf[i] = byte(key >> (i * 8))
		}
		_, err := hasher.Write(buf[:])
		if err != nil {
			return 0
		}
	}
	g.hash = int64(hasher.Sum64())
	return g.hash
}
