package gograph

import "testing"

func TestMazeToGraphProvider_Build(t *testing.T) {
	response := AldousBroderMazeGenerator{}.Build(NewMazeGeneratorRequest(10, 10))
	graph := MazeToGraphProvider{}.Build(response.Maze)
	print("Num Vertices:", len(graph.GetVertices()), "\tNum Edges:", len(graph.GetEdges()))
}
