package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
	"testing"
)

func TestMazeToGraphProvider_Build(t *testing.T) {
	response := AldousBroderMazeGenerator(NewMazeGeneratorRequest(10, 10))
	graph := MazeToGraphProvider{response.Maze}.Build()
	print("Num Vertices:", len(graph.GetVertices()), "\tNum Edges:", len(graph.GetEdges()))
}

func TestBoundedGraphProvider_Build(t *testing.T) {
	provider := BoundedGraphProvider{
		BoundingBox: gomath.BoundingBox{10, 10, 20, 20},
		Width:       10,
		Height:      10,
	}
	graph := provider.Build()
	print("Num Vertices:", len(graph.GetVertices()), "\tNum Edges:", len(graph.GetEdges()))
}
