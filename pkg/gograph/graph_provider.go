package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
)

type IMazeToGraphProvider interface {
	Build(maze Maze) Graph
}

type MazeToGraphProvider struct{}

func (m MazeToGraphProvider) Build(maze Maze) Graph {
	retGraph := NewSimpleGraph()
	cellToVertexMap := map[int64]Vertex{}

	allCells := maze.Flatten()
	for _, cell := range allCells {
		vertex := SimpleVertex{
			Spatial: gomath.Point{Values: []float64{float64(cell.Col), float64(cell.Row)}},
			Edges:   make([]Edge, 0),
		}
		retGraph.AddVertex(vertex)
		cellToVertexMap[cell.Hash()] = vertex
	}

	allConnections := make([]MazeConnection, 0)
	for _, cell := range allCells {
		tempConnections := maze.GetConnections(cell.Row, cell.Col)
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
			retGraph.AddEdge(edge)
			calls++
		}
	}
	return retGraph
}
