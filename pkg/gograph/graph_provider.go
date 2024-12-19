package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
	"math/rand"
)

type GraphProvider interface {
	Build() Graph
}

type MazeToGraphProvider struct {
	Maze Maze
}

func (m MazeToGraphProvider) Build() Graph {
	retGraph := NewSimpleGraph()
	cellToVertexMap := map[int64]Vertex{}

	allCells := m.Maze.Flatten()
	for _, cell := range allCells {
		vertex := NewSimpleVertex(gomath.Point{Values: []float64{float64(cell.Col), float64(cell.Row)}}, make([]Edge, 0)...)
		cellToVertexMap[cell.Hash()] = &vertex
	}

	allConnections := make([]MazeConnection, 0)
	for _, cell := range allCells {
		tempConnections := m.Maze.GetConnections(cell.Row, cell.Col)
		for _, connection := range tempConnections {
			if !connection.IsWall {
				allConnections = append(allConnections, connection)
			}
		}
	}
	calls := 0
	for _, connection := range allConnections {
		from := cellToVertexMap[HashMazeCoordinate(connection.From)]
		to := cellToVertexMap[HashMazeCoordinate(connection.To)]
		if from != nil && to != nil {
			edge := NewSimpleEdge(from, to, -1)
			from.AddEdge(edge)
			to.AddEdge(edge.Reverse())
			retGraph.AddEdge(edge)
			retGraph.AddEdge(edge.Reverse())
			retGraph.AddVertex(from)
			retGraph.AddVertex(to)
			calls++
		}
	}
	return retGraph
}

type RandomPruneGraphProvider struct {
	InternalProvider GraphProvider
	PruneRatio       float64
}

func (r RandomPruneGraphProvider) Build() Graph {
	graph := r.InternalProvider.Build()
	toRemoveVertices := make([]Vertex, 0)
	toRemoveEdges := make([]Edge, 0)
	for _, vertex := range graph.GetVertices() {
		if rand.Float64() < r.PruneRatio {
			toRemoveVertices = append(toRemoveVertices, vertex)
			for _, other := range graph.GetVertices() {
				localRemove := make([]Edge, 0)
				for _, edge := range other.GetEdges() {
					if vertex.Hash() == ToVertex(edge.To()).Hash() {
						toRemoveEdges = append(toRemoveEdges, edge)
						toRemoveEdges = append(toRemoveEdges, edge.Reverse())
						localRemove = append(localRemove, edge)
					}
				}
				for _, edge := range localRemove {
					other.RemoveEdge(edge)
					vertex.RemoveEdge(edge.Reverse())
				}
			}
		}
	}
	for _, vertex := range toRemoveVertices {
		graph.RemoveVertex(vertex)
	}
	for _, edge := range toRemoveEdges {
		graph.RemoveEdge(edge)
	}
	toRemoveVertices = make([]Vertex, 0)
	for _, vertex := range graph.GetVertices() {
		if len(vertex.GetEdges()) == 0 {
			toRemoveVertices = append(toRemoveVertices, vertex)
		}
	}
	for _, vertex := range toRemoveVertices {
		graph.RemoveVertex(vertex)
	}
	return graph
}

type BoundedGraphProvider struct {
	BoundingBox gomath.BoundingBox
	Width       int
	Height      int
}

func (b BoundedGraphProvider) Build() Graph {
	lengthX := b.BoundingBox.MaxX - b.BoundingBox.MinX
	lengthY := b.BoundingBox.MaxY - b.BoundingBox.MinY
	dx := lengthX / float64(b.Width)
	dy := lengthY / float64(b.Height)

	graph := NewSimpleGraph()

	vertexMatrix := make([][]Vertex, b.Height)
	for ROW := 0; ROW < b.Height; ROW++ {
		vertexMatrix[ROW] = make([]Vertex, b.Width)
		for COL := 0; COL < b.Width; COL++ {
			vertex := NewSimpleVertex(gomath.Point{Values: []float64{float64(COL)*dx + b.BoundingBox.MinX, float64(ROW)*dy + b.BoundingBox.MinY}}, make([]Edge, 0)...)
			vertexMatrix[ROW][COL] = &vertex
		}
	}

	for ROW := 0; ROW < b.Height; ROW++ {
		for COL := 0; COL < b.Width; COL++ {
			if ROW > 0 {
				edge := NewSimpleEdge(vertexMatrix[ROW-1][COL], vertexMatrix[ROW][COL], -1)
				vertexMatrix[ROW-1][COL].AddEdge(edge)
				vertexMatrix[ROW][COL].AddEdge(edge.Reverse())
				graph.AddEdge(edge)
				graph.AddEdge(edge.Reverse())
			}
			if COL > 0 {
				edge := NewSimpleEdge(vertexMatrix[ROW][COL-1], vertexMatrix[ROW][COL], -1)
				vertexMatrix[ROW][COL-1].AddEdge(edge)
				vertexMatrix[ROW][COL].AddEdge(edge.Reverse())
				graph.AddEdge(edge)
				graph.AddEdge(edge.Reverse())
			}
		}
	}
	for _, row := range vertexMatrix {
		for _, vertex := range row {
			graph.AddVertex(vertex)
		}
	}
	return graph
}
