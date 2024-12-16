package gograph

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
}

type SimpleGraph struct {
	id       int64
	edges    map[int64]Edge
	vertices map[int64]Vertex
}

func NewSimpleGraph() *SimpleGraph {
	return &SimpleGraph{
		id:       -1,
		edges:    make(map[int64]Edge),
		vertices: make(map[int64]Vertex),
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
