package gograph

type VertexVisitor interface {
	Visit(v Vertex)
}

type EdgeVisitor interface {
	Visit(e Edge)
}
